<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useHotelStore } from '@/stores'
import { regionApi } from '@/api'
import type { HotelSearchParams, Region } from '@/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import { Search, MapPin } from '@lucide/vue'

const router = useRouter()
const hotelStore = useHotelStore()

const regions = ref<Region[]>([])
const keyword = ref('')
const selectedRegion = ref<string>('')
const selectedStar = ref<string>('')

const searchParams = computed<HotelSearchParams>(() => {
  const params: HotelSearchParams = {}
  if (selectedRegion.value && selectedRegion.value !== 'all') {
    params.regionId = Number(selectedRegion.value)
  }
  if (selectedStar.value && selectedStar.value !== 'all') {
    params.starLevel = Number(selectedStar.value)
  }
  if (keyword.value.trim()) params.keyword = keyword.value.trim()
  return params
})

function handleSearch() {
  hotelStore.fetchHotels(searchParams.value)
}

function goToHotel(id: string) {
  router.push(`/hotel/${id}`)
}

function renderStars(level: number | undefined): string {
  if (!level) return ''
  return '\u2605'.repeat(level) + '\u2606'.repeat(5 - level)
}

onMounted(async () => {
  hotelStore.fetchHotels()
  try {
    const res = await regionApi.list()
    if (res.data.success && res.data.data) {
      regions.value = res.data.data
    }
  } catch {
    // regions load failure is non-critical
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Hero / Search Section -->
    <section class="bg-blue-900 px-4 py-12">
      <div class="container mx-auto max-w-4xl">
        <h1 class="mb-2 text-center text-3xl font-semibold text-white">
          发现理想酒店
        </h1>
        <p class="mb-8 text-center text-blue-200">
          搜索全国优质酒店，轻松预订
        </p>

        <Card class="border-0 shadow-lg">
          <CardContent class="p-6">
            <div class="grid gap-4 md:grid-cols-4">
              <!-- Keyword -->
              <div class="md:col-span-1">
                <Label class="mb-1.5 block text-sm">关键词</Label>
                <div class="relative">
                  <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                  <Input
                    v-model="keyword"
                    placeholder="酒店名称..."
                    class="pl-9"
                    @keyup.enter="handleSearch"
                  />
                </div>
              </div>

              <!-- Region -->
              <div>
                <Label class="mb-1.5 block text-sm">地区</Label>
                <Select v-model="selectedRegion">
                  <SelectTrigger>
                    <SelectValue placeholder="全部地区" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部地区</SelectItem>
                    <SelectItem
                      v-for="region in regions"
                      :key="region.id"
                      :value="String(region.id)"
                    >
                      {{ region.regionName }}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <!-- Star Level -->
              <div>
                <Label class="mb-1.5 block text-sm">星级</Label>
                <Select v-model="selectedStar">
                  <SelectTrigger>
                    <SelectValue placeholder="全部星级" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部星级</SelectItem>
                    <SelectItem value="5">五星级</SelectItem>
                    <SelectItem value="4">四星级</SelectItem>
                    <SelectItem value="3">三星级</SelectItem>
                    <SelectItem value="2">二星级</SelectItem>
                    <SelectItem value="1">一星级</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <!-- Search Button -->
              <div class="flex items-end">
                <Button
                  class="w-full bg-blue-900 hover:bg-blue-800"
                  @click="handleSearch"
                >
                  <Search class="mr-2 h-4 w-4" />
                  搜索
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </section>

    <!-- Results Section -->
    <section class="px-4 py-8">
      <div class="container mx-auto max-w-6xl">
        <!-- Loading State -->
        <div v-if="hotelStore.loading" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card v-for="i in 6" :key="i">
            <Skeleton class="h-48 w-full rounded-t-lg" />
            <CardContent class="p-4">
              <Skeleton class="mb-2 h-5 w-3/4" />
              <Skeleton class="mb-2 h-4 w-1/2" />
              <Skeleton class="mb-3 h-4 w-full" />
              <div class="flex justify-between">
                <Skeleton class="h-5 w-20" />
                <Skeleton class="h-5 w-16" />
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Empty State -->
        <div
          v-else-if="hotelStore.hotels.length === 0"
          class="flex flex-col items-center justify-center py-20 text-center"
        >
          <Search class="mb-4 h-12 w-12 text-gray-300" />
          <h3 class="mb-1 text-lg font-medium text-gray-700">暂无搜索结果</h3>
          <p class="text-sm text-gray-500">请尝试调整筛选条件</p>
        </div>

        <!-- Hotel Cards -->
        <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Card
            v-for="hotel in hotelStore.hotels"
            :key="hotel.id"
            class="cursor-pointer overflow-hidden transition-shadow hover:shadow-md"
            @click="goToHotel(hotel.id)"
          >
            <!-- Image -->
            <div class="relative h-48 w-full overflow-hidden bg-gray-100">
              <img
                v-if="hotel.images?.[0]?.imageUrl"
                :src="hotel.images[0].imageUrl"
                :alt="hotel.hotelName"
                class="h-full w-full object-cover"
              />
              <div
                v-else
                class="flex h-full w-full items-center justify-center text-gray-400"
              >
                <MapPin class="h-8 w-8" />
              </div>
              <Badge
                v-if="hotel.starLevel"
                class="absolute right-2 top-2 bg-white/90 text-amber-600 hover:bg-white/90"
              >
                {{ renderStars(hotel.starLevel) }}
              </Badge>
            </div>

            <CardContent class="p-4">
              <h3 class="mb-1 truncate text-lg font-semibold text-gray-900">
                {{ hotel.hotelName }}
              </h3>
              <div class="mb-2 flex items-center gap-1 text-sm text-gray-500">
                <MapPin class="h-3.5 w-3.5 shrink-0" />
                <span class="truncate">{{ hotel.address }}</span>
              </div>
              <p
                v-if="hotel.description"
                class="mb-3 line-clamp-2 text-sm text-gray-500"
              >
                {{ hotel.description }}
              </p>
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-500">
                  {{ hotel.telephone }}
                </span>
                <span class="text-xs text-gray-400">
                  {{ hotel.images?.length ?? 0 }} 张图片
                </span>
              </div>
            </CardContent>
          </Card>
        </div>

        <!-- Pagination -->
        <div
          v-if="hotelStore.pagination && hotelStore.pagination.totalPages > 1"
          class="mt-8 flex items-center justify-center gap-2"
        >
          <Button
            variant="outline"
            size="sm"
            :disabled="!hotelStore.pagination.hasPrev"
            @click="hotelStore.fetchHotels({ ...searchParams, page: (hotelStore.pagination?.currentPage ?? 1) - 1 })"
          >
            上一页
          </Button>
          <span class="text-sm text-gray-500">
            {{ hotelStore.pagination.currentPage }} / {{ hotelStore.pagination.totalPages }}
          </span>
          <Button
            variant="outline"
            size="sm"
            :disabled="!hotelStore.pagination.hasNext"
            @click="hotelStore.fetchHotels({ ...searchParams, page: (hotelStore.pagination?.currentPage ?? 1) + 1 })"
          >
            下一页
          </Button>
        </div>
      </div>
    </section>
  </div>
</template>
