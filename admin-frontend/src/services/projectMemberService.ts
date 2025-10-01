import api from '../utils/api';
import { ApiResponse } from '../types/api';
import { Project } from '../types/project';

// 项目成员相关类型定义
export interface ProjectMember {
  id: number;
  project_id: number;
  user_id: number;
  role: 'owner' | 'editor' | 'viewer';
  created_at: string;
  updated_at: string;
}

export interface ProjectMemberInfo {
  id: number;
  user_id: number;
  username: string;
  email: string;
  role: 'owner' | 'editor' | 'viewer';
}

export interface AddProjectMemberRequest {
  user_id: number;
  role: 'owner' | 'editor' | 'viewer';
}

export interface UpdateProjectMemberRequest {
  role: 'owner' | 'editor' | 'viewer';
}

export const projectMemberService = {
  // 添加项目成员
  addMember: async (projectId: number, memberData: AddProjectMemberRequest): Promise<ProjectMember> => {
    const response: ApiResponse<ProjectMember> = await api.post(`/api/projects/${projectId}/members`, memberData);
    return response.data;
  },

  // 获取项目成员列表
  getProjectMembers: async (projectId: number): Promise<ProjectMemberInfo[]> => {
    const response: ApiResponse<ProjectMemberInfo[]> = await api.get(`/api/projects/${projectId}/members`);
    return response.data;
  },

  // 获取用户参与的项目列表
  getUserProjects: async (userId: number): Promise<Project[]> => {
    const response: ApiResponse<Project[]> = await api.get(`/api/user-projects/${userId}`);
    return response.data;
  },

  // 更新成员角色
  updateMemberRole: async (projectId: number, userId: number, roleData: UpdateProjectMemberRequest): Promise<ProjectMember> => {
    const response: ApiResponse<ProjectMember> = await api.put(`/api/projects/${projectId}/members/${userId}`, roleData);
    return response.data;
  },

  // 移除项目成员
  removeMember: async (projectId: number, userId: number): Promise<void> => {
    await api.delete(`/api/projects/${projectId}/members/${userId}`);
  },

  // 检查用户权限
  checkPermission: async (projectId: number, userId: number, requiredRole: string): Promise<boolean> => {
    const response: ApiResponse<{ has_permission: boolean }> = await api.get(
      `/api/projects/${projectId}/members/${userId}/permission`,
      { params: { required_role: requiredRole } }
    );
    return response.data.has_permission;
  }
};
