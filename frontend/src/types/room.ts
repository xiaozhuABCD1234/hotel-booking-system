export interface Room {
  id: string
  hotelId: string
  roomType: string
  price: number
  totalCount: number
  availableCount: number
  description?: string
  imageUrl?: string
  createAt: string
  updateAt: string
  status: number
  hotel?: Hotel
}

export interface RoomDetails {
  roomId: string
  hotelId: string
  hotelName: string
  roomType: string
  price: number
  totalCount: number
  availableCount: number
  description?: string
  imageUrl?: string
  province?: string
  city?: string
  district?: string
  hotelStatus: number
  roomStatus: number
}

export interface RoomSearchParams {
  page?: number
  pageSize?: number
  hotelId?: string
}
