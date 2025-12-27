import api from './api'
import type {
  User,
  UserListParams,
  UserListResponse,
  CreateUserRequest,
  UpdateUserRequest,
  ResetPasswordRequest,
} from '@/types/api'

/**
 * 获取用户列表
 * @param params 查询参数（分页、搜索）
 * @returns 用户列表响应
 */
export const getUsers = async (params: UserListParams = {}): Promise<UserListResponse> => {
  return api.get('/users', { params })
}

/**
 * 获取用户详情
 * @param id 用户 ID
 * @returns 用户详情
 */
export const getUser = async (id: number): Promise<User> => {
  return api.get(`/users/${id}`)
}

/**
 * 创建用户
 * @param data 用户信息
 * @returns 创建的用户
 */
export const createUser = async (data: CreateUserRequest): Promise<User> => {
  return api.post('/users', data)
}

/**
 * 更新用户
 * @param id 用户 ID
 * @param data 更新信息
 * @returns 更新后的用户
 */
export const updateUser = async (id: number, data: UpdateUserRequest): Promise<User> => {
  return api.put(`/users/${id}`, data)
}

/**
 * 删除用户
 * @param id 用户 ID
 */
export const deleteUser = async (id: number): Promise<void> => {
  return api.delete(`/users/${id}`)
}

/**
 * 重置用户密码
 * @param id 用户 ID
 * @param data 新密码
 */
export const resetUserPassword = async (id: number, data: ResetPasswordRequest): Promise<void> => {
  return api.post(`/users/${id}/reset-password`, data)
}
