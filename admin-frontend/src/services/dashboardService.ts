import api from '../utils/api';

export interface DashboardStats {
  project_count: number;
  translation_count: number;
  language_count: number;
  user_count: number;
}

/**
 * 获取仪表板统计数据
 */
export async function getDashboardStats() {
  return api.get("/api/dashboard/stats");
}
