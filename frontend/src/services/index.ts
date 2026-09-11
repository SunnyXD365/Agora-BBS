import api from '@/lib/api';
import {
  ApiResponse,
  AuthData,
  Bookmark,
  Category,
  PageData,
  Post,
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
};

export const topicApi = {
  getCategories: () =>
    api.get<ApiResponse<Category[]>, ApiResponse<Category[]>>('/categories'),
  getTopics: (params?: { category_id?: number; page?: number; page_size?: number }) =>
    api.get<ApiResponse<PageData<Topic>>, ApiResponse<PageData<Topic>>>('/topics', { params }),
  getTopicDetail: (id: number) =>
    api.get<ApiResponse<Topic>, ApiResponse<Topic>>(`/topics/${id}`),
  createTopic: (data: { category_id: number; title: string; content?: string; structured_content?: StructuredContent }) =>
    api.post<ApiResponse<{ topic_id: number }>, ApiResponse<{ topic_id: number }>>('/topics', data),
};

export const postApi = {
  getPosts: (topicId: number, params?: { page?: number; page_size?: number }) =>
    api.get<ApiResponse<PageData<Post>>, ApiResponse<PageData<Post>>>(`/topics/${topicId}/posts`, { params }),
  createPost: (topicId: number, data: { content: string; parent_id?: number; post_type?: Post['post_type'] }) =>
    api.post<ApiResponse<{ post_id: number }>, ApiResponse<{ post_id: number }>>(`/topics/${topicId}/posts`, data),
};

export const likeApi = {
  toggleLike: (target_type: 'topic' | 'post', target_id: number, isLike: boolean) => {
    const payload = { target_type, target_id };
    return isLike
      ? api.post<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>('/likes', payload)
      : api.delete<ApiResponse<{ status: string }>, ApiResponse<{ status: string }>>('/likes', { data: payload });
  },
};

export const bookmarkApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    api.get<ApiResponse<PageData<Bookmark>>, ApiResponse<PageData<Bookmark>>>('/bookmarks', { params }),
  create: (topicId: number) =>
    api.put<ApiResponse<{ bookmarked: boolean }>, ApiResponse<{ bookmarked: boolean }>>(`/bookmarks/${topicId}`),
  remove: (topicId: number) =>
    api.delete<ApiResponse<{ bookmarked: boolean }>, ApiResponse<{ bookmarked: boolean }>>(`/bookmarks/${topicId}`),
};
