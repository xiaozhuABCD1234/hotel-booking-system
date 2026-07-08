import type { User } from './user'
import type { Hotel } from './hotel'

export interface Review {
  id: string
  userId: string
  orderId: string
  hotelId: string
  rating: number
  content?: string
  createAt: string
  updateAt: string
  user?: User
  hotel?: Hotel
}

export interface ReviewFull {
  reviewId: string
  userId: string
  username: string
  hotelName: string
  rating: number
  content?: string
  createAt: string
}

export interface CreateReviewRequest {
  orderId: string
  hotelId: string
  rating: number
  content?: string
  images?: string[]
}

export interface UpdateReviewRequest {
  rating?: number
  content?: string
  images?: string[]
}
