<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { userApi, reportApi } from "@/api";
import { useAuthStore } from "@/stores";
import type { User } from "@/types";
import { toast } from "vue-sonner";
import { getApiErrorMessage } from "@/lib/utils";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Lock, Crown } from "@lucide/vue";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const educationOptions = [
  "小学",
  "初中",
  "高中",
  "中专",
  "大专",
  "本科",
  "硕士",
  "博士",
  "其他",
];

interface VipInfo {
  userId: string;
  username: string;
  vipLevelName?: string;
  discountRate?: number;
  points: number;
  pointsToNextLevel?: number;
}

const auth = useAuthStore();

const user = ref<User | null>(null);
const vipInfo = ref<VipInfo | null>(null);

const userLoading = ref(false);
const userError = ref(false);
const vipLoading = ref(false);
const savingProfile = ref(false);
const changingPassword = ref(false);

const profileForm = ref({
  phone: "",
  email: "",
  realName: "",
  idCard: "",
  occupation: "",
  education: "",
  income: undefined as number | undefined,
});

const passwordForm = ref({
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
});

const passwordErrors = reactive<{
  newPassword?: string;
  confirmPassword?: string;
}>({});

function vipLevelBadgeClass(name?: string): string {
  if (!name) return "bg-gray-100 text-gray-800 border-gray-200";
  const lower = name.toLowerCase();
  if (lower.includes("gold") || lower.includes("金"))
    return "bg-yellow-100 text-yellow-800 border-yellow-200";
  if (lower.includes("silver") || lower.includes("银"))
    return "bg-gray-200 text-gray-700 border-gray-300";
  if (lower.includes("bronze") || lower.includes("铜"))
    return "bg-orange-100 text-orange-800 border-orange-200";
  return "bg-gray-100 text-gray-800 border-gray-200";
}

async function fetchUser() {
  const userId = auth.user?.id;
  if (!userId) return;
  userLoading.value = true;
  userError.value = false;
  try {
    const res = await userApi.getById(userId);
    if (res.data.data) {
      user.value = res.data.data;
      profileForm.value = {
        phone: res.data.data.phone ?? "",
        email: res.data.data.email ?? "",
        realName: res.data.data.realName ?? "",
        idCard: res.data.data.idCard ?? "",
        occupation: res.data.data.occupation ?? "",
        education: res.data.data.education ?? "",
        income: res.data.data.income ?? undefined,
      };
    }
  } catch (e: unknown) {
    userError.value = true;
    toast.error(getApiErrorMessage(e, "获取用户信息失败"));
  } finally {
    userLoading.value = false;
  }
}

async function fetchVipInfo() {
  const userId = auth.user?.id;
  if (!userId) return;
  vipLoading.value = true;
  try {
    const res = await reportApi.userVipList({ userId });
    if (res.data.data) {
      const list = res.data.data as VipInfo[];
      vipInfo.value = list.length > 0 ? list[0] : null;
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "获取VIP信息失败"));
  } finally {
    vipLoading.value = false;
  }
}

async function saveProfile() {
  const userId = auth.user?.id;
  if (!userId) return;
  savingProfile.value = true;
  try {
    await userApi.update(userId, {
      phone: profileForm.value.phone || undefined,
      email: profileForm.value.email || undefined,
      realName: profileForm.value.realName || undefined,
      idCard: profileForm.value.idCard || undefined,
      occupation: profileForm.value.occupation || undefined,
      education: profileForm.value.education || undefined,
      income:
        typeof profileForm.value.income === "number"
          ? profileForm.value.income
          : undefined,
    });
    toast.success("个人信息已更新");
    await fetchUser();
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "更新失败"));
  } finally {
    savingProfile.value = false;
  }
}

function validatePassword(): boolean {
  passwordErrors.newPassword = undefined;
  passwordErrors.confirmPassword = undefined;
  let valid = true;

  if (passwordForm.value.newPassword.length < 6) {
    passwordErrors.newPassword = "新密码至少6个字符";
    valid = false;
  }
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    passwordErrors.confirmPassword = "两次输入的密码不一致";
    valid = false;
  }

  return valid;
}

