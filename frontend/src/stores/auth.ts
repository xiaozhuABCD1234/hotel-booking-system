import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, userApi } from '@/api'
import type { UserInfo, LoginRequest, RegisterRequest } from '@/types'

/** Decode JWT payload (synchronous, no signature verification). */
function decodeToken(token: string): { userId: string; role: string } | null {
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return null
    const payload = JSON.parse(
      atob(parts[1].replace(/-/g, '+').replace(/_/g, '/')),
    )
    return { userId: payload.user_id, role: payload.role }
  } catch {
    return null
  }
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))

  const isLoggedIn = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  /** Set basic user info from JWT synchronously, then fetch full profile async. */
  async function initUserFromToken(token: string) {
    const decoded = decodeToken(token)
    if (!decoded) return

    console.log('[Auth] 用户角色:', decoded.role)

    // Set basic info from JWT immediately so isAdmin works synchronously
    user.value = {
      id: decoded.userId,
      username: '',
      role: decoded.role as 'user' | 'admin',
      points: 0,
    }

    // Fetch full user profile in background
    try {
      const res = await userApi.getById(decoded.userId)
      if (res.data.success && res.data.data) {
        const u = res.data.data
        user.value = {
          id: u.id,
          username: u.username,
          phone: u.phone,
          email: u.email,
          role: u.role as 'user' | 'admin',
          vipLevelId: u.vipLevelId,
          points: u.points,
        }
      }
    } catch {
      // Silently ignore — token may be expired; API interceptor handles 401 redirect
    }
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('accessToken', access)
    localStorage.setItem('refreshToken', refresh)
  }

  async function login(data: LoginRequest) {
    const res = await authApi.login(data)
    if (res.data.success && res.data.data) {
      const { accessToken: at, refreshToken: rt } = res.data.data
      setTokens(at, rt)
      // Set user info immediately so route guard can check role
      initUserFromToken(at)
    }
    return res.data
  }

  async function register(data: RegisterRequest) {
    const res = await authApi.register(data)
    if (res.data.success && res.data.data) {
      const { accessToken: at, refreshToken: rt } = res.data.data
      setTokens(at, rt)
      initUserFromToken(at)
    }
    return res.data
  }

  async function logout() {
    try {
      await authApi.logout(accessToken.value!, refreshToken.value!)
    } catch {
      // Ignore API errors — state cleanup happens in finally
    } finally {
      user.value = null
      accessToken.value = null
      refreshToken.value = null
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')
    }
  }

  function setUser(userInfo: UserInfo) {
    user.value = userInfo
  }

  // Restore user from stored token on page refresh
  if (accessToken.value) {
    initUserFromToken(accessToken.value)
  }

  return {
    user,
    accessToken,
    refreshToken,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
    setUser,
  }
})
