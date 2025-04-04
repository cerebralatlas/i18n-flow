import api from '../utils/api';
import {
  TranslationRequest,
  BatchTranslationRequest,
} from '../types/translation';

export const translationService = {
  // 获取项目的翻译列表
  getTranslationsByProject: async (
    projectId: number,
    page: number = 1,
    pageSize: number = 10,
    keyword: string = ''
  ) => {
    const response = await api.get(`/api/translations/by-project/${projectId}`, {
      params: { page, page_size: pageSize, keyword }
    });
    return response.data;
  },

  // 获取翻译详情
  getTranslationById: async (id: number) => {
    const response = await api.get(`/api/translations/${id}`);
    return response.data;
  },

  // 创建翻译
  createTranslation: async (translation: TranslationRequest) => {
    const response = await api.post('/api/translations', translation);
    return response.data;
  },

  // 批量创建翻译
  batchCreateTranslations: async (request: BatchTranslationRequest) => {
    const response = await api.post('/api/translations/batch', request);
    return response.data;
  },

  // 更新翻译
  updateTranslation: async (id: number, translation: TranslationRequest) => {
    const response = await api.put(`/api/translations/${id}`, translation);
    return response.data;
  },

  // 删除翻译
  deleteTranslation: async (id: number) => {
    const response = await api.delete(`/api/translations/${id}`);
    return response.data;
  },

  // 导出项目翻译
  exportTranslations: async (projectId: number, format: string = 'json') => {
    const response = await api.get(`/api/exports/project/${projectId}`, {
      params: { format }
    });
    return response.data;
  },

  // 导入项目翻译
  importTranslations: async (projectId: number, data: Record<string, Record<string, string>>) => {
    const response = await api.post(`/api/imports/project/${projectId}`, data);
    return response.data;
  },

  // 获取所有语言
  getLanguages: async () => {
    const response = await api.get('/api/languages');
    return response.data;
  },

  // 批量删除翻译
  batchDeleteTranslations: async (ids: number[]) => {
    return await api.post('/api/translations/batch-delete', ids);
  }
};