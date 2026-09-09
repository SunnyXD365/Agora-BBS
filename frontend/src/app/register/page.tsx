'use client';

import { useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import Link from 'next/link';

export default function RegisterPage() {
  const { register } = useAuth();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);

    try {
      await register({ username, email, password });
    } catch (err: any) {
      setError(err.message || '注册失败，请稍后重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-h-[80vh] flex items-center justify-center">
      <div className="w-full max-w-md p-6 bg-[#141414] border border-gray-800 rounded-lg">
        <h1 className="text-2xl font-bold text-center mb-6 text-gray-200">注册账号</h1>

        {error && (
          <div className="mb-4 p-3 bg-red-900/30 text-red-400 text-sm rounded-md border border-red-800">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">用户名</label>
            <input
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-gray-200 placeholder-gray-500"
              placeholder="设置用户名"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">邮箱</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-gray-200 placeholder-gray-500"
              placeholder="your@email.com"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">密码</label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-gray-200 placeholder-gray-500"
              placeholder="设置登录密码"
            />
          </div>

          <button
            type="submit"
            disabled={submitting}
            className="w-full py-2 bg-gray-700 text-gray-100 rounded-md hover:bg-gray-600 font-medium transition disabled:bg-gray-800 disabled:text-gray-500"
          >
            {submitting ? '注册中...' : '确认注册'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-gray-500">
          已有账号？{' '}
          <Link href="/login" className="text-gray-300 hover:text-gray-100 underline">
            直接登录
          </Link>
        </p>
      </div>
    </div>
  );
}
