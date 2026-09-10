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

export interface User {
  id: number;
  username: string;
  email?: string;
  avatar?: string;
  trust_score: number;
  unlock_level?: number;
  created_at?: string;
  updated_at?: string;
}

export interface Category {
  id: number;
  name: string;
  slug: string;
  description?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Topic {
  id: number;
  category_id: number;
  user_id: number;
  author_name?: string;
  title: string;
  content: string;
  view_count: number;
  post_count: number;
  like_count: number;
  is_sticky: boolean;
  is_essence: boolean;
  status: string;
  created_at: string;
  updated_at?: string;
}

export interface Post {
  id: number;
  topic_id: number;
  user_id: number;
  author_name?: string;
  author_avatar?: string;
  parent_id: number | null;
  content: string;
  like_count: number;
  status: string;
  created_at: string;
  updated_at?: string;
  children?: Post[];
}

// Request DTOs
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

export interface CreatePostParams {
  topic_id: number;
  parent_id?: number | null;
  content: string;
}
