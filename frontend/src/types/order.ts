import type { Room } from './room'
import type { Hotel } from './hotel'
import type { User } from './user'

export type OrderStatus = 'pending' | 'confirmed' | 'checked_in' | 'cancelled' | 'completed'

export interface Order {
  id: string
  userId: string
  roomId: string
  hotelId: string
  checkInDate: string
  checkOutDate: string
  guestName: string
  guestPhone: string
  guestIdCard: string
  roomCount: number
  totalPrice: number
  status: OrderStatus
  createAt: string
  updateAt: string
  room?: Room
  hotel?: Hotel
  user?: User
}

export interface OrderFull {
  orderId: string
  userId: string
  username: string
  hotelName: string
  roomType: string
  checkInDate: string
  checkOutDate: string
  guestName: string
  guestPhone: string
  guestIdCard: string
  roomCount: number
  totalPrice: number
  status: OrderStatus
  createAt: string
}

export interface MyOrder {
  orderId: string
  hotelName: string
  roomType: string
  price: number
  checkInDate: string
  checkOutDate: string
  guestName: string
  guestPhone: string
  guestIdCard: string
  roomCount: number
  totalPrice: number
  status: OrderStatus
  createAt: string
}

export interface CreateOrderRequest {
  roomId: string
  checkInDate: string
  checkOutDate: string
  guestName: string
  guestPhone: string
  guestIdCard: string
  roomCount: number
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus
}
