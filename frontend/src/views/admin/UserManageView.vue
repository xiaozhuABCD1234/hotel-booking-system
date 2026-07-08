<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/user'
import type { User } from '@/types'
import { toast } from 'vue-sonner'
import { getApiErrorMessage } from '@/lib/utils'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from '@/components/ui/dialog'
import { Users, Search, X } from '@lucide/vue'

const users = ref<User[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(1)
const totalItems = ref(0)

const editDialogOpen = ref(false)
const editUser = ref<User | null>(null)
const editRole = ref<'customer' | 'vip' | 'hotel_manager' | 'admin'>('customer')

const deleteDialogOpen = ref(false)
const deleteTargetUser = ref<User | null>(null)

const roleBadgeClass = (role: string): string => {
  switch (role) {
    case 'admin': return 'bg-purple-100 text-purple-800 border-purple-200'
    case 'hotel_manager': return 'bg-blue-100 text-blue-800 border-blue-200'
    case 'vip': return 'bg-amber-100 text-amber-800 border-amber-200'
    default: return 'bg-gray-100 text-gray-800 border-gray-200'
  }
}

const roleLabel = (role: string): string => {
  switch (role) {
    case 'admin': return '系统管理员'
    case 'hotel_manager': return '酒店管理员'
    case 'vip': return 'VIP用户'
    default: return '普通用户'
  }
}

const statusBadgeClass = (status: number): string => {
  if (status === 1) return 'bg-green-100 text-green-800 border-green-200'
  return 'bg-red-100 text-red-800 border-red-200'
}

const statusLabel = (status: number): string => {
  if (status === 1) return '正常'
  return '禁用'
}

const vipLevelLabel = (levelId?: number): string => {
  if (!levelId || levelId === 0) return '-'
  return `Lv.${levelId}`
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await userApi.list(currentPage.value, pageSize.value)
    if (res.data.data) {
      users.value = res.data.data ?? []
      const p = res.data.pagination
      if (p) {
        totalPages.value = p.totalPages
        totalItems.value = p.totalItems
      }
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '获取用户列表失败'))
  } finally {
    loading.value = false
  }
}

function openEdit(user: User) {
  editUser.value = user
  editRole.value = user.role
  editDialogOpen.value = true
}

async function confirmEdit() {
  if (!editUser.value) return
  try {
    await userApi.update(editUser.value.id, { role: editRole.value })
    toast.success('用户角色已更新')
    editDialogOpen.value = false
    await fetchUsers()
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '更新角色失败'))
  }
}

function openDelete(user: User) {
  deleteTargetUser.value = user
  deleteDialogOpen.value = true
}

async function confirmDelete() {
  if (!deleteTargetUser.value) return
  try {
    await userApi.delete(deleteTargetUser.value.id)
    toast.success('用户已删除')
    deleteDialogOpen.value = false
    await fetchUsers()
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '删除用户失败'))
  }
}

function goToPage(page: number) {
  currentPage.value = page
  fetchUsers()
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold flex items-center gap-2">
        <Users class="h-6 w-6" />
        用户管理
      </h1>
      <div class="text-sm text-muted-foreground">
        共 {{ totalItems }} 位用户
      </div>
    </div>

    <Card>
      <CardContent class="pt-6">
        <div v-if="loading" class="space-y-3">
          <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
        </div>
        <div v-else class="rounded-md border overflow-x-auto">
          <Table class="w-full">
            <TableHeader>
              <TableRow>
                <TableHead class="whitespace-nowrap">用户名</TableHead>
                <TableHead class="whitespace-nowrap">手机号</TableHead>
                <TableHead class="whitespace-nowrap">邮箱</TableHead>
                <TableHead class="whitespace-nowrap">角色</TableHead>
                <TableHead class="text-right whitespace-nowrap">积分</TableHead>
                <TableHead class="whitespace-nowrap">VIP等级</TableHead>
                <TableHead class="whitespace-nowrap">状态</TableHead>
                <TableHead class="text-right whitespace-nowrap w-[100px]">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="users.length === 0">
                <TableCell colspan="8" class="text-center py-8 text-muted-foreground">
                  暂无数据
                </TableCell>
              </TableRow>
              <TableRow
                v-for="user in users"
                :key="user.id"
                class="hover:bg-muted/50"
              >
                <TableCell class="font-medium whitespace-nowrap">{{ user.username }}</TableCell>
                <TableCell class="whitespace-nowrap">{{ user.phone ?? '-' }}</TableCell>
                <TableCell class="whitespace-nowrap">{{ user.email ?? '-' }}</TableCell>
                <TableCell class="whitespace-nowrap">
                  <Badge
                    :class="roleBadgeClass(user.role)"
                    variant="outline"
                  >
                    {{ roleLabel(user.role) }}
                  </Badge>
                </TableCell>
                <TableCell class="text-right font-medium whitespace-nowrap">
                  {{ user.points }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  <span class="text-sm whitespace-nowrap">{{ vipLevelLabel(user.vipLevelId) }}</span>
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  <Badge
                    :class="statusBadgeClass(user.status)"
                    variant="outline"
                  >
                    {{ statusLabel(user.status) }}
                  </Badge>
                </TableCell>
                <TableCell class="text-right whitespace-nowrap">
                  <div class="flex items-center justify-end gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      class="h-8 w-8"
                      @click="openEdit(user)"
                    >
                      <Search class="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      class="h-8 w-8"
                      @click="openDelete(user)"
                    >
                      <X class="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <div
          v-if="!loading && totalPages > 1"
          class="flex items-center justify-between mt-4"
        >
          <div class="text-sm text-muted-foreground">
            第 {{ currentPage }} / {{ totalPages }} 页
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage <= 1"
              @click="goToPage(currentPage - 1)"
            >
              上一页
            </Button>
            <Button
              variant="outline"
              size="sm"
              :disabled="currentPage >= totalPages"
              @click="goToPage(currentPage + 1)"
            >
              下一页
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>

    <!-- Edit Role Dialog -->
    <Dialog v-model:open="editDialogOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>编辑用户角色</DialogTitle>
          <DialogDescription>
            修改用户 {{ editUser?.username }} 的角色
          </DialogDescription>
        </DialogHeader>
        <div class="space-y-4">
          <div>
            <Label>当前角色</Label>
            <Badge
              v-if="editUser"
              :class="roleBadgeClass(editUser.role)"
              variant="outline"
              class="ml-2"
            >
              {{ roleLabel(editUser.role) }}
            </Badge>
          </div>
          <div>
            <Label>新角色</Label>
            <Select v-model="editRole">
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="customer">普通用户</SelectItem>
                <SelectItem value="vip">VIP用户</SelectItem>
                <SelectItem value="hotel_manager">酒店管理员</SelectItem>
                <SelectItem value="admin">系统管理员</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" @click="editDialogOpen = false">取消</Button>
          <Button @click="confirmEdit">确认更新</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Confirm Dialog -->
    <Dialog v-model:open="deleteDialogOpen">
      <DialogContent class="max-w-sm">
        <DialogHeader>
          <DialogTitle>确认删除</DialogTitle>
          <DialogDescription>
            确定要删除用户 {{ deleteTargetUser?.username }} 吗？此操作不可撤销。
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="deleteDialogOpen = false">取消</Button>
          <Button variant="destructive" @click="confirmDelete">删除</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
