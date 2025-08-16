import api from '../utils/api';

export interface DashboardStats {
  total_projects: number;
  total_translations: number;
  total_languages: number;
  total_keys: number;
}

/**
 * 获取仪表板统计数据
 */
export async function getDashboardStats() {
  return api.get("/api/dashboard/stats");
}
