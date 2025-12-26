import api from './api'
import type { LoginParams, LoginResponse, RefreshTokenParams } from '@/types/api'

/**
 * 用户登录
 * @param params 登录参数（用户名和密码）
 * @returns 登录响应（包含 token、refresh_token 和用户信息）
 */
export const login = async (params: LoginParams): Promise<LoginResponse> => {
  return api.post('/login', params)
}

/**
 * 刷新访问令牌
 * @param params 刷新令牌参数
 * @returns 新的登录响应（包含新的 token、refresh_token 和用户信息）
 */
export const refreshToken = async (params: RefreshTokenParams): Promise<LoginResponse> => {
  return api.post('/refresh', params)
}

