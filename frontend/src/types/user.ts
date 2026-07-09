export interface VipLevel {
  level: number
  levelName: string
  minPoints: number
  discountRate: number
}

export interface User {
  id: string
  username: string
  realName?: string
  phone?: string
  email?: string
  role: 'customer' | 'vip' | 'hotel_manager' | 'admin'
  vipLevelId?: number
  vipLevel?: VipLevel
  points: number
  idCard?: string
  occupation?: string
  education?: string
  income?: number
  createAt: string
  updateAt: string
  status: number
}

export interface UserVip {
  userId: string
  username: string
  phone?: string
  email?: string
  realName?: string
  idCard?: string
  role: string
  vipLevelName?: string
  discountRate?: number
  points: number
  pointsToNextLevel?: number
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
  realName?: string
  idCard?: string
  occupation?: string
  education?: string
  income?: number
  oldPassword?: string
  password?: string
  role?: 'customer' | 'vip' | 'hotel_manager' | 'admin'
}
