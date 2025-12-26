import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import router from '@/router'
import type { User, LoginResponse } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')
  const user = ref<User | null>(null)

  // 初始化时从 localStorage 加载用户数据
  try {
    const userData = localStorage.getItem('user')
    if (userData && userData !== 'undefined' && userData !== 'null') {
      user.value = JSON.parse(userData)
    }
  } catch (error) {
    console.warn('Failed to parse user data from localStorage:', error)
    user.value = null
  }

  const isAuthenticated = computed(() => !!token.value)

  /**
   * 设置认证信息
   * @param data 登录响应数据
   */
  const setAuth = (data: LoginResponse) => {
    token.value = data.token
    refreshToken.value = data.refresh_token
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('refresh_token', data.refresh_token)
    localStorage.setItem('user', JSON.stringify(data.user))
  }

  /**
   * 清除认证信息
   */
  const clearAuth = () => {
    token.value = ''
    refreshToken.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
  }

  /**
   * 退出登录
   */
  const logout = () => {
    clearAuth()
    router.push('/login')
  }

  return {
    token,
    refreshToken,
    user,
    isAuthenticated,
    logout,
    setAuth,
    clearAuth,
  }
})

