// 统一 API 响应包装
export interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
}

// 用户类型
export interface UserProfile {
  id: number;
  username: string;
  email: string;
  avatar: string;
  role: string;
  trust_score: number;
  unlock_level: number;
  created_at: string;
}

export interface AuthData {
  token: string;
  user: UserProfile;
}

// 板块类型
export interface Category {
  id: number;
  name: string;
  slug: string;
  description: string;
  sort_order: number;
}

// 主题帖类型
export interface Topic {
  id: number;
  category_id: number;
  user_id: number;
  author_name: string;
  title: string;
  content: string;
  view_count: number;
  post_count: number;
  like_count: number;
  created_at: string;
  updated_at: string;
}

// 回复楼层类型
export interface Post {
  id: number;
  topic_id: number;
  user_id: number;
  author_name: string;
  parent_id: number | null;
  content: string;
  like_count: number;
  created_at: string;
}