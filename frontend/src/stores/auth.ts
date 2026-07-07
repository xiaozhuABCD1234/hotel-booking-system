import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api'
import type { UserInfo, LoginRequest, RegisterRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))

  const isLoggedIn = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('accessToken', access)
    localStorage.setItem('refreshToken', refresh)
  }

  async function login(data: LoginRequest) {
    const res = await authApi.login(data)
    if (res.data.success && res.data.data) {
      const { tokens, user: userInfo } = res.data.data
      setTokens(tokens.accessToken, tokens.refreshToken)
      user.value = userInfo
    }
    return res.data
  }

  async function register(data: RegisterRequest) {
    const res = await authApi.register(data)
    if (res.data.success && res.data.data) {
      const { tokens, user: userInfo } = res.data.data
      setTokens(tokens.accessToken, tokens.refreshToken)
      user.value = userInfo
    }
    return res.data
  }

  async function logout() {
    try {
      await authApi.logout()
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
