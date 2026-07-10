<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useHotelStore, useAuthStore, useOrderStore } from "@/stores";
import { roomApi } from "@/api";
import type { CreateOrderRequest, Room } from "@/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Skeleton } from "@/components/ui/skeleton";
import { Separator } from "@/components/ui/separator";
import { toast } from "vue-sonner";
import { getApiErrorMessage } from "@/lib/utils";
import { Calendar, Users, ArrowLeft } from "@lucide/vue";

const route = useRoute();
const router = useRouter();
const hotelStore = useHotelStore();
const authStore = useAuthStore();
const orderStore = useOrderStore();

const roomId = computed(() => route.params.roomId as string);

const room = ref<Room | null>(null);
const submitting = ref(false);

const form = ref({
  checkInDate: "",
  checkOutDate: "",
  guestName: "",
  guestPhone: "",
  guestIdCard: "",
  roomCount: 1,
});

const today = computed(() => {
  const d = new Date();
  return d.toISOString().split("T")[0];
});

const minCheckOut = computed(() => {
  if (!form.value.checkInDate) return today.value;
  const d = new Date(form.value.checkInDate);
  d.setDate(d.getDate() + 1);
  return d.toISOString().split("T")[0];
});

const nights = computed(() => {
  if (!form.value.checkInDate || !form.value.checkOutDate) return 0;
  const inDate = new Date(form.value.checkInDate);
  const outDate = new Date(form.value.checkOutDate);
  const diff = Math.ceil(
    (outDate.getTime() - inDate.getTime()) / (1000 * 60 * 60 * 24),
  );
  return diff > 0 ? diff : 0;
});

const totalPrice = computed(() => {
  if (!room.value || nights.value === 0) return 0;
  return room.value.price * nights.value * form.value.roomCount;
});

const isFormValid = computed(() => {
  return (
    form.value.checkInDate &&
    form.value.checkOutDate &&
    form.value.guestName.trim() &&
    form.value.guestPhone.trim() &&
    form.value.guestIdCard.trim() &&
    form.value.roomCount > 0 &&
    nights.value > 0
  );
});

async function handleSubmit() {
  if (!isFormValid.value || !room.value) return;

  submitting.value = true;
  try {
    const data: CreateOrderRequest = {
      roomId: roomId.value,
      checkInDate: form.value.checkInDate,
      checkOutDate: form.value.checkOutDate,
      guestName: form.value.guestName.trim(),
      guestPhone: form.value.guestPhone.trim(),
      guestIdCard: form.value.guestIdCard.trim(),
      roomCount: form.value.roomCount,
      totalPrice: totalPrice.value,
    };

    const res = await orderStore.createOrder(data);
    if (res.success) {
      toast.success("预订成功！");
      await router.push("/orders");
    } else {
      toast.error(res.message || "预订失败，请重试");
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, "预订失败，请重试"));
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  if (!authStore.isLoggedIn) {
    router.push({ name: "Login", query: { redirect: route.fullPath } });
    return;
  }

  const existingRoom = hotelStore.rooms.find((r) => r.id === roomId.value);
  if (existingRoom) {
    room.value = existingRoom;
  } else {
    try {
      const res = await roomApi.getById(roomId.value);
      if (res.data.success && res.data.data) {
        room.value = res.data.data;
        if (room.value.hotelId) {
          await hotelStore.fetchHotelById(room.value.hotelId);
        }
      }
    } catch (e: unknown) {
      toast.error(getApiErrorMessage(e, "无法加载房间信息"));
    }
  }
});
</script>

