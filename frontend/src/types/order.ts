import type { Room } from './room'
import type { User } from './user'

export type OrderStatus = 'pending' | 'booked' | 'checked_in' | 'cancelled' | 'completed'

export interface Order {
  id: string
  userId: string
  roomId: string
  quantity: number
  checkInDate: string
  checkOutDate: string
  totalPrice: number
  discount: number
  actualPrice: number
  status: OrderStatus
  createAt: string
  updateAt: string
  room?: Room
  user?: User
  guests?: OrderGuest[]
}

export interface OrderGuest {
  orderId: string
  idCard: string
  person?: {
    idCard: string
    name: string
    phone?: string
  }
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
