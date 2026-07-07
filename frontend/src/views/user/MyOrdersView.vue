<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore, useOrderStore } from '@/stores'
import type { Order, OrderStatus } from '@/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
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
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { toast } from 'vue-sonner'
import { ShoppingBag, ArrowLeft, AlertCircle } from '@lucide/vue'

const router = useRouter()
const authStore = useAuthStore()
const orderStore = useOrderStore()

const cancellingId = ref<string | null>(null)
const dialogOrderId = ref<string | null>(null)

const orders = computed(() => orderStore.orders)
const loading = computed(() => orderStore.loading)

const statusConfig: Record<OrderStatus, { label: string; class: string }> = {
  pending: { label: '待确认', class: 'bg-yellow-100 text-yellow-700 hover:bg-yellow-100' },
  confirmed: { label: '已确认', class: 'bg-green-100 text-green-700 hover:bg-green-100' },
  checked_in: { label: '已入住', class: 'bg-blue-100 text-blue-700 hover:bg-blue-100' },
  cancelled: { label: '已取消', class: 'bg-red-100 text-red-700 hover:bg-red-100' },
  completed: { label: '已完成', class: 'bg-gray-100 text-gray-700 hover:bg-gray-100' },
}

function getStatusBadge(status: OrderStatus) {
  return statusConfig[status] ?? { label: status, class: 'bg-gray-100 text-gray-600 hover:bg-gray-100' }
}

function canCancel(order: Order): boolean {
  return order.status === 'pending' || order.status === 'confirmed'
}

function openCancelDialog(orderId: string) {
  dialogOrderId.value = orderId
}

async function confirmCancel() {
  if (!dialogOrderId.value) return
  cancellingId.value = dialogOrderId.value
  try {
    const res = await orderStore.cancelOrder(dialogOrderId.value)
    if (res.success) {
      toast.success('订单已取消')
      refreshOrders()
    } else {
      toast.error(res.message || '取消失败')
    }
  } catch {
    toast.error('操作失败，请重试')
  } finally {
    cancellingId.value = null
    dialogOrderId.value = null
  }
}

function refreshOrders() {
  if (authStore.user) {
    orderStore.fetchMyOrders(authStore.user.id)
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '—'
  return dateStr.substring(0, 10)
}

onMounted(() => {
  if (!authStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: '/orders' } })
    return
  }
  if (authStore.user) {
    orderStore.fetchMyOrders(authStore.user.id)
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="container mx-auto max-w-5xl px-4 py-6">
      <h1 class="mb-6 text-2xl font-semibold text-gray-900">
        我的订单
      </h1>

      <!-- Loading -->
      <div v-if="loading" class="space-y-4">
        <Card v-for="i in 4" :key="i">
          <CardContent class="flex items-center gap-4 p-4">
            <Skeleton class="h-6 w-20" />
            <Skeleton class="h-5 w-32" />
            <Skeleton class="h-5 flex-1" />
            <Skeleton class="h-5 w-24" />
            <Skeleton class="h-8 w-16" />
          </CardContent>
        </Card>
      </div>

      <!-- Empty -->
      <div
        v-else-if="orders.length === 0"
        class="flex flex-col items-center justify-center py-20 text-center"
      >
        <ShoppingBag class="mb-4 h-12 w-12 text-gray-300" />
        <h3 class="mb-1 text-lg font-medium text-gray-700">暂无订单</h3>
        <p class="mb-4 text-sm text-gray-500">快去预订心仪的酒店吧</p>
        <Button
          class="bg-blue-900 hover:bg-blue-800"
          @click="router.push('/')"
        >
          <ArrowLeft class="mr-2 h-4 w-4" />
          浏览酒店
        </Button>
      </div>

      <!-- Orders Table -->
      <Card v-else>
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>酒店</TableHead>
                  <TableHead>房型</TableHead>
                  <TableHead>入住 / 退房</TableHead>
                  <TableHead>住客</TableHead>
                  <TableHead>间数</TableHead>
                  <TableHead>总价</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead class="text-right">操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="order in orders" :key="order.id">
                  <TableCell class="font-medium">
                    {{ order.hotel?.hotelName ?? '—' }}
                  </TableCell>
                  <TableCell>
                    {{ order.room?.roomType ?? '—' }}
                  </TableCell>
                  <TableCell>
                    <div class="text-sm">
                      <div>{{ formatDate(order.checkInDate) }}</div>
                      <div class="text-gray-400">{{ formatDate(order.checkOutDate) }}</div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <div class="text-sm">
                      <div>{{ order.guestName }}</div>
                      <div class="text-gray-400">{{ order.guestPhone }}</div>
                    </div>
                  </TableCell>
                  <TableCell>{{ order.roomCount }}</TableCell>
                  <TableCell class="font-medium">
                    ¥{{ order.totalPrice.toFixed(2) }}
                  </TableCell>
                  <TableCell>
                    <Badge
                      :class="getStatusBadge(order.status).class"
                    >
                      {{ getStatusBadge(order.status).label }}
                    </Badge>
                  </TableCell>
                  <TableCell class="text-right">
                    <Button
                      v-if="canCancel(order)"
                      variant="outline"
                      size="sm"
                      class="text-red-600 hover:bg-red-50 hover:text-red-700"
                      :disabled="cancellingId === order.id"
                      @click="openCancelDialog(order.id)"
                    >
                      取消
                    </Button>
                    <span v-else class="text-xs text-gray-400">—</span>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <!-- Pagination -->
      <div
        v-if="orderStore.pagination && orderStore.pagination.totalPages > 1"
        class="mt-6 flex items-center justify-center gap-2"
      >
        <Button
          variant="outline"
          size="sm"
          :disabled="!orderStore.pagination.hasPrev"
          @click="orderStore.fetchMyOrders(authStore.user!.id, (orderStore.pagination?.currentPage ?? 1) - 1)"
        >
          上一页
        </Button>
        <span class="text-sm text-gray-500">
          {{ orderStore.pagination.currentPage }} / {{ orderStore.pagination.totalPages }}
        </span>
        <Button
          variant="outline"
          size="sm"
          :disabled="!orderStore.pagination.hasNext"
          @click="orderStore.fetchMyOrders(authStore.user!.id, (orderStore.pagination?.currentPage ?? 1) + 1)"
        >
          下一页
        </Button>
      </div>
    </div>

    <!-- Cancel Confirmation Dialog -->
    <Dialog :open="dialogOrderId !== null" @update:open="(v) => { if (!v) dialogOrderId = null }">
      <DialogContent>
        <DialogHeader>
          <DialogTitle class="flex items-center gap-2">
            <AlertCircle class="h-5 w-5 text-red-500" />
            确认取消订单
          </DialogTitle>
          <DialogDescription>
            取消后无法恢复，确定要取消该订单吗？
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="dialogOrderId = null">
            返回
          </Button>
          <Button
            variant="destructive"
            :disabled="cancellingId !== null"
            @click="confirmCancel"
          >
            {{ cancellingId ? '取消中...' : '确认取消' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
