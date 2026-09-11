import api from '@/lib/api';
import {
  ApiResponse,
  AdminContent,
  AdminLLMJob,
  AdminOverview,
  AdminTrustLog,
  AdminUser,
  AuthData,
  Bookmark,
  Category,
  CommentCluster,
  ContextualFeedback,
  FeedbackStance,
  FeedbackSummary,
  FeedbackTag,
  PageData,
  Post,
  GovernancePolicy,
  ReadingSession,
  ReviewSubmission,
  ReviewTask,
  StructuredContent,
  Topic,
  UserProfile,
} from '@/types/api';

export function getErrorMessage(error: unknown, fallback: string): string {
  if (typeof error === 'object' && error !== null && 'msg' in error) {
    const msg = (error as { msg?: unknown }).msg;
    if (typeof msg === 'string' && msg) return msg;
  }
  return fallback;
}

export const authApi = {
  register: (data: Record<string, string>) =>
    api.post<ApiResponse<AuthData>, ApiResponse<AuthData>>('/auth/register', data),
  login: (data: Record<string, string>) =>
    api.post<ApiResponse<AuthData>, ApiResponse<AuthData>>('/auth/login', data),
  getMe: () =>
    api.get<ApiResponse<UserProfile>, ApiResponse<UserProfile>>('/users/me'),
  saveOnboarding: (data: { statement: string; background_tag: string }) =>
    api.put<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>('/users/me/onboarding', data),
};

export const topicApi = {
  getCategories: () =>
    api.get<ApiResponse<Category[]>, ApiResponse<Category[]>>('/categories'),
  getTopics: (params?: { category_id?: number; page?: number; page_size?: number }) =>
    api.get<ApiResponse<PageData<Topic>>, ApiResponse<PageData<Topic>>>('/topics', { params }),
  getTopicDetail: (id: number) =>
    api.get<ApiResponse<Topic>, ApiResponse<Topic>>(`/topics/${id}`),
  createTopic: (data: { category_id: number; title: string; content?: string; structured_content?: StructuredContent }) =>
    api.post<ApiResponse<{ topic_id: number; status: string; cooling_ends_at: string }>, ApiResponse<{ topic_id: number; status: string; cooling_ends_at: string }>>('/topics', data),
  updateCooling: (id: number, data: { title: string; structured_content: StructuredContent }) =>
    api.patch<ApiResponse<Topic>, ApiResponse<Topic>>(`/topics/${id}`, data),
  recallCooling: (id: number) =>
    api.delete<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>(`/topics/${id}`),
};

export const postApi = {
  getPosts: (topicId: number, params?: { page?: number; page_size?: number }) =>
    api.get<ApiResponse<PageData<Post>>, ApiResponse<PageData<Post>>>(`/topics/${topicId}/posts`, { params }),
  createPost: (topicId: number, data: { content: string; parent_id?: number; post_type?: Post['post_type'] }) =>
    api.post<ApiResponse<{ post_id: number; status: string; cooling_ends_at: string }>, ApiResponse<{ post_id: number; status: string; cooling_ends_at: string }>>(`/topics/${topicId}/posts`, data),
  updateCooling: (id: number, data: { content: string; post_type: Post['post_type'] }) =>
    api.patch<ApiResponse<Post>, ApiResponse<Post>>(`/posts/${id}`, data),
  recallCooling: (id: number) =>
    api.delete<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>(`/posts/${id}`),
};

export const likeApi = {
  toggleLike: (target_type: 'topic' | 'post', target_id: number, isLike: boolean) => {
    const payload = { target_type, target_id };
    return isLike
      ? api.post<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>('/likes', payload)
      : api.delete<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>('/likes', { data: payload });
  },
};

export const governanceApi = {
  policy: () => api.get<ApiResponse<GovernancePolicy>, ApiResponse<GovernancePolicy>>('/governance/policy'),
  startReading: (topicId: number) => api.post<ApiResponse<ReadingSession>, ApiResponse<ReadingSession>>('/reading-sessions', { topic_id: topicId }),
  heartbeat: (id: string, progress: number, replyFocused: boolean) =>
    api.patch<ApiResponse<ReadingSession>, ApiResponse<ReadingSession>>(`/reading-sessions/${id}/heartbeat`, { progress, reply_focused: replyFocused }),
  completeReading: (id: string) => api.post<ApiResponse<ReadingSession>, ApiResponse<ReadingSession>>(`/reading-sessions/${id}/complete`),
};

