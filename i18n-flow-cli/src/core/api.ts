import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import { getConfig, I18nFlowConfig } from './config';

// API响应类型
export interface ApiResponse<T> {
  data: T;
  status: number;
  statusText: string;
}

// 键推送请求
export interface KeysPushRequest {
  project_id: string;
  keys: string[];
  defaults: Record<string, string>;
}

// 键推送响应
export interface KeysPushResponse {
  added: string[];
  existed: string[];
  failed: string[];
}


// API客户端类
export class ApiClient {
  private client!: AxiosInstance;

  constructor() {
    const config = getConfig();
    this.updateClient(config);
  }

  /**
   * 更新API客户端
   */
  updateClient(config: I18nFlowConfig) {
    this.client = axios.create({
      baseURL: `${config.serverUrl}/api`,
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': config.apiKey
      },
      timeout: 30000
    });
  }

  /**
   * 测试API连接
   */
  async testConnection(config: I18nFlowConfig): Promise<boolean> {
    try {
      if (config) {
        this.updateClient(config);
      }
      const response = await this.client.get('/cli/auth');
      return response.status === 200 && response.data.status === 'ok';
    } catch (error) {
      return false;
    }
  }

  /**
   * 获取翻译数据
   */
  async getTranslations(projectId?: string, locale?: string): Promise<Record<string, any>> {
    const params: Record<string, string> = {};
    if (projectId) params.project_id = projectId;
    if (locale) params.locale = locale;

    const response = await this.client.get('/cli/translations', { params });
    return response.data;
  }

  /**
   * 推送新的翻译键
   */
  async pushKeys(request: KeysPushRequest): Promise<KeysPushResponse> {
    const response = await this.client.post('/cli/keys', request);
    return response.data;
  }
}

// 导出API客户端实例
export default new ApiClient();