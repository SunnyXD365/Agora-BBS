'use client';

import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { UserProfile } from '@/types/api';
import { authApi } from '@/services';

interface AuthContextType {
  user: UserProfile | null;
  token: string | null;
  isLoading: boolean;
  login: (token: string, user: UserProfile) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    setIsLoading(false);
  }, []);

  // 在挂载后读取浏览器存储，避免服务端渲染与首次 hydration 的状态不一致。
  useEffect(() => {
    let cancelled = false;
    const restoreSession = async () => {
      // 让状态更新发生在异步恢复流程中，而不是 effect 的同步执行阶段。
      await Promise.resolve();
      if (cancelled) return;
      const storedToken = localStorage.getItem('token');
      if (!storedToken) {
        setIsLoading(false);
        return;
      }

      try {
        const res = await authApi.getMe();
        if (cancelled) return;
        if (res.code === 0) {
          setToken(storedToken);
          setUser(res.data);
        } else {
          logout();
        }
      } catch {
        if (!cancelled) logout();
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    };
    void restoreSession();
    return () => { cancelled = true; };
  }, [logout]);

  const login = (newToken: string, newUser: UserProfile) => {
    localStorage.setItem('token', newToken);
    setToken(newToken);
    setUser(newUser);
    setIsLoading(false);
  };

  return (
    <AuthContext.Provider value={{ user, token, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