export const bookmarkApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    api.get<ApiResponse<PageData<Bookmark>>, ApiResponse<PageData<Bookmark>>>('/bookmarks', { params }),
  create: (topicId: number) =>
    api.put<ApiResponse<{ bookmarked: boolean }>, ApiResponse<{ bookmarked: boolean }>>(`/bookmarks/${topicId}`),
  remove: (topicId: number) =>
    api.delete<ApiResponse<{ bookmarked: boolean }>, ApiResponse<{ bookmarked: boolean }>>(`/bookmarks/${topicId}`),
};

export const feedbackApi = {
  summary: (targetType: 'topic' | 'post', targetId: number) =>
    api.get<ApiResponse<FeedbackSummary>, ApiResponse<FeedbackSummary>>('/feedbacks/summary', { params: { target_type: targetType, target_id: targetId } }),
  upsert: (data: { target_type: 'topic' | 'post'; target_id: number; stance: FeedbackStance; tag: FeedbackTag; reason: string }) =>
    api.post<ApiResponse<ContextualFeedback>, ApiResponse<ContextualFeedback>>('/feedbacks', data),
  withdraw: (id: number) =>
    api.delete<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>(`/feedbacks/${id}`),
  clusters: (topicId: number) =>
    api.get<ApiResponse<CommentCluster[]>, ApiResponse<CommentCluster[]>>(`/topics/${topicId}/clusters`),
};

export const reviewApi = {
  list: () => api.get<ApiResponse<ReviewTask[]>, ApiResponse<ReviewTask[]>>('/reviews/tasks'),
  submit: (id: number, data: { appropriateness: boolean; sincerity: boolean; reason: string }) =>
    api.post<ApiResponse<ReviewSubmission>, ApiResponse<ReviewSubmission>>(`/reviews/tasks/${id}`, data),
};

export const adminApi = {
  overview: () => api.get<ApiResponse<AdminOverview>, ApiResponse<AdminOverview>>('/admin/overview'),
  users: (params?: { page?: number; page_size?: number }) => api.get<ApiResponse<PageData<AdminUser>>, ApiResponse<PageData<AdminUser>>>('/admin/users', { params }),
  setUserStatus: (id: number, status: 'active' | 'suspended') => api.patch<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>(`/admin/users/${id}/status`, { status }),
  contents: (type: 'topic' | 'post', params?: { page?: number; page_size?: number; status?: string }) => api.get<ApiResponse<PageData<AdminContent>>, ApiResponse<PageData<AdminContent>>>('/admin/contents', { params: { type, ...params } }),
  setContentVisibility: (type: 'topic' | 'post', id: number, hidden: boolean) => api.patch<ApiResponse<{ hidden: boolean }>, ApiResponse<{ hidden: boolean }>>(`/admin/contents/${type}/${id}/visibility`, { hidden }),
  llmJobs: (params?: { page?: number; page_size?: number; status?: string }) => api.get<ApiResponse<PageData<AdminLLMJob>>, ApiResponse<PageData<AdminLLMJob>>>('/admin/llm-jobs', { params }),
  retryLLMJob: (id: number) => api.post<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>(`/admin/llm-jobs/${id}/retry`),
  trustLogs: (params?: { page?: number; page_size?: number; user_id?: number }) => api.get<ApiResponse<PageData<AdminTrustLog>>, ApiResponse<PageData<AdminTrustLog>>>('/admin/trust-logs', { params }),
  categories: () => api.get<ApiResponse<Category[]>, ApiResponse<Category[]>>('/admin/categories'),
  createCategory: (data: Omit<Category, 'id'>) => api.post<ApiResponse<Category>, ApiResponse<Category>>('/admin/categories', data),
  updateCategory: (id: number, data: Omit<Category, 'id'>) => api.patch<ApiResponse<Category>, ApiResponse<Category>>(`/admin/categories/${id}`, data),
};
