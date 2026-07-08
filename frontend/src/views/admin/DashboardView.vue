<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
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
import { Hotel, DoorOpen, ShoppingBag, Users } from '@lucide/vue'
import { hotelApi, roomApi, orderApi, userApi } from '@/api'
import type { Order, OrderStatus } from '@/types'

const loading = ref(true)
const hotelCount = ref(0)
const roomCount = ref(0)
const orderCount = ref(0)
const userCount = ref(0)
const recentOrders = ref<Order[]>([])

const statusLabels: Record<OrderStatus, string> = {
  pending: '待确认',
  booked: '已预订',
  checked_in: '已入住',
  cancelled: '已取消',
  completed: '已完成',
}

const statusVariants: Record<OrderStatus, 'default' | 'secondary' | 'destructive' | 'outline'> = {
  pending: 'secondary',
  booked: 'default',
  checked_in: 'default',
  cancelled: 'destructive',
  completed: 'outline',
}

onMounted(async () => {
  try {
    const [hotelsRes, roomsRes, ordersRes, usersRes] = await Promise.all([
      hotelApi.list({ page: 1, pageSize: 1 }),
      roomApi.list({ page: 1, pageSize: 1 }),
      orderApi.list(1, 5),
      userApi.list(1, 1),
    ])

    hotelCount.value = hotelsRes.data.pagination?.totalItems ?? 0
    roomCount.value = roomsRes.data.pagination?.totalItems ?? 0
    orderCount.value = ordersRes.data.pagination?.totalItems ?? 0
    userCount.value = usersRes.data.pagination?.totalItems ?? 0
    recentOrders.value = ordersRes.data.data ?? []
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  } finally {
    loading.value = false
  }
})

const stats = [
  { title: '酒店数量', value: () => hotelCount.value, icon: Hotel, color: 'text-blue-600' },
  { title: '客房数量', value: () => roomCount.value, icon: DoorOpen, color: 'text-green-600' },
  { title: '订单数量', value: () => orderCount.value, icon: ShoppingBag, color: 'text-orange-600' },
  { title: '用户数量', value: () => userCount.value, icon: Users, color: 'text-purple-600' },
]
</script>

<template>
  <div class="space-y-6">
    <h1 class="text-2xl font-semibold">仪表盘</h1>

    <!-- Stat Cards -->
    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card v-for="stat in stats" :key="stat.title" class="shadow-sm">
        <CardHeader class="flex flex-row items-center justify-between pb-2">
          <CardTitle class="text-sm font-medium text-muted-foreground">
            {{ stat.title }}
          </CardTitle>
          <component :is="stat.icon" :class="['h-5 w-5', stat.color]" />
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-20" />
          <div v-else class="text-2xl font-bold">{{ stat.value() }}</div>
        </CardContent>
      </Card>
    </div>

    <!-- Recent Orders -->
    <Card class="shadow-sm">
      <CardHeader>
        <CardTitle>最近订单</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>客人姓名</TableHead>
              <TableHead>入住日期</TableHead>
              <TableHead>离店日期</TableHead>
              <TableHead>总价</TableHead>
              <TableHead>状态</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <template v-if="loading">
              <TableRow v-for="i in 5" :key="i">
                <TableCell v-for="j in 5" :key="j">
                  <Skeleton class="h-4 w-full" />
                </TableCell>
              </TableRow>
            </template>
            <template v-else-if="recentOrders.length === 0">
              <TableRow>
                <TableCell :colspan="5" class="text-center text-muted-foreground">
                  暂无数据
                </TableCell>
              </TableRow>
            </template>
            <template v-else>
              <TableRow v-for="order in recentOrders" :key="order.id">
                <TableCell>{{ order.guestName }}</TableCell>
                <TableCell>{{ order.checkInDate }}</TableCell>
                <TableCell>{{ order.checkOutDate }}</TableCell>
                <TableCell>¥{{ order.totalPrice.toFixed(2) }}</TableCell>
                <TableCell>
                  <Badge :variant="statusVariants[order.status]">
                    {{ statusLabels[order.status] }}
                  </Badge>
                </TableCell>
              </TableRow>
            </template>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