async function changePassword() {
  const userId = auth.user?.id;
  if (!userId || !validatePassword()) return;
  changingPassword.value = true;
  try {
    await userApi.update(userId, {
      oldPassword: passwordForm.value.oldPassword,
      password: passwordForm.value.newPassword,
    });
    toast.success("密码已修改");
    passwordForm.value = {
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    };
    passwordErrors.newPassword = undefined;
    passwordErrors.confirmPassword = undefined;
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "修改密码失败"));
  } finally {
    changingPassword.value = false;
  }
}

onMounted(() => {
  fetchUser();
  fetchVipInfo();
});
</script>

<template>
  <div class="container mx-auto max-w-5xl">
    <h1 class="mb-6 text-2xl font-semibold text-gray-900">个人中心</h1>

    <Tabs default-value="profile" class="w-full">
      <TabsList>
        <TabsTrigger value="profile">个人信息</TabsTrigger>
        <TabsTrigger value="vip">VIP信息</TabsTrigger>
        <TabsTrigger value="password">修改密码</TabsTrigger>
      </TabsList>

      <!-- Tab 1: Personal Info -->
      <TabsContent value="profile" class="mt-4">
        <Card>
          <CardHeader>
            <div class="flex items-center gap-4">
              <Avatar class="h-14 w-14">
                <AvatarFallback class="text-lg">
                  {{
                    user?.username?.charAt(0)?.toUpperCase() ??
                    (userLoading ? "?" : "!")
                  }}
                </AvatarFallback>
              </Avatar>
              <div>
                <CardTitle class="text-lg">
                  <span
                    v-if="userLoading"
                    class="inline-block h-5 w-24 animate-pulse rounded bg-muted"
                  />
                  <span v-else-if="userError">无法加载用户信息</span>
                  <span v-else>{{ user?.username ?? "未知用户" }}</span>
                </CardTitle>
                <CardDescription class="text-gray-500"
                  >用户名不可修改</CardDescription
                >
              </div>
            </div>
          </CardHeader>
          <Separator />
          <CardContent class="space-y-4 pt-6">
            <div v-if="userLoading" class="space-y-3">
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
            </div>
            <FieldGroup v-else class="gap-4">
              <Field>
                <FieldLabel for="phone">手机号</FieldLabel>
                <Input
                  id="phone"
                  v-model="profileForm.phone"
                  placeholder="请输入手机号"
                  type="tel"
                />
                <FieldDescription>用于接收预订通知和找回密码</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="email">邮箱</FieldLabel>
                <Input
                  id="email"
                  v-model="profileForm.email"
                  placeholder="请输入邮箱"
                  type="email"
                />
                <FieldDescription>用于接收订单确认和优惠信息</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="realName">真实姓名</FieldLabel>
                <Input
                  id="realName"
                  v-model="profileForm.realName"
                  placeholder="请输入真实姓名"
                />
                <FieldDescription>预订时用于确认入住人身份</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="idCard">身份证号</FieldLabel>
                <Input
                  id="idCard"
                  v-model="profileForm.idCard"
                  placeholder="请输入身份证号"
                />
                <FieldDescription>用于酒店入住登记验证</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="education">学历</FieldLabel>
                <Select v-model="profileForm.education">
                  <SelectTrigger id="education">
                    <SelectValue placeholder="请选择学历" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="opt in educationOptions"
                      :key="opt"
                      :value="opt"
                    >
                      {{ opt }}
                    </SelectItem>
                  </SelectContent>
                </Select>
                <FieldDescription>用于客户偏好统计分析</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="occupation">职业</FieldLabel>
                <Input
                  id="occupation"
                  v-model="profileForm.occupation"
                  placeholder="请输入职业"
                />
                <FieldDescription>用于客户偏好统计分析</FieldDescription>
              </Field>
              <Field>
                <FieldLabel for="income">月收入</FieldLabel>
                <Input
                  id="income"
                  v-model="profileForm.income"
                  placeholder="请输入月收入"
                  type="number"
                  step="0.01"
                  min="0"
                />
                <FieldDescription>用于客户偏好统计分析</FieldDescription>
              </Field>
            </FieldGroup>
          </CardContent>
          <CardFooter v-if="!userLoading">
            <Button
              class="bg-blue-900 hover:bg-blue-800"
              :disabled="savingProfile"
              @click="saveProfile"
            >
              {{ savingProfile ? "保存中..." : "保存修改" }}
            </Button>
          </CardFooter>
        </Card>
      </TabsContent>

      <!-- Tab 2: VIP Info -->
      <TabsContent value="vip" class="mt-4">
        <Card>
          <CardHeader>
            <CardTitle class="text-lg flex items-center gap-2">
              <Crown class="h-5 w-5" />
              VIP信息
            </CardTitle>
            <CardDescription class="text-gray-500"
              >查看您的会员等级和积分信息</CardDescription
            >
          </CardHeader>
          <Separator />
          <CardContent class="pt-6">
            <div v-if="vipLoading" class="space-y-3">
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
              <div class="h-10 w-full animate-pulse rounded bg-muted" />
            </div>
            <div v-else-if="!vipInfo" class="py-8 text-center text-gray-500">
              暂无VIP信息
            </div>
            <div v-else class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-500">会员等级</span>
                <Badge
                  v-if="vipInfo.vipLevelName"
                  :class="vipLevelBadgeClass(vipInfo.vipLevelName)"
                >
                  {{ vipInfo.vipLevelName }}
                </Badge>
                <span v-else class="text-gray-400">-</span>
              </div>
              <Separator />
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-500">当前积分</span>
                <span class="text-lg font-semibold">{{ vipInfo.points }}</span>
              </div>
              <Separator />
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-500">折扣率</span>
                <span class="font-medium">
                  {{
                    vipInfo.discountRate
                      ? `${(vipInfo.discountRate * 100).toFixed(0)}%`
                      : "-"
                  }}
                </span>
              </div>
              <Separator />
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-500"
                  >距下一等级积分</span
                >
                <span class="font-medium">
                  {{ vipInfo.pointsToNextLevel ?? "-" }}
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <!-- Tab 3: Change Password -->
      <TabsContent value="password" class="mt-4">
        <Card>
          <CardHeader>
            <CardTitle class="text-lg flex items-center gap-2">
              <Lock class="h-5 w-5" />
              修改密码
            </CardTitle>
            <CardDescription class="text-gray-500"
              >定期修改密码有助于保护账户安全</CardDescription
            >
          </CardHeader>
          <Separator />
          <CardContent class="pt-6">
            <FieldGroup class="gap-4">
              <Field>
                <FieldLabel for="oldPassword">当前密码</FieldLabel>
                <Input
                  id="oldPassword"
                  v-model="passwordForm.oldPassword"
                  placeholder="请输入当前密码"
                  type="password"
                />
                <FieldDescription>输入您的当前密码以验证身份</FieldDescription>
              </Field>
              <Field :data-invalid="!!passwordErrors.newPassword || undefined">
                <FieldLabel for="newPassword">新密码</FieldLabel>
                <Input
                  id="newPassword"
                  v-model="passwordForm.newPassword"
                  placeholder="请输入新密码（至少6个字符）"
                  type="password"
                  :aria-invalid="!!passwordErrors.newPassword || undefined"
                />
                <FieldDescription
                  >建议使用字母、数字和特殊字符的组合</FieldDescription
                >
                <FieldError v-if="passwordErrors.newPassword">
                  {{ passwordErrors.newPassword }}
                </FieldError>
              </Field>
              <Field
                :data-invalid="!!passwordErrors.confirmPassword || undefined"
              >
                <FieldLabel for="confirmPassword">确认新密码</FieldLabel>
                <Input
                  id="confirmPassword"
                  v-model="passwordForm.confirmPassword"
                  placeholder="请再次输入新密码"
                  type="password"
                  :aria-invalid="!!passwordErrors.confirmPassword || undefined"
                />
                <FieldError v-if="passwordErrors.confirmPassword">
                  {{ passwordErrors.confirmPassword }}
                </FieldError>
              </Field>
            </FieldGroup>
          </CardContent>
          <CardFooter>
            <Button
              class="bg-blue-900 hover:bg-blue-800"
              :disabled="changingPassword"
              @click="changePassword"
            >
              {{ changingPassword ? "修改中..." : "确认修改" }}
            </Button>
          </CardFooter>
        </Card>
      </TabsContent>
    </Tabs>
  </div>
</template>
