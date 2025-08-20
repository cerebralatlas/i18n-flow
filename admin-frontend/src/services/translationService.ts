import api from '../utils/api';
import {
  TranslationRequest,
  BatchTranslationRequest,
} from '../types/translation';
import { TranslationMatrix } from '../components/translation/TranslationTable';
import { ApiResponse, PaginatedResponse } from '../types/api';

// 数据转换函数：将后端key-language映射转换为前端TranslationMatrix数组
const transformMatrixData = (matrixData: Record<string, Record<string, string>> | null | undefined): TranslationMatrix[] => {
  // 防御性检查
  if (!matrixData || typeof matrixData !== 'object') {
    console.warn('Matrix data is null, undefined, or not an object:', matrixData);
    return [];
  }

  return Object.entries(matrixData).map(([keyName, languages]) => ({
    key_name: keyName,
    context: '', // 后端暂未返回context信息
    languages: languages && typeof languages === 'object' 
      ? Object.entries(languages).reduce((acc, [langCode, value]) => {
          acc[langCode] = {
            id: 0, // 临时设置为0，因为后端返回的matrix数据中没有具体的翻译记录ID
            value: value || ''
          };
          return acc;
        }, {} as Record<string, { id: number; value: string }>)
      : {}
  }));
};

export const translationService = {
  // 获取项目的翻译列表
  getTranslationsByProject: async (
    projectId: number,
    page: number = 1,
    pageSize: number = 10,
    keyword: string = ''
  ) => {
    const response: ApiResponse = await api.get(`/api/translations/by-project/${projectId}`, {
      params: { page, page_size: pageSize, keyword }
    });
    // 处理统一响应格式
    return response?.success ? response.data : response;
  },

  // 获取翻译详情
  getTranslationById: async (id: number) => {
    const response: ApiResponse = await api.get(`/api/translations/${id}`);
    return response.data;
  },

  // 创建翻译
  createTranslation: async (translation: TranslationRequest) => {
    const response: ApiResponse = await api.post('/api/translations', translation);
    return response.data;
  },

  // 批量创建翻译
  batchCreateTranslations: async (request: BatchTranslationRequest) => {
    const response: ApiResponse = await api.post('/api/translations/batch', request);
    return response.data;
  },

  // 更新翻译
  updateTranslation: async (id: number, translation: TranslationRequest) => {
    const response: ApiResponse = await api.put(`/api/translations/${id}`, translation);
    return response.data;
  },

  // 删除翻译
  deleteTranslation: async (id: number) => {
    await api.delete(`/api/translations/${id}`);
  },

  // 导出项目翻译
  exportTranslations: async (projectId: number, format: string = 'json') => {
    const response: ApiResponse = await api.get(`/api/exports/project/${projectId}`, {
      params: { format }
    });
    return response.data;
  },

  // 导入项目翻译
  importTranslations: async (projectId: number, data: Record<string, Record<string, string>>) => {
    const response: ApiResponse = await api.post(`/api/imports/project/${projectId}`, data);
    return response.data;
  },

  // 获取所有语言
  getLanguages: async () => {
    const response: ApiResponse = await api.get('/api/languages');
    return response.data;
  },

  // 批量删除翻译
  batchDeleteTranslations: async (ids: number[]) => {
    await api.post('/api/translations/batch-delete', ids);
  },

  // 获取翻译矩阵
  getTranslationMatrix: async (projectId: number, page: number, pageSize: number, keyword: string = "") => {
    const backendResponse: ApiResponse = await api.get(`/api/translations/matrix/by-project/${projectId}`, {
      params: { page, page_size: pageSize, keyword }
    });
    
    // 防御性检查 - API拦截器已经返回了response.data，所以backendResponse就是{ success, data, meta }
    if (!backendResponse) {
      console.warn('Backend response is null or undefined');
      return {
        data: [],
        total: 0,
        current: 1,
        pageSize: pageSize,
        totalPages: 0
      };
    }

    // 添加调试日志
    console.log('Backend response:', backendResponse);
    
    // 转换数据格式
    const transformedData = transformMatrixData(backendResponse.data);
    
    // 返回前端期望的格式
    return {
      data: transformedData,
      total: backendResponse.meta?.total_count || 0,
      current: backendResponse.meta?.page || 1,
      pageSize: backendResponse.meta?.page_size || pageSize,
      totalPages: backendResponse.meta?.total_pages || 0
    };
  }
};