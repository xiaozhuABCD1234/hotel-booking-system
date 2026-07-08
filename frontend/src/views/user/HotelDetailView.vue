<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useHotelStore, useAuthStore } from '@/stores'
import type { Room } from '@/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'
import {
  ChevronLeft,
  MapPin,
  Phone,
  BedDouble,
} from '@lucide/vue'

const route = useRoute()
const router = useRouter()
const hotelStore = useHotelStore()
const authStore = useAuthStore()

const hotelId = computed(() => route.params.id as string)
const currentImageIndex = ref(0)
const roomsLoading = ref(false)

const hotel = computed(() => hotelStore.currentHotel)
const rooms = computed(() => hotelStore.rooms)
const images = computed(() => hotel.value?.images ?? [])

function renderStars(level: number | undefined): string {
  if (!level) return ''
  return '\u2605'.repeat(level) + '\u2606'.repeat(5 - level)
}

function prevImage() {
  if (images.value.length === 0) return
  currentImageIndex.value =
    (currentImageIndex.value - 1 + images.value.length) % images.value.length
}

function nextImage() {
  if (images.value.length === 0) return
  currentImageIndex.value = (currentImageIndex.value + 1) % images.value.length
}

function goToBooking(room: Room) {
  if (!authStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: route.fullPath } })
    return
  }
  router.push(`/booking/${room.id}`)
}

