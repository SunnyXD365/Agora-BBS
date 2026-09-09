'use client';

import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import CategoryNav from './CategoryNav';

export default function Header() {
  const { user, logout, loading } = useAuth();

  return (
    <header className="border-b border-gray-800 bg-[#0a0a0a]">
      <div className="max-w-5xl mx-auto px-4 h-16 flex items-center justify-between">
        {/* Logo */}
        <Link href="/" className="text-xl font-bold text-gray-100 hover:text-gray-300 transition">
          Agora BBS
        </Link>

        {/* 板块导航 */}
        <div className="hidden md:block">
          <CategoryNav />
        </div>

        {/* 导航区域 */}
        <div className="flex items-center space-x-4">
          {loading ? (
            <div className="text-sm text-gray-500">加载中...</div>
          ) : user ? (
            <>
              <Link
                href="/topics/new"
                className="px-3 py-1.5 bg-gray-700 text-gray-100 text-sm font-medium rounded-md hover:bg-gray-600 transition"
              >
                + 发布新帖
              </Link>
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <span className="font-semibold">{user.username}</span>
                <span className="text-xs px-2 py-0.5 bg-gray-800 rounded-full text-gray-400">
                  声望: {user.trust_score}
                </span>
              </div>
              <button
                onClick={logout}
                className="text-sm text-gray-400 hover:text-red-400 transition"
              >
                退出登录
              </button>
            </>
          ) : (
            <>
              <Link
                href="/login"
                className="text-sm text-gray-300 hover:text-gray-100 transition font-medium"
              >
                登录
              </Link>
              <Link
                href="/register"
                className="px-3 py-1.5 bg-gray-700 text-gray-100 text-sm font-medium rounded-md hover:bg-gray-600 transition"
              >
                注册
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
