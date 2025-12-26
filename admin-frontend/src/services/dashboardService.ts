import api from './api'
import type { DashboardStats } from '@/types/api'

/**
 * 获取仪表板统计信息
 * @returns 仪表板统计数据
 */
export const getDashboardStats = async (): Promise<DashboardStats> => {
  return api.get('/dashboard/stats')
}