onMounted(async () => {
  try {
    await hotelStore.fetchHotelById(hotelId.value)
    roomsLoading.value = true
    await hotelStore.fetchRooms(hotelId.value)
  } finally {
    roomsLoading.value = false
  }
})
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <div class="container mx-auto max-w-5xl px-4 py-6">
      <!-- Back Button -->
      <Button
        variant="ghost"
        size="sm"
        class="mb-4 text-gray-600"
        @click="router.push('/')"
      >
        <ChevronLeft class="mr-1 h-4 w-4" />
        返回列表
      </Button>

      <!-- Loading State -->
      <div v-if="!hotel && !hotelStore.loading" class="py-20 text-center">
        <p class="text-gray-500">酒店不存在或已下架</p>
        <Button variant="outline" class="mt-4" @click="router.push('/')">
          返回首页
        </Button>
      </div>

      <div v-else-if="hotelStore.loading || !hotel">
        <Skeleton class="mb-4 h-80 w-full rounded-lg" />
        <Skeleton class="mb-2 h-8 w-1/2" />
        <Skeleton class="mb-4 h-5 w-1/3" />
        <Skeleton class="h-32 w-full" />
      </div>

      <!-- Hotel Content -->
      <div v-else-if="hotel">
        <!-- Image Gallery -->
        <div class="relative mb-6 overflow-hidden rounded-lg bg-gray-100">
          <div class="aspect-[16/9] w-full">
            <img
              v-if="images[currentImageIndex]?.imageUrl"
              :src="images[currentImageIndex].imageUrl"
              :alt="hotel.hotelName"
              class="h-full w-full object-cover"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center text-gray-400"
            >
              <MapPin class="h-16 w-16" />
            </div>
          </div>

          <!-- Image Navigation -->
          <template v-if="images.length > 1">
            <Button
              variant="outline"
              size="icon"
              class="absolute left-3 top-1/2 -translate-y-1/2 bg-white/80 shadow-sm"
              @click="prevImage"
            >
              <ChevronLeft class="h-5 w-5" />
            </Button>
            <Button
              variant="outline"
              size="icon"
              class="absolute right-3 top-1/2 -translate-y-1/2 bg-white/80 shadow-sm"
              @click="nextImage"
            >
              <ChevronLeft class="h-5 w-5 rotate-180" />
            </Button>
            <div class="absolute bottom-3 left-1/2 -translate-x-1/2 rounded-full bg-black/50 px-3 py-1 text-xs text-white">
              {{ currentImageIndex + 1 }} / {{ images.length }}
            </div>
          </template>
        </div>

        <!-- Thumbnail Strip -->
        <div v-if="images.length > 1" class="mb-6 flex gap-2 overflow-x-auto pb-2">
          <button
            v-for="(img, idx) in images"
            :key="idx"
            class="h-16 w-24 shrink-0 overflow-hidden rounded-md border-2 transition-colors"
            :class="idx === currentImageIndex ? 'border-blue-900' : 'border-transparent'"
            @click="currentImageIndex = idx"
          >
            <img :src="img.imageUrl" :alt="`${hotel.hotelName} ${idx + 1}`" class="h-full w-full object-cover" />
          </button>
        </div>

        <!-- Hotel Info -->
        <Card class="mb-6">
          <CardContent class="p-6">
            <div class="mb-3 flex items-start justify-between">
              <div>
                <h1 class="text-2xl font-semibold text-gray-900">
                  {{ hotel.hotelName }}
                </h1>
                <div
                  v-if="hotel.starLevel"
                  class="mt-1 text-amber-500"
                >
                  {{ renderStars(hotel.starLevel) }}
                </div>
              </div>
            </div>

            <div class="space-y-2 text-sm text-gray-600">
              <div class="flex items-center gap-2">
                <MapPin class="h-4 w-4 shrink-0 text-gray-400" />
                <span>{{ hotel.address }}</span>
              </div>
              <div class="flex items-center gap-2">
                <Phone class="h-4 w-4 shrink-0 text-gray-400" />
                <span>{{ hotel.telephone }}</span>
              </div>
            </div>

            <Separator class="my-4" />

            <div v-if="hotel.description">
              <h3 class="mb-2 text-sm font-medium text-gray-900">酒店简介</h3>
              <p class="text-sm leading-relaxed text-gray-600">
                {{ hotel.description }}
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- Room List -->
        <div>
          <h2 class="mb-4 text-xl font-semibold text-gray-900">
            客房类型
          </h2>

          <!-- Rooms Loading -->
          <div v-if="roomsLoading" class="space-y-4">
            <Card v-for="i in 3" :key="i">
              <CardContent class="flex items-center gap-4 p-4">
                <Skeleton class="h-20 w-28 rounded-md" />
                <div class="flex-1">
                  <Skeleton class="mb-2 h-5 w-1/3" />
                  <Skeleton class="mb-2 h-4 w-1/2" />
                  <Skeleton class="h-4 w-1/4" />
                </div>
              </CardContent>
            </Card>
          </div>

          <!-- Rooms Empty -->
          <div
            v-else-if="rooms.length === 0"
            class="flex flex-col items-center py-12 text-center"
          >
            <BedDouble class="mb-3 h-10 w-10 text-gray-300" />
            <p class="text-gray-500">暂无可用客房</p>
          </div>

          <!-- Room Cards -->
          <div v-else class="space-y-4">
            <Card
              v-for="room in rooms"
              :key="room.id"
              class="overflow-hidden"
            >
              <CardContent class="flex flex-col gap-4 p-4 sm:flex-row">
                <!-- Room Image -->
                <div class="h-32 w-full shrink-0 overflow-hidden rounded-md bg-gray-100 sm:h-28 sm:w-40">
                  <img
                    v-if="room.imageUrl"
                    :src="room.imageUrl"
                    :alt="room.typeName"
                    class="h-full w-full object-cover"
                  />
                  <div
                    v-else
                    class="flex h-full w-full items-center justify-center text-gray-300"
                  >
                    <BedDouble class="h-8 w-8" />
                  </div>
                </div>

                <!-- Room Info -->
                <div class="flex flex-1 flex-col justify-between">
                  <div>
                    <div class="mb-1 flex items-center gap-2">
                      <h3 class="text-base font-semibold text-gray-900">
                        {{ room.typeName }}
                      </h3>
                      <Badge
                        :variant="room.availableQuantity > 0 ? 'default' : 'secondary'"
                        :class="room.availableQuantity > 0 ? 'bg-green-100 text-green-700 hover:bg-green-100' : ''"
                      >
                        {{ room.availableQuantity > 0 ? `剩余 ${room.availableQuantity} 间` : '已满' }}
                      </Badge>
                    </div>
                    <p v-if="room.description" class="mb-2 text-sm text-gray-500">
                      {{ room.description }}
                    </p>
                  </div>

                  <div class="flex items-center justify-between">
                    <span class="text-xl font-bold text-blue-900">
                      ¥{{ room.price }}
                      <span class="text-xs font-normal text-gray-400">/晚</span>
                    </span>
                    <Button
                      size="sm"
                      class="bg-blue-900 hover:bg-blue-800"
                      :disabled="room.availableQuantity <= 0"
                      @click="goToBooking(room)"
                    >
                      预订
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
