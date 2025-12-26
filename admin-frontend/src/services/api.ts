import axios, { AxiosError } from 'axios'
import type { APIResponse } from '@/types/api'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    const apiResponse = response.data as APIResponse

    // 验证后端是否返回了标准的 APIResponse 结构
    if (typeof apiResponse === 'object' && apiResponse !== null) {
      // 如果有 success 字段，检查是否成功
      if ('success' in apiResponse) {
        if (apiResponse.success) {
          // 成功：如果有 meta 字段（分页数据），返回包含 data 和 meta 的对象
          if (apiResponse.meta) {
            return {
              data: apiResponse.data,
              meta: apiResponse.meta,
            }
          }
          // 否则只返回 data 字段（解包一层）
          return apiResponse.data
        } else {
          // success 为 false，抛出错误
          const errorMessage = apiResponse.error?.message || '请求失败'
          const error = new Error(errorMessage)
          error.name = apiResponse.error?.code || 'API_ERROR'
          return Promise.reject(error)
        }
      }
      // 如果没有 success 字段，直接返回原数据（兼容某些特殊接口）
      return apiResponse
    }

    // 如果响应不是对象，直接返回
    return response.data
  },
  (error: AxiosError) => {
    const { response } = error

    if (response) {
      // 尝试从响应中获取后端返回的错误信息
      const apiResponse = response.data as APIResponse

      switch (response.status) {
        case 401: {
          // 未授权：清除认证信息并跳转到登录页
          localStorage.removeItem('token')
          localStorage.removeItem('refresh_token')
          localStorage.removeItem('user')
          window.location.href = '/login'

          const err = new Error(apiResponse?.error?.message || '未授权，请重新登录')
          err.name = 'UNAUTHORIZED'
          return Promise.reject(err)
        }
        case 400: {
          const err = new Error(apiResponse?.error?.message || '请求参数错误')
          err.name = apiResponse?.error?.code || 'BAD_REQUEST'
          return Promise.reject(err)
        }
        case 403: {
          const err = new Error(apiResponse?.error?.message || '无权访问')
          err.name = 'FORBIDDEN'
          return Promise.reject(err)
        }
        case 404: {
          const err = new Error(apiResponse?.error?.message || '资源不存在')
          err.name = 'NOT_FOUND'
          return Promise.reject(err)
        }
        case 409: {
          const err = new Error(apiResponse?.error?.message || '资源冲突')
          err.name = 'CONFLICT'
          return Promise.reject(err)
        }
        case 500: {
          const err = new Error(apiResponse?.error?.message || '服务器内部错误')
          err.name = 'INTERNAL_SERVER_ERROR'
          return Promise.reject(err)
        }
        default: {
          const err = new Error(apiResponse?.error?.message || '请求失败')
          err.name = apiResponse?.error?.code || 'REQUEST_ERROR'
          return Promise.reject(err)
        }
      }
    }

    // 网络错误或请求超时
    const err = new Error('网络连接失败，请检查网络设置')
    err.name = 'NETWORK_ERROR'
    return Promise.reject(err)
  }
)

export default api
