// API 响应基础类型
export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  meta?: {
    page: number;
    page_size: number;
    total_count: number;
    total_pages: number;
  };
}

// 分页响应类型
export interface PaginatedResponse<T = any> {
  success: boolean;
  data: T[];
  meta: {
    page: number;
    page_size: number;
    total_count: number;
    total_pages: number;
  };
}
