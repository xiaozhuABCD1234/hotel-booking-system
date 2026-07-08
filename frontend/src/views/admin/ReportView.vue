<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { reportApi } from '@/api/region'
import type { GuestBookingStats } from '@/types'
import { toast } from 'vue-sonner'
import { getApiErrorMessage } from '@/lib/utils'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { BarChart3, TrendingUp, Users } from '@lucide/vue'

interface HotelSummaryRow {
  hotelId: string
  hotelName: string
  city?: string
  starLevel?: number
  avgRating: number
  minPrice: number
  reviewCount: number
  roomCount: number
}

interface UserVipRow {
  userId: string
  username: string
  vipLevelName?: string
  discountRate?: number
  points: number
}

const activeTab = ref('hotel')

const hotelData = ref<HotelSummaryRow[]>([])
const hotelLoading = ref(false)
const hotelSortByRating = ref(false)

const vipData = ref<UserVipRow[]>([])
const vipLoading = ref(false)

const guestStats = ref<GuestBookingStats[]>([])
const topGuests = ref<GuestBookingStats[]>([])
const guestLoading = ref(false)

const sortedHotelData = () => {
  if (!hotelSortByRating.value) return hotelData.value
  return [...hotelData.value].sort((a, b) => b.avgRating - a.avgRating)
}

const vipLevelBadgeClass = (name?: string): string => {
  if (!name) return 'bg-gray-100 text-gray-800 border-gray-200'
  const lower = name.toLowerCase()
  if (lower.includes('gold') || lower.includes('金'))
    return 'bg-yellow-100 text-yellow-800 border-yellow-200'
  if (lower.includes('silver') || lower.includes('银'))
    return 'bg-gray-200 text-gray-700 border-gray-300'
  if (lower.includes('bronze') || lower.includes('铜'))
    return 'bg-orange-100 text-orange-800 border-orange-200'
  return 'bg-gray-100 text-gray-800 border-gray-200'
}

async function fetchHotelSummaries() {
  hotelLoading.value = true
  try {
    const res = await reportApi.hotelSummaries()
    if (res.data.data) {
      hotelData.value = res.data.data as HotelSummaryRow[]
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '获取酒店统计失败'))
  } finally {
    hotelLoading.value = false
  }
}

async function fetchUserVip() {
  vipLoading.value = true
  try {
    const res = await reportApi.userVipList()
    if (res.data.data) {
      vipData.value = res.data.data as UserVipRow[]
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '获取用户VIP数据失败'))
  } finally {
    vipLoading.value = false
  }
}

async function fetchGuestData() {
  guestLoading.value = true
  try {
    const [statsRes, topRes] = await Promise.all([
      reportApi.guestStats(),
      reportApi.topGuests(20),
    ])
    if (statsRes.data.data) {
      guestStats.value = statsRes.data.data as GuestBookingStats[]
    }
    if (topRes.data.data) {
      topGuests.value = topRes.data.data
    }
  } catch (e: unknown) {
    toast.error(getApiErrorMessage(e, '获取入住人分析数据失败'))
  } finally {
    guestLoading.value = false
  }
}

function handleTabChange(tab: string | number) {
  const t = String(tab)
  activeTab.value = t
  if (t === 'hotel' && hotelData.value.length === 0) fetchHotelSummaries()
  if (t === 'vip' && vipData.value.length === 0) fetchUserVip()
  if (t === 'guest' && topGuests.value.length === 0) fetchGuestData()
}

