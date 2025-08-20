import api from '../utils/api';
import { ProjectFormData, Project } from '../types/project';
import { ApiResponse, PaginatedResponse } from '../types/api';

export const projectService = {
  // 获取项目列表
  getProjects: async (page: number = 1, pageSize: number = 10, keyword: string = ''): Promise<PaginatedResponse<Project>> => {
    const response: PaginatedResponse<Project> = await api.get('/api/projects', {
      params: { page, page_size: pageSize, keyword }
    });
    return response;
  },

  // 获取项目详情
  getProjectById: async (id: number): Promise<Project> => {
    const response: ApiResponse<Project> = await api.get(`/api/projects/detail/${id}`);
    return response.data;
  },

  // 创建项目
  createProject: async (project: ProjectFormData): Promise<Project> => {
    const response: ApiResponse<Project> = await api.post('/api/projects', project);
    return response.data;
  },

  // 更新项目
  updateProject: async (id: number, project: ProjectFormData): Promise<Project> => {
    const response: ApiResponse<Project> = await api.put(`/api/projects/update/${id}`, project);
    return response.data;
  },

  // 删除项目
  deleteProject: async (id: number): Promise<void> => {
    await api.delete(`/api/projects/delete/${id}`);
  }
};