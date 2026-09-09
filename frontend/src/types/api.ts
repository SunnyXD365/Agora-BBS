/**
 * 通用响应结构
 */
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

export interface Pagination {
  page: number;
  page_size: number;
  total: number;
}

export interface PaginatedData<T> {
  list: T[];
  pagination: Pagination;
}

/**
 * 业务数据实体
 */
// 没有完全实现数据库字段
export interface User {
  id: number;
  username: string;
  email?: string;
  trust_score: number;
  unlock_level?: number;
  created_at?: string;
  updated_at?: string;
}

export interface Topic {
  id: number;
  category_id: number;
  author_id: number;
  author_name?: string;// 这个可以通过id查询
  title: string;
  content?: string;
  structured_content?: string;// 格式有待考查
  status: 'published' | 'draft' | 'archived' | string;// 这里和数据库的意图不同
  view_count: number;
  reply_count: number;
  created_at: string;
  updated_at?: string;
}

export interface Post {
  id: number;
  topic_id: number;
  author_id: number;
  parent_id: number | null;
  content: string;
  post_type: 'reply' | 'topic' | string;// 这里和数据库的意图不同
  status: 'published' | 'hidden' | string;// 这里和数据库的意图不同
  cooling_ends_at?: string | null;
  created_at: string;
  updated_at?: string;
}

/**
 * 请求/响应 DTO
 */

// Auth 模块
export interface RegisterParams {
  username: string;
  password: string;
  email: string;
}

export interface LoginParams {
  username: string;
  password: string;
}

export interface LoginResponseData {
  token: string;
  user: User;
}

// Topics 模块
export interface GetTopicsParams {
  page?: number;
  page_size?: number;
  category_id?: number;
}

export interface CreateTopicParams {
  category_id: number;
  title: string;
  content: string;
}

// Posts 模块
export interface CreatePostParams {
  parent_id: number | null;
  content: string;
}