onMounted(() => {
  fetchHotelSummaries()
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold flex items-center gap-2">
        <BarChart3 class="h-6 w-6" />
        数据报表
      </h1>
    </div>

    <Tabs default-value="hotel" @update:model-value="handleTabChange">
      <TabsList>
        <TabsTrigger value="hotel">酒店统计</TabsTrigger>
        <TabsTrigger value="vip">用户VIP</TabsTrigger>
        <TabsTrigger value="guest">入住人分析</TabsTrigger>
      </TabsList>

      <!-- Tab 1: Hotel Statistics -->
      <TabsContent value="hotel" class="mt-4">
        <Card>
          <CardHeader class="pb-4">
            <div class="flex items-center justify-between">
              <CardTitle class="text-lg flex items-center gap-2">
                <BarChart3 class="h-5 w-5" />
                酒店统计概览
              </CardTitle>
              <button
                class="text-sm text-muted-foreground hover:text-foreground transition-colors flex items-center gap-1"
                @click="hotelSortByRating = !hotelSortByRating"
              >
                <TrendingUp class="h-4 w-4" />
                {{ hotelSortByRating ? '取消排序' : '按评分排序' }}
              </button>
            </div>
          </CardHeader>
          <CardContent>
            <div v-if="hotelLoading" class="space-y-3">
              <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
            </div>
            <div v-else class="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>酒店名称</TableHead>
                    <TableHead>城市</TableHead>
                    <TableHead class="text-center">星级</TableHead>
                    <TableHead class="text-center">平均评分</TableHead>
                    <TableHead class="text-right">最低价</TableHead>
                    <TableHead class="text-center">评价数</TableHead>
                    <TableHead class="text-center">房间数</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-if="sortedHotelData().length === 0">
                    <TableCell colspan="7" class="text-center py-8 text-muted-foreground">
                      暂无数据
                    </TableCell>
                  </TableRow>
                  <TableRow
                    v-for="hotel in sortedHotelData()"
                    :key="hotel.hotelId"
                    class="hover:bg-muted/50"
                  >
                    <TableCell class="font-medium">{{ hotel.hotelName }}</TableCell>
                    <TableCell>{{ hotel.city ?? '-' }}</TableCell>
                    <TableCell class="text-center">
                      <span v-if="hotel.starLevel">
                        {{ '★'.repeat(hotel.starLevel) }}
                      </span>
                      <span v-else>-</span>
                    </TableCell>
                    <TableCell class="text-center">
                      <Badge v-if="hotel.avgRating > 0" variant="secondary">
                        {{ hotel.avgRating.toFixed(1) }}
                      </Badge>
                      <span v-else>-</span>
                    </TableCell>
                    <TableCell class="text-right font-medium">
                      ¥{{ hotel.minPrice.toFixed(0) }}
                    </TableCell>
                    <TableCell class="text-center">{{ hotel.reviewCount }}</TableCell>
                    <TableCell class="text-center">{{ hotel.roomCount }}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <!-- Tab 2: User VIP -->
      <TabsContent value="vip" class="mt-4">
        <Card>
          <CardHeader class="pb-4">
            <CardTitle class="text-lg flex items-center gap-2">
              <Users class="h-5 w-5" />
              用户VIP等级分布
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div v-if="vipLoading" class="space-y-3">
              <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
            </div>
            <div v-else class="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>用户名</TableHead>
                    <TableHead>VIP等级</TableHead>
                    <TableHead class="text-center">折扣率</TableHead>
                    <TableHead class="text-right">积分</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-if="vipData.length === 0">
                    <TableCell colspan="6" class="text-center py-8 text-muted-foreground">
                      暂无数据
                    </TableCell>
                  </TableRow>
                  <TableRow
                    v-for="user in vipData"
                    :key="user.userId"
                    class="hover:bg-muted/50"
                  >
                    <TableCell class="font-medium">{{ user.username }}</TableCell>
                    <TableCell>
                      <Badge
                        v-if="user.vipLevelName"
                        :class="vipLevelBadgeClass(user.vipLevelName)"
                        variant="outline"
                      >
                        {{ user.vipLevelName }}
                      </Badge>
                      <span v-else class="text-muted-foreground">-</span>
                    </TableCell>
                    <TableCell class="text-center">
                      {{ user.discountRate ? `${(user.discountRate * 100).toFixed(0)}%` : '-' }}
                    </TableCell>
                    <TableCell class="text-right font-medium">{{ user.points }}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <!-- Tab 3: Guest Analysis -->
      <TabsContent value="guest" class="mt-4">
        <div class="space-y-6">
          <!-- Stat Cards -->
          <div v-if="guestStats.length > 0" class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium text-muted-foreground">
                  入住人群分组数
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">{{ guestStats.length }}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium text-muted-foreground">
                  总预订次数
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  {{ guestStats.reduce((sum, g) => sum + g.totalOrders, 0) }}
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader class="pb-2">
                <CardTitle class="text-sm font-medium text-muted-foreground">
                  总消费金额
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="text-2xl font-bold">
                  ¥{{ guestStats.reduce((sum, g) => sum + g.totalAmount, 0).toFixed(0) }}
                </div>
              </CardContent>
            </Card>
          </div>

          <!-- Top Guests Table -->
          <Card>
            <CardHeader class="pb-4">
              <CardTitle class="text-lg flex items-center gap-2">
                <TrendingUp class="h-5 w-5" />
                热门入住人排行
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div v-if="guestLoading" class="space-y-3">
                <Skeleton v-for="i in 5" :key="i" class="h-12 w-full" />
              </div>
              <div v-else class="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead class="w-[50px]">#</TableHead>
                      <TableHead>姓名</TableHead>
                      <TableHead>性别</TableHead>
                      <TableHead class="text-center">年龄</TableHead>
                      <TableHead class="text-center">预订次数</TableHead>
                      <TableHead class="text-right">累计消费</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    <TableRow v-if="topGuests.length === 0">
                    <TableCell colspan="4" class="text-center py-8 text-muted-foreground">
                        暂无数据
                      </TableCell>
                    </TableRow>
                    <TableRow
                      v-for="(guest, index) in topGuests"
                      :key="guest.personIdCard"
                      class="hover:bg-muted/50"
                    >
                      <TableCell class="font-medium text-muted-foreground">
                        {{ index + 1 }}
                      </TableCell>
                      <TableCell class="font-medium">{{ guest.personName }}</TableCell>
                      <TableCell>{{ guest.gender }}</TableCell>
                      <TableCell class="text-center">{{ guest.age }}</TableCell>
                      <TableCell class="text-center font-medium">
                        {{ guest.totalOrders }}
                      </TableCell>
                      <TableCell class="text-right font-medium">
                        ¥{{ guest.totalAmount.toFixed(2) }}
                      </TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </div>
            </CardContent>
          </Card>
        </div>
      </TabsContent>
    </Tabs>
  </div>
</template>
