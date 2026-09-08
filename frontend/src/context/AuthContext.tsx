'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, LoginParams, RegisterParams } from '@/types/api';
import { authApi, setToken, removeToken, getToken } from '@/lib/api';
import { useRouter } from 'next/navigation';

const USER_KEY = 'agora_user_info';

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (params: LoginParams) => Promise<void>;
  register: (params: RegisterParams) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  // 初始化检查本地用户数据
  useEffect(() => {
    const token = getToken();
    const savedUser = localStorage.getItem(USER_KEY);
    if (token && savedUser) {
      try {
        setUser(JSON.parse(savedUser));
      } catch (e) {
        console.error('Failed to parse cached user:', e);
        removeToken();
        localStorage.removeItem(USER_KEY);
      }
    }
    setLoading(false);
  }, []);

  // 登录
  const login = async (params: LoginParams) => {
    const res = await authApi.login(params);
    const { token, user: userData } = res.data;
    
    setToken(token);
    localStorage.setItem(USER_KEY, JSON.stringify(userData));
    setUser(userData);
    router.push('/');
  };

  // 注册（注册成功后直接自动登录或引导跳转）
  const register = async (params: RegisterParams) => {
    await authApi.register(params);
    // 注册完成后自动登录
    await login({ username: params.username, password: params.password });
  };

  // 退出登录
  const logout = () => {
    removeToken();
    localStorage.removeItem(USER_KEY);
    setUser(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
