import api from '../utils/api';
import { ProjectFormData } from '../types/project';

export const projectService = {
  // 获取项目列表
  getProjects: async (page: number = 1, pageSize: number = 10, keyword: string = '') => {
    const response = await api.get('/api/projects', {
      params: { page, page_size: pageSize, keyword }
    });
    return response?.success ? response.data : response;
  },

  // 获取项目详情
  getProjectById: async (id: number) => {
    const response = await api.get(`/api/projects/detail/${id}`);
    return response?.success ? response.data : response;
  },

  // 创建项目
  createProject: async (project: ProjectFormData) => {
    const response = await api.post('/api/projects', project);
    return response?.success ? response.data : response;
  },

  // 更新项目
  updateProject: async (id: number, project: ProjectFormData) => {
    const response = await api.put(`/api/projects/update/${id}`, project);
    return response?.success ? response.data : response;
  },

  // 删除项目
  deleteProject: async (id: number) => {
    const response = await api.delete(`/api/projects/delete/${id}`);
    return response?.success ? response.data : response;
  }
};