import { defineStore } from 'pinia'
import { ref } from 'vue'
import { hotelApi, roomApi } from '@/api'
import type { Hotel, HotelSearchParams, Room, Pagination } from '@/types'

export const useHotelStore = defineStore('hotel', () => {
  const hotels = ref<Hotel[]>([])
  const pagination = ref<Pagination | null>(null)
  const loading = ref(false)
  const currentHotel = ref<Hotel | null>(null)
  const rooms = ref<Room[]>([])

  async function fetchHotels(params: HotelSearchParams = {}) {
    loading.value = true
    try {
      const res = await hotelApi.list({ page: 1, pageSize: 12, ...params })
      if (res.data.success && res.data.data) {
        hotels.value = res.data.data
        pagination.value = res.data.pagination ?? null
      } else {
        hotels.value = []
        pagination.value = null
      }
    } catch {
      hotels.value = []
      pagination.value = null
    } finally {
      loading.value = false
    }
  }

  async function fetchHotelById(id: string) {
    const res = await hotelApi.getById(id)
    if (res.data.success && res.data.data) {
      currentHotel.value = res.data.data
    }
    return res.data
  }

  async function fetchRooms(hotelId: string) {
    try {
      const res = await roomApi.list({ hotelId, pageSize: 100 })
      if (res.data.success && res.data.data) {
        rooms.value = res.data.data
      } else {
        rooms.value = []
      }
      return res.data
    } catch {
      rooms.value = []
      return undefined
    }
  }

  async function createHotel(data: Partial<Hotel>) {
    const res = await hotelApi.create(data)
    return res.data
  }

  async function updateHotel(id: string, data: Partial<Hotel>) {
    const res = await hotelApi.update(id, data)
    return res.data
  }

  async function deleteHotel(id: string) {
    const res = await hotelApi.delete(id)
    return res.data
  }

  async function createRoom(data: Partial<Room>) {
    const res = await roomApi.create(data)
    return res.data
  }

  async function updateRoom(id: string, data: Partial<Room>) {
    const res = await roomApi.update(id, data)
    return res.data
  }

  async function deleteRoom(id: string) {
    const res = await roomApi.delete(id)
    return res.data
  }

  return {
    hotels,
    pagination,
    loading,
    currentHotel,
    rooms,
    fetchHotels,
    fetchHotelById,
    fetchRooms,
    createHotel,
    updateHotel,
    deleteHotel,
    createRoom,
    updateRoom,
    deleteRoom,
  }
})
