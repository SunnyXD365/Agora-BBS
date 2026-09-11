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
  const [token, setToken] = useState<string | null>(() =>
    typeof window === 'undefined' ? null : localStorage.getItem('token'),
  );
  const [isLoading, setIsLoading] = useState(() =>
    typeof window !== 'undefined' && Boolean(localStorage.getItem('token')),
  );

  const logout = useCallback(() => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    setIsLoading(false);
  }, []);

  // 初始化：自动恢复登录态
  useEffect(() => {
    if (token) {
      authApi.getMe()
        .then((res) => {
          if (res.code === 0) {
            setUser(res.data);
          } else {
            logout();
          }
        })
        .catch(() => logout())
        .finally(() => setIsLoading(false));
    }
  }, [logout, token]);

  const login = (newToken: string, newUser: UserProfile) => {
    localStorage.setItem('token', newToken);
    setToken(newToken);
    setUser(newUser);
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
