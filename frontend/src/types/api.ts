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
export interface User {
  id: number;
  username: string;
  email?: string;
  trust_score: number;
  created_at?: string;
}

export interface Topic {
  id: number;
  category_id: number;
  author_id: number;
  author_name?: string;
  title: string;
  content?: string;
  structured_content?: string;
  status: 'published' | 'draft' | 'archived' | string;
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
  post_type: 'reply' | 'topic' | string;
  status: 'published' | 'hidden' | string;
  created_at: string;
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

// 3. Posts 模块
export interface CreatePostParams {
  parent_id: number | null;
  content: string;
}
