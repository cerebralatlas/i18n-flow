// 翻译请求参数
export interface TranslationRequest {
  project_id: number;
  key_name: string;
  context?: string;
  language_id: number;
  value: string;
}

// 批量翻译请求
export interface BatchTranslationRequest {
  project_id: number;
  key_name: string;
  context?: string;
  translations: Record<string, string>; // 语言代码 -> 翻译值
}

// 翻译响应
export interface TranslationResponse {
  id: number;
  project_id: number;
  key_name: string;
  context: string;
  language_id: number;
  value: string;
  status: string;
  created_at: string;
  updated_at: string;
  project_name: string;
  language_code: string;
  language_name: string;
}

// 语言类型
export interface Language {
  id: number;
  code: string;
  name: string;
  is_default: boolean;
  status: string;
  created_at: string;
  updated_at: string;
}