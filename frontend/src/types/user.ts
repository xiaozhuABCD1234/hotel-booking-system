export interface User {
  id: string
  username: string
  realName?: string
  phone?: string
  email?: string
  role: 'customer' | 'vip' | 'hotel_manager' | 'admin'
  vipLevelId?: number
  points: number
  createAt: string
  updateAt: string
  status: number
}

export interface UserVip {
  userId: string
  username: string
  phone?: string
  email?: string
  role: string
  vipLevelName?: string
  discountRate?: number
  points: number
}

export interface CreateUserRequest {
  username: string
  password: string
  phone?: string
  email?: string
  role?: 'customer' | 'vip' | 'hotel_manager' | 'admin'
}

export interface UpdateUserRequest {
  phone?: string
  email?: string
  role?: 'customer' | 'vip' | 'hotel_manager' | 'admin'
}
