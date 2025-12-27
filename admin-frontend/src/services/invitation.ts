import api from './api'
import type {
  ValidateInvitationResponse,
  RegisterParams,
  RegisterResponse,
  CreateInvitationParams,
  CreateInvitationResponse,
  InvitationListResponse,
  Invitation,
} from '@/types/api'

/**
 * 验证邀请码是否有效
 * @param code 邀请码
 * @returns 验证结果
 */
export const validateInvitation = async (code: string): Promise<ValidateInvitationResponse> => {
  return api.get(`/invitations/${code}/validate`)
}

/**
 * 使用邀请码注册新用户
 * @param params 注册参数
 * @returns 注册响应
 */
export const registerWithInvitation = async (params: RegisterParams): Promise<RegisterResponse> => {
  return api.post('/register', params)
}

/**
 * 创建邀请链接（管理员功能）
 * @param params 创建邀请参数
 * @returns 创建结果
 */
export const createInvitation = async (params: CreateInvitationParams): Promise<CreateInvitationResponse> => {
  return api.post('/invitations', params)
}

/**
 * 获取邀请列表（管理员功能）
 * @param page 页码
 * @param pageSize 每页数量
 * @returns 邀请列表
 */
export const getInvitations = async (
  page: number = 1,
  pageSize: number = 20
): Promise<InvitationListResponse> => {
  return api.get('/invitations', { params: { page, page_size: pageSize } })
}

/**
 * 获取邀请详情（管理员功能）
 * @param code 邀请码
 * @returns 邀请详情
 */
export const getInvitation = async (code: string): Promise<Invitation> => {
  return api.get(`/invitations/${code}`)
}

/**
 * 撤销邀请（管理员功能）
 * @param code 邀请码
 */
export const revokeInvitation = async (code: string): Promise<void> => {
  return api.delete(`/invitations/${code}`)
}
