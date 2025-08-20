import api from '../utils/api';
import { ApiResponse } from '../types/api';

export interface DashboardStats {
  total_projects: number;
  total_translations: number;
  total_languages: number;
  total_keys: number;
}

/**
 * 获取仪表板统计数据
 */
export async function getDashboardStats(): Promise<DashboardStats> {
  const response: ApiResponse<DashboardStats> = await api.get("/api/dashboard/stats");
  return response.data;
}
