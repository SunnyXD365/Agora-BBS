'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';

export default function Navbar() {
  const { user, logout, isLoading } = useAuth();
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  return (
    <header className="border-b bg-white shadow-sm">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-6">
          <Link href="/" className="text-xl font-bold text-gray-900 hover:opacity-80">
            Agora BBS
          </Link>
        </div>

        <nav className="flex items-center gap-4">
          <Link href="/guide" className="text-sm text-gray-600 hover:text-gray-900">使用说明</Link>
          {isLoading ? (
            <div className="h-8 w-24 animate-pulse rounded bg-gray-200" />
          ) : user ? (
            <div className="flex items-center gap-4">
              <Link href="/bookmarks" className="text-sm text-gray-600 hover:text-gray-900">收藏</Link>
              {user.capabilities.includes('review') && <Link href="/reviews" className="text-sm text-gray-600 hover:text-gray-900">匿名盲审</Link>}
              {user.role === 'admin' && <Link href="/admin" className="text-sm font-semibold text-red-700 hover:text-red-900">管理后台</Link>}
              <Link href="/profile" className="text-sm text-gray-600 hover:text-gray-900">成长中心</Link>
              <span className="text-sm text-gray-600">
                你好，<strong className="text-gray-900">{user.username}</strong>
              </span>
              <button
                onClick={handleLogout}
                className="rounded-md bg-gray-100 px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-200 transition-colors"
              >
                退出登录
              </button>
            </div>
          ) : (
            <div className="flex items-center gap-3">
              <Link
                href="/login"
                className="text-sm font-medium text-gray-700 hover:text-gray-900 px-3 py-1.5"
              >
                登录
              </Link>
              <Link
                href="/register"
                className="rounded-md bg-blue-600 px-3.5 py-1.5 text-sm font-medium text-white hover:bg-blue-500 transition-colors"
              >
                注册
              </Link>
            </div>
          )}
        </nav>
      </div>
    </header>
  );
}
