import {
  ApiResponse,
  User,
  Topic,
  Post,
  PaginatedData,
  RegisterParams,
  LoginParams,
  LoginResponseData,
  GetTopicsParams,
  CreateTopicParams,
  CreatePostParams,
} from '@/types/api';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api/v1';
const USE_MOCK = process.env.NEXT_PUBLIC_USE_MOCK === 'true';

// Token 本地存储 Key
const TOKEN_KEY = 'agora_jwt_token';

// ===== Token 管理辅助函数 =====
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

// ===== Fetch 统一请求封装 =====
interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | undefined>;
}

async function fetchClient<T>(endpoint: string, options: RequestOptions = {}): Promise<ApiResponse<T>> {
  const { params, headers, ...customConfig } = options;

  // 构建带 Query 参数的 URL
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

  // 设置请求头
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

    // 401 权限失效处理
    if (response.status === 401) {
      removeToken();
      if (typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
        window.location.href = '/login';
      }
      throw new Error('登录已过期，请重新登录');
    }

    const data: ApiResponse<T> = await response.json();

    if (!response.ok || data.code !== 200) {
      throw new Error(data.message || '请求失败，请稍后重试');
    }

    return data;
  } catch (error: any) {
    console.error(`API Error [${endpoint}]:`, error);
    throw error;
  }
}

// ===== 接口定义集合 =====

// Auth API
export const authApi = {
  register: (data: RegisterParams) => fetchClient<User>('/auth/register', { method: 'POST', body: JSON.stringify(data) }),
  login: (data: LoginParams) => fetchClient<LoginResponseData>('/auth/login', { method: 'POST', body: JSON.stringify(data) }),
};

// Topics API
export const topicsApi = {
  // 获取主题列表
  getTopics: (params?: GetTopicsParams) => {
    if (USE_MOCK) {
      return Promise.resolve({
        code: 200,
        message: 'success',
        data: {
          list: [
            {
              id: 1,
              category_id: 1,
              author_id: 1,
              author_name: 'sunny',
              title: '关于论坛架构的讨论 (Mock)',
              status: 'published',
              view_count: 42,
              reply_count: 5,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
          ],
          pagination: { page: 1, page_size: 20, total: 1 },
        },
      });
    }
    return fetchClient<PaginatedData<Topic>>('/topics', { params: params as any });
  },

  // 获取主题详情
  getTopicDetail: (id: string | number) => {
    if (USE_MOCK) {
      return Promise.resolve({
        code: 200,
        message: 'success',
        data: {
          id: Number(id),
          category_id: 1,
          author_id: 1,
          author_name: 'sunny',
          title: '关于论坛架构的讨论 (Mock)',
          content: '这是一个基于 Go + Next.js + Temporal 的论坛系统...',
          structured_content: '{}',
          status: 'published',
          view_count: 43,
          reply_count: 5,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      });
    }
    return fetchClient<Topic>(`/topics/${id}`);
  },

  // 发布主题帖
  createTopic: (data: CreateTopicParams) =>
    fetchClient<Topic>('/topics', { method: 'POST', body: JSON.stringify(data) }),
};

// Posts API
export const postsApi = {
  // 发表回复
  createPost: (topicId: string | number, data: CreatePostParams) =>
    fetchClient<Post>(`/topics/${topicId}/posts`, { method: 'POST', body: JSON.stringify(data) }),
};