<template>
  <div class="container mx-auto max-w-3xl">
    <!-- Back -->
    <Button
      variant="ghost"
      size="sm"
      class="mb-4 text-gray-600"
      @click="router.back()"
    >
      <ArrowLeft class="mr-1 h-4 w-4" />
      返回
    </Button>

    <!-- Loading -->
    <div v-if="!room" class="space-y-4">
      <Skeleton class="h-8 w-1/3" />
      <Skeleton class="h-64 w-full" />
    </div>

    <!-- Booking Form -->
    <div v-else>
      <h1 class="mb-6 text-2xl font-semibold text-gray-900">预订确认</h1>

      <div class="grid gap-6 lg:grid-cols-5">
        <!-- Left: Form -->
        <div class="lg:col-span-3">
          <Card>
            <CardHeader>
              <CardTitle class="text-lg">入住信息</CardTitle>
            </CardHeader>
            <CardContent class="space-y-5">
              <!-- Dates -->
              <div class="grid gap-4 sm:grid-cols-2">
                <div>
                  <Label class="mb-1.5 flex items-center gap-1.5 text-sm">
                    <Calendar class="h-3.5 w-3.5" />
                    入住日期
                  </Label>
                  <Input type="date" v-model="form.checkInDate" :min="today" />
                </div>
                <div>
                  <Label class="mb-1.5 flex items-center gap-1.5 text-sm">
                    <Calendar class="h-3.5 w-3.5" />
                    退房日期
                  </Label>
                  <Input
                    type="date"
                    v-model="form.checkOutDate"
                    :min="minCheckOut"
                  />
                </div>
              </div>

              <!-- Room Count -->
              <div>
                <Label class="mb-1.5 flex items-center gap-1.5 text-sm">
                  <Users class="h-3.5 w-3.5" />
                  房间数量
                </Label>
                <Input
                  type="number"
                  v-model.number="form.roomCount"
                  :min="1"
                  :max="room.availableQuantity"
                />
              </div>

              <Separator />

              <!-- Guest Info -->
              <div class="space-y-4">
                <h3 class="text-sm font-medium text-gray-900">住客信息</h3>

                <div>
                  <Label class="mb-1.5 block text-sm">姓名</Label>
                  <Input
                    v-model="form.guestName"
                    placeholder="请输入入住人姓名"
                  />
                </div>

                <div>
                  <Label class="mb-1.5 block text-sm">手机号</Label>
                  <Input
                    v-model="form.guestPhone"
                    placeholder="请输入手机号"
                    type="tel"
                  />
                </div>

                <div>
                  <Label class="mb-1.5 block text-sm">身份证号</Label>
                  <Input
                    v-model="form.guestIdCard"
                    placeholder="请输入身份证号"
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Right: Summary -->
        <div class="lg:col-span-2">
          <Card class="sticky top-20">
            <CardHeader>
              <CardTitle class="text-lg">订单摘要</CardTitle>
            </CardHeader>
            <CardContent class="space-y-4">
              <!-- Room Summary -->
              <div>
                <p class="font-medium text-gray-900">{{ room.typeName }}</p>
                <p class="text-sm text-gray-500">¥{{ room.price }} / 晚</p>
              </div>

              <Separator />

              <!-- Booking Details -->
              <div class="space-y-2 text-sm">
                <div class="flex justify-between text-gray-600">
                  <span>入住</span>
                  <span>{{ form.checkInDate || "—" }}</span>
                </div>
                <div class="flex justify-between text-gray-600">
                  <span>退房</span>
                  <span>{{ form.checkOutDate || "—" }}</span>
                </div>
                <div class="flex justify-between text-gray-600">
                  <span>房间数</span>
                  <span>{{ form.roomCount }} 间</span>
                </div>
                <div class="flex justify-between text-gray-600">
                  <span>入住天数</span>
                  <span>{{ nights }} 晚</span>
                </div>
              </div>

              <Separator />

              <!-- Total -->
              <div class="flex items-center justify-between">
                <span class="text-sm font-medium text-gray-700">总计</span>
                <span class="text-2xl font-bold text-blue-900">
                  ¥{{ totalPrice.toFixed(2) }}
                </span>
              </div>

              <!-- Submit -->
              <Button
                class="w-full bg-blue-900 hover:bg-blue-800"
                :disabled="!isFormValid || submitting"
                @click="handleSubmit"
              >
                {{ submitting ? "提交中..." : "确认预订" }}
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>
