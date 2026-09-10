import {
  ApiResponse,
  User,
  Category,
  Topic,
  Post,
  RegisterParams,
  LoginParams,
  LoginResponseData,
  GetTopicsParams,
  CreateTopicParams,
  CreatePostParams,
} from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api';
const USE_MOCK = process.env.NEXT_PUBLIC_USE_MOCK === 'true';

const TOKEN_KEY = 'agora_jwt_token';

export const getToken = (): string | null => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(TOKEN_KEY);
};

export const setToken = (token: string): void => {
  if (typeof window !== 'undefined') {
    localStorage.setItem(TOKEN_KEY, token);
  }
};

export const removeToken = (): void => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem(TOKEN_KEY);
  }
};

interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | undefined>;
}

async function fetchClient<T>(endpoint: string, options: RequestOptions = {}): Promise<ApiResponse<T>> {
  const { params, headers, ...customConfig } = options;

  let url = `${API_BASE_URL}${endpoint}`;
  if (params) {
    const searchParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        searchParams.append(key, String(value));
      }
    });
    const queryString = searchParams.toString();
    if (queryString) {
      url += `?${queryString}`;
    }
  }

  const requestHeaders: HeadersInit = {
    'Content-Type': 'application/json',
    ...headers,
  };

  const token = getToken();
  if (token) {
    (requestHeaders as Record<string, string>)['Authorization'] = `Bearer ${token}`;
  }

  const config: RequestInit = {
    method: 'GET',
    headers: requestHeaders,
    ...customConfig,
  };

  try {
    const response = await fetch(url, config);

    if (response.status === 401) {
      removeToken();
      if (typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
        window.location.href = '/login';
      }
      throw new Error('登录已过期，请重新登录');
    }

    const data: ApiResponse<T> = await response.json();

    // 后端成功 code 约定为 0
    if (!response.ok || data.code !== 0) {
      throw new Error(data.message || '请求失败，请稍后重试');
    }

    return data;
  } catch (error: any) {
    console.error(`API Error [${endpoint}]:`, error);
    throw error;
  }
}

// Auth API
export const authApi = {
  register: (data: RegisterParams) => fetchClient<User>('/auth/register', { method: 'POST', body: JSON.stringify(data) }),
  login: (data: LoginParams) => fetchClient<LoginResponseData>('/auth/login', { method: 'POST', body: JSON.stringify(data) }),
  getMe: () => fetchClient<User>('/auth/me'),
};

// Categories API
export const categoriesApi = {
  getCategories: () => fetchClient<Category[]>('/categories'),
};

// Topics API
export const topicsApi = {
  getTopics: (params?: GetTopicsParams) => fetchClient<Topic[]>('/topics', { params: params as any }),
  getTopicDetail: (id: string | number) => fetchClient<Topic>(`/topics/${id}`),
  createTopic: (data: CreateTopicParams) => fetchClient<Topic>('/topics', { method: 'POST', body: JSON.stringify(data) }),
};

// Posts API
export const postsApi = {
  // 匹配后端的 GET /api/posts?topic_id=x
  getPosts: (topicId: string | number) => fetchClient<Post[]>('/posts', { params: { topic_id: topicId } }),
  // 匹配后端的 POST /api/posts 传参 { topic_id, parent_id, content }
  createPost: (data: CreatePostParams) => fetchClient<Post>('/posts', { method: 'POST', body: JSON.stringify(data) }),
};

// Likes API
export const likesApi = {
  toggleLike: (targetType: 'topic' | 'post', targetId: number) =>
    fetchClient<{ is_liked: boolean; like_count: number }>('/likes/toggle', {
      method: 'POST',
      body: JSON.stringify({ target_type: targetType, target_id: targetId }),
    }),
};
