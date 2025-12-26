/**
 * 后端统一 API 响应结构
 */
export interface APIResponse<T = any> {
  success: boolean
  data?: T
  error?: ErrorInfo
  meta?: PaginationMeta
}

/**
 * 错误信息结构
 */
export interface ErrorInfo {
  code: string
  message: string
  details?: string
}

/**
 * 分页元数据
 */
export interface PaginationMeta {
  page: number
  page_size: number
  total_count: number
  total_pages: number
}

/**
 * 用户信息
 */
export interface User {
  id: number
  username: string
  email?: string
  role?: string
  status?: string
  created_at: string
  updated_at: string
}

/**
 * 登录响应数据
 */
export interface LoginResponse {
  token: string
  refresh_token: string
  user: User
}

/**
 * 登录请求参数
 */
export interface LoginParams {
  username: string
  password: string
}

/**
 * 刷新 Token 请求参数
 */
export interface RefreshTokenParams {
  refresh_token: string
}

/**
 * 仪表板统计信息
 */
export interface DashboardStats {
  total_projects: number
  total_languages: number
  total_keys: number
  total_translations: number
}

/**
 * 项目实体
 */
export interface Project {
  id: number
  name: string
  slug: string
  description?: string
  status: 'active' | 'archived'
  created_at: string
  updated_at: string
}

/**
 * 创建项目请求参数
 */
export interface CreateProjectRequest {
  name: string
  description?: string
}

/**
 * 更新项目请求参数
 */
export interface UpdateProjectRequest {
  name?: string
  description?: string
  status?: 'active' | 'archived'
}

/**
 * 项目列表查询参数
 */
export interface ProjectListParams {
  page?: number
  page_size?: number
  keyword?: string
}

/**
 * 项目列表响应数据
 */
export interface ProjectListResponse {
  data: Project[]
  meta: PaginationMeta
}
