/** Login request */
export interface LoginRequest {
  username: string
  password: string
}

/** Register request */
export interface RegisterRequest {
  username: string
  password: string
  phone?: string
  email?: string
}

/** Auth tokens */
export interface TokenPair {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

/** Login response data */
export interface LoginResponse {
  tokens: TokenPair
  user: UserInfo
}

export interface UserInfo {
  id: string
  username: string
  phone?: string
  email?: string
  role: 'user' | 'admin'
  vipLevelId?: number
  points: number
}
