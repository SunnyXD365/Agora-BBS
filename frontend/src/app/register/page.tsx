'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { authApi, getErrorMessage } from '@/services';

export default function RegisterPage() {
  const router = useRouter();
  const { login } = useAuth();

  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError('两次输入的密码不一致');
      return;
    }

    setLoading(true);

    try {
      const res = await authApi.register({ username, email, password });
      if (res.code === 0 && res.data?.token && res.data?.user) {
        // 注册成功自动完成登录并跳转到首页
        login(res.data.token, res.data.user);
        router.push('/');
      } else {
        setError(res.msg || '注册失败');
      }
    } catch (err: unknown) {
      setError(getErrorMessage(err, '注册失败，请检查填写内容'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-[calc(100vh-12rem)] items-center justify-center">
      <div className="paper-card w-full max-w-md rounded-2xl p-8">
        <p className="text-center text-xs font-semibold uppercase tracking-[0.24em] text-[var(--accent-ink)]">
          统一身份入口
        </p>
        <h1 className="mt-2 text-center text-2xl font-bold">注册 Agora BBS</h1>
        <p className="mt-2 text-center text-xs text-[var(--text-muted)]">
          创建新账号，开启理性讨论与社区治理体验
        </p>

        {error && (
          <div className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <Field label="用户名">
            <input
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="mt-1 w-full rounded-lg border p-2.5 text-sm"
              placeholder="设置你的用户名"
            />
          </Field>

          <Field label="邮箱">
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="mt-1 w-full rounded-lg border p-2.5 text-sm"
              placeholder="your@email.com"
            />
          </Field>

          <Field label="密码">
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="mt-1 w-full rounded-lg border p-2.5 text-sm"
              placeholder="设置登录密码"
            />
          </Field>

          <Field label="确认密码">
            <input
              type="password"
              required
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="mt-1 w-full rounded-lg border p-2.5 text-sm"
              placeholder="再次输入密码"
            />
          </Field>

          <button
            disabled={loading}
            className="paper-btn-primary w-full rounded-lg py-2.5 text-sm font-medium disabled:opacity-50"
          >
            {loading ? '注册中…' : '注册并登录'}
          </button>
        </form>

        <p className="mt-6 text-center text-xs text-[var(--text-muted)]">
          已有账号？{' '}
          <Link href="/login" className="font-semibold text-[var(--accent-ink)] hover:underline">
            直接登录
          </Link>
        </p>
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block text-xs font-medium">
      {label}
      {children}
    </label>
  );
}