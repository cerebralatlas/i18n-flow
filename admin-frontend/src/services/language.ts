import api from './api'
import type { Language, CreateLanguageRequest } from '@/types/translation'

/**
 * 获取所有语言列表
 */
export const getLanguages = async (): Promise<Language[]> => {
  const response = await api.get('/languages')
  return response as unknown as Language[]
}

/**
 * 创建新语言
 */
export const createLanguage = async (data: CreateLanguageRequest): Promise<Language> => {
  const response = await api.post('/languages', data)
  return response as unknown as Language
}

/**
 * 更新语言
 */
export const updateLanguage = async (id: number, data: CreateLanguageRequest): Promise<Language> => {
  const response = await api.put(`/languages/${id}`, data)
  return response as unknown as Language
}

/**
 * 删除语言
 */
export const deleteLanguage = async (id: number): Promise<void> => {
  await api.delete(`/languages/${id}`)
}
