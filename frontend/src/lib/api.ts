import {
  ApiResponse,
  User,
  Category,
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

// Categories API
export const categoriesApi = {
  // 获取板块列表
  getCategories: () => {
    if (USE_MOCK) {
      return Promise.resolve({
        code: 200,
        message: 'success',
        data: [
          { id: 1, name: '综合讨论', slug: 'general', description: '自由交流各类技术与生活话题' },
          { id: 2, name: '技术干货', slug: 'tech', description: 'Go、Next.js、云原生与架构实践' },
          { id: 3, name: '反馈与建议', slug: 'feedback', description: '社区功能建议与 Bug 反馈' },
        ],
      });
    }
    return fetchClient<Category[]>('/categories');
  },

  // 获取板块详情
  getCategory: (slug: string) => {
    if (USE_MOCK) {
      const mockCategories = [
        { id: 1, name: '综合讨论', slug: 'general', description: '自由交流各类技术与生活话题' },
        { id: 2, name: '技术干货', slug: 'tech', description: 'Go、Next.js、云原生与架构实践' },
        { id: 3, name: '反馈与建议', slug: 'feedback', description: '社区功能建议与 Bug 反馈' },
      ];
      const category = mockCategories.find((c) => c.slug === slug);
      if (!category) {
        return Promise.reject(new Error('板块不存在'));
      }
      return Promise.resolve({ code: 200, message: 'success', data: category });
    }
    return fetchClient<Category>(`/categories/${slug}`);
  },
};

// Topics API
export const topicsApi = {
  // 获取主题列表
  getTopics: (params?: GetTopicsParams) => {
    if (USE_MOCK) {
      const mockTopics = [
        {
          id: 1,
          category_id: 1,
          author_id: 1,
          author_name: 'admin',
          title: '欢迎来到 Agora-BBS 社区！',
          content: '这是一个基于 Go + Next.js + Temporal 架构构建的现代论坛系统。欢迎在此畅所欲言！',
          status: 'published',
          view_count: 102,
          reply_count: 2,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 2,
          category_id: 2,
          author_id: 2,
          author_name: 'alice',
          title: 'Go 1.23 特性解析与最佳实践',
          content: 'Go 1.23 带来了不少迭代器相关的增强，本文来聊聊如何在实际项目中使用 iterator...',
          status: 'published',
          view_count: 45,
          reply_count: 0,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 3,
          category_id: 3,
          author_id: 3,
          author_name: 'bob',
          title: '建议增加暗黑模式支持',
          content: '希望前端能支持 Dark Mode 切换，夜间看社区有点烫眼。',
          status: 'published',
          view_count: 12,
          reply_count: 1,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      ];

      // 按板块筛选
      let filtered = mockTopics;
      if (params?.category_id) {
        filtered = mockTopics.filter((t) => t.category_id === params.category_id);
      }

      // 分页
      const page = params?.page || 1;
      const pageSize = params?.page_size || 20;
      const start = (page - 1) * pageSize;
      const end = start + pageSize;
      const list = filtered.slice(start, end);

      return Promise.resolve({
        code: 200,
        message: 'success',
        data: {
          list,
          pagination: { page, page_size: pageSize, total: filtered.length },
        },
      });
    }
    return fetchClient<PaginatedData<Topic>>('/topics', { params: params as any });
  },

  // 获取主题详情
  getTopicDetail: (id: string | number) => {
    if (USE_MOCK) {
      const mockTopics = [
        {
          id: 1,
          category_id: 1,
          author_id: 1,
          author_name: 'admin',
          title: '欢迎来到 Agora-BBS 社区！',
          content: '这是一个基于 Go + Next.js + Temporal 架构构建的现代论坛系统。欢迎在此畅所欲言！',
          structured_content: '{}',
          status: 'published',
          view_count: 102,
          reply_count: 2,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 2,
          category_id: 2,
          author_id: 2,
          author_name: 'alice',
          title: 'Go 1.23 特性解析与最佳实践',
          content: 'Go 1.23 带来了不少迭代器相关的增强，本文来聊聊如何在实际项目中使用 iterator...',
          structured_content: '{}',
          status: 'published',
          view_count: 45,
          reply_count: 0,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 3,
          category_id: 3,
          author_id: 3,
          author_name: 'bob',
          title: '建议增加暗黑模式支持',
          content: '希望前端能支持 Dark Mode 切换，夜间看社区有点烫眼。',
          structured_content: '{}',
          status: 'published',
          view_count: 12,
          reply_count: 1,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      ];
      const topic = mockTopics.find((t) => t.id === Number(id));
      if (!topic) {
        return Promise.reject(new Error('帖子不存在'));
      }
      return Promise.resolve({ code: 200, message: 'success', data: topic });
    }
    return fetchClient<Topic>(`/topics/${id}`);
  },

  // 发布主题帖
  createTopic: (data: CreateTopicParams) =>
    fetchClient<Topic>('/topics', { method: 'POST', body: JSON.stringify(data) }),
};

// Posts API
export const postsApi = {
  // 获取某主题下的回复列表
  getPosts: (topicId: string | number) => {
    if (USE_MOCK) {
      const mockPosts = [
        {
          id: 1,
          topic_id: 1,
          author_id: 2,
          author_name: 'alice',
          parent_id: null,
          content: '支持！期待后续更多 AI 功能落地。',
          post_type: 'reply',
          status: 'published',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 2,
          topic_id: 1,
          author_id: 3,
          author_name: 'bob',
          parent_id: null,
          content: '环境一键拉起体验很好，赞！',
          post_type: 'reply',
          status: 'published',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 3,
          topic_id: 3,
          author_id: 1,
          author_name: 'admin',
          parent_id: null,
          content: '收到建议，暗黑模式已列入后续路线图。',
          post_type: 'reply',
          status: 'published',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
      ];
      const filtered = mockPosts.filter((p) => p.topic_id === Number(topicId));
      return Promise.resolve({ code: 200, message: 'success', data: filtered });
    }
    return fetchClient<Post[]>(`/topics/${topicId}/posts`);
  },

  // 发表回复
  createPost: (topicId: string | number, data: CreatePostParams) =>
    fetchClient<Post>(`/topics/${topicId}/posts`, { method: 'POST', body: JSON.stringify(data) }),
};
