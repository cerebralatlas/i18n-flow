import api from './api'
import type {
  Project,
  CreateProjectRequest,
  UpdateProjectRequest,
  ProjectListParams,
  ProjectListResponse,
} from '@/types/api'

/**
 * 获取项目列表
 * @param params 查询参数（分页、关键词搜索）
 * @returns 项目列表和分页信息
 */
export const getProjects = async (params?: ProjectListParams): Promise<ProjectListResponse> => {
  return api.get('/projects', { params })
}

/**
 * 根据 ID 获取项目详情
 * @param id 项目 ID
 * @returns 项目详情
 */
export const getProjectById = async (id: number): Promise<Project> => {
  return api.get(`/projects/detail/${id}`)
}

/**
 * 创建新项目
 * @param data 项目信息
 * @returns 创建的项目
 */
export const createProject = async (data: CreateProjectRequest): Promise<Project> => {
  return api.post('/projects', data)
}

/**
 * 更新项目
 * @param id 项目 ID
 * @param data 更新的项目信息
 * @returns 更新后的项目
 */
export const updateProject = async (
  id: number,
  data: UpdateProjectRequest
): Promise<Project> => {
  return api.put(`/projects/update/${id}`, data)
}

/**
 * 删除项目
 * @param id 项目 ID
 */
export const deleteProject = async (id: number): Promise<void> => {
  return api.delete(`/projects/delete/${id}`)
}
