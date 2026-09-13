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

export interface LoginData {
  token?: string;
  user?: UserProfile;
  requires_email_verification: boolean;
  challenge_id?: string;
  masked_email?: string;
  expires_in_seconds?: number;
  development_verification_code?: string;
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

export interface MyContentItem {
  id: number;
  type: 'topic' | 'post';
  topic_id: number;
  title: string;
  excerpt: string;
  status: string;
  post_type?: Post['post_type'];
  cooling_ends_at?: string;
  created_at: string;
  updated_at: string;
}

export interface SearchResult {
  id: number;
  type: 'topic' | 'post';
  topic_id: number;
  title: string;
  excerpt: string;
  author_name: string;
  created_at: string;
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
  topic_id?: number;
  resource_type: 'topic' | 'guide';
  resource_key: string;
  progress: number;
  reading_seconds: number;
  reply_dwell_seconds: number;
  bottom_reached: boolean;
  requires_reply_dwell: boolean;
  eligible: boolean;
  completed: boolean;
  last_heartbeat_at: string;
}

export type FeedbackStance = 'support' | 'challenge';
export type FeedbackTag = 'logical' | 'new_perspective' | 'well_sourced' | 'empathetic' | 'factual_concern' | 'reasoning_gap' | 'inappropriate' | 'legacy_support';

export interface ContextualFeedback {
  id: number;
  user_id: number;
  target_type: 'topic' | 'post';
  target_id: number;
  stance: FeedbackStance;
  tag: FeedbackTag;
  reason: string;
  status: string;
  llm_audit_status: string;
  created_at: string;
  updated_at: string;
}

export interface FeedbackSummary {
  support: number;
  challenge: number;
  score: number;
  tags: Partial<Record<FeedbackTag, number>>;
  mine?: ContextualFeedback;
}

export interface CommentCluster {
  id: number;
  topic_id: number;
  tag: string;
  summary: string;
  weight: number;
  post_ids: number[];
  generation: number;
  created_at: string;
}

export interface ReviewTask {
  id: number;
  batch_id: number;
  subject_type: 'user' | 'topic';
  subject: {
    statement?: string;
    background_tag?: string;
    title?: string;
    claim?: string;
    evidence?: string;
    uncertainty?: string;
  };
  task_status: string;
  deadline: string;
  created_at: string;
}

export interface ReviewSubmission {
  id: number;
  batch_id: number;
  result: 'pass' | 'reject';
  llm_check_status: string;
  completed_at: string;
}

export interface AdminDailyTrend { date: string; users: number; topics: number; posts: number; feedback: number; }
export interface AdminOverview {
  users_total: number;
  topics_total: number;
  posts_total: number;
  active_users_7_days: number;
  new_users_today: number;
  feedback_total: number;
  bookmarks_total: number;
  suspended_users: number;
  verified_read_hours: number;
  content_status: Record<string, number>;
  trust_distribution: Record<string, number>;
  level_distribution: Record<string, number>;
  feedback_distribution: Record<string, number>;
  review_total: number;
  review_completed: number;
  review_expired: number;
  review_fair: number;
  llm_calls: number;
  llm_success: number;
  llm_average_ms: number;
  llm_prompt_tokens: number;
  llm_output_tokens: number;
  trend: AdminDailyTrend[];
  trend_days: number;
  categories: { name: string; topics: number; posts: number }[];
}
export interface AdminUser { id: number; username: string; email: string; role: string; status: string; unlock_level: number; trust_score: number; verified_read_seconds: number; onboarding_status: string; created_at: string; }
export interface AdminContent { id: number; type: 'topic' | 'post'; title: string; excerpt: string; author_name: string; status: string; created_at: string; }
export interface AdminLLMJob { id: number; job_type: string; aggregate_type: string; aggregate_id: number; status: string; attempts: number; model: string; error_message: string; prompt_tokens: number; completion_tokens: number; latency_ms: number; created_at: string; }
export interface AdminTrustLog { id: number; user_id: number; username: string; event_type: string; score_delta: number; reason: string; reference_type: string; reference_id?: number; created_at: string; }
