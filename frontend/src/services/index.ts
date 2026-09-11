import api from '@/lib/api';
import { ApiResponse, AuthData, Category, Topic, Post, UserProfile } from '@/types/api';

export const authApi = {
  register: (data: Record<string, string>) =>
    api.post<any, ApiResponse<AuthData>>('/auth/register', data),
  login: (data: Record<string, string>) =>
    api.post<any, ApiResponse<AuthData>>('/auth/login', data),
  getMe: () =>
    api.get<any, ApiResponse<UserProfile>>('/users/me'),
};

export const topicApi = {
  getCategories: () =>
    api.get<any, ApiResponse<Category[]>>('/categories'),
  getTopics: (params?: { category_id?: number; page?: number; page_size?: number }) =>
    api.get<any, ApiResponse<Topic[]>>('/topics', { params }),
  getTopicDetail: (id: number) =>
    api.get<any, ApiResponse<Topic>>(`/topics/${id}`),
  createTopic: (data: { category_id: number; title: string; content: string }) =>
    api.post<any, ApiResponse<{ topic_id: number }>>('/topics', data),
};

export const postApi = {
  getPosts: (topicId: number, params?: { page?: number; page_size?: number }) =>
    api.get<any, ApiResponse<Post[]>>(`/topics/${topicId}/posts`, { params }),
  createPost: (topicId: number, data: { content: string; parent_id?: number }) =>
    api.post<any, ApiResponse<{ post_id: number }>>(`/topics/${topicId}/posts`, data),
};

export const likeApi = {
  toggleLike: (target_type: 'topic' | 'post', target_id: number, isLike: boolean) => {
    const payload = { target_type, target_id };
    return isLike
      ? api.post<any, ApiResponse<{ status: string }>>('/likes', payload)
      : api.delete<any, ApiResponse<{ status: string }>>('/likes', { data: payload });
  },
};