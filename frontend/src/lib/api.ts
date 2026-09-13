import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器：自动注入 JWT
api.interceptors.request.use(
  (config) => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器：提取 data payload，捕获 401
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const requiresAdminReverification = error.response?.data?.code === 40391;
    if ((error.response?.status === 401 || requiresAdminReverification) && typeof window !== 'undefined') {
      localStorage.removeItem('token');
      if (requiresAdminReverification && window.location.pathname.startsWith('/admin')) {
        window.dispatchEvent(new Event('agora:admin-reverification-required'));
      }
    }
    return Promise.reject(error.response?.data || error);
  }
);

export default api;
