// 统一 API 响应包装
export interface ApiResponse<T> {
  code: number;
  msg: string;
  data: T;
  request_id: string;
}

export interface PageData<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}

// 用户类型
export interface UserProfile {
  id: number;
  username: string;
  email: string;
  avatar: string;
  role: string;
  status: string;
  unlock_level: number;
  verified_read_seconds: number;
  capabilities: string[];
  onboarding_statement: string;
  background_tag: string;
  onboarding_status: 'not_submitted' | 'pending_review' | 'approved' | 'rejected';
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
  is_active: boolean;
  requires_review: boolean;
}

export interface StructuredContent {
  claim: string;
  evidence: string;
  uncertainty: string;
}

// 主题帖类型
export interface Topic {
  id: number;
  category_id: number;
  user_id: number;
  author_name: string;
  title: string;
  content: string;
  structured_content: StructuredContent;
  status: string;
  cooling_ends_at?: string;
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
  post_type: 'debate' | 'evidence' | 'experience' | 'thanks';
  status: string;
  cooling_ends_at?: string;
  like_count: number;
  created_at: string;
}

export interface Bookmark {
  id: number;
  user_id: number;
  topic_id: number;
  created_at: string;
  topic: Topic;
}

export interface GovernancePolicy {
  cooling_seconds: number;
  reply_dwell_seconds: number;
  heartbeat_seconds: number;
  long_topic_chars: number;
  level_1_read_seconds: number;
  level_2_read_seconds: number;
  level_3_read_seconds: number;
}

export interface ReadingSession {
  id: string;
  topic_id: number;
  progress: number;
  reading_seconds: number;
  reply_dwell_seconds: number;
  bottom_reached: boolean;
  eligible: boolean;
  completed: boolean;
  last_heartbeat_at: string;
}
