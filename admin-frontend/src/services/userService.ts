import api from '../utils/api';
import { ApiResponse, PaginatedResponse } from '../types/api';

// 用户相关类型定义
export interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'member' | 'viewer';
  status: 'active' | 'disabled';
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
  role: 'admin' | 'member' | 'viewer';
}

export interface UpdateUserRequest {
  username?: string;
  email?: string;
  role?: 'admin' | 'member' | 'viewer';
  status?: 'active' | 'disabled';
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export const userService = {
  // 获取用户列表
  getUsers: async (page: number = 1, pageSize: number = 10, keyword: string = ''): Promise<PaginatedResponse<User>> => {
    const response: PaginatedResponse<User> = await api.get('/api/users', {
      params: { page, page_size: pageSize, keyword }
    });
    return response;
  },

  // 获取用户详情
  getUserById: async (id: number): Promise<User> => {
    const response: ApiResponse<User> = await api.get(`/api/users/${id}`);
    return response.data;
  },

  // 创建用户
  createUser: async (user: CreateUserRequest): Promise<User> => {
    const response: ApiResponse<User> = await api.post('/api/users', user);
    return response.data;
  },

  // 更新用户
  updateUser: async (id: number, user: UpdateUserRequest): Promise<User> => {
    const response: ApiResponse<User> = await api.put(`/api/users/${id}`, user);
    return response.data;
  },

  // 删除用户
  deleteUser: async (id: number): Promise<void> => {
    await api.delete(`/api/users/${id}`);
  },

  // 修改密码
  changePassword: async (passwordData: ChangePasswordRequest): Promise<void> => {
    await api.post('/api/user/change-password', passwordData);
  },

  // 获取当前用户信息
  getCurrentUser: async (): Promise<User> => {
    const response: ApiResponse<User> = await api.get('/api/user/info');
    return response.data;
  }
};
