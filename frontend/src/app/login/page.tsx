'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { authApi, getErrorMessage } from '@/services';

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [challenge, setChallenge] = useState<{ id: string; email: string; devCode?: string } | null>(null);
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const submitCredentials = async (event: React.FormEvent) => {
    event.preventDefault(); setError(''); setLoading(true);
    try {
      const res = await authApi.login({ username, password });
      if (res.data.requires_email_verification && res.data.challenge_id) {
        setChallenge({ id: res.data.challenge_id, email: res.data.masked_email || '管理员邮箱', devCode: res.data.development_verification_code });
        return;
      }
      if (!res.data.token || !res.data.user) throw new Error('登录响应不完整');
      login(res.data.token, res.data.user);
      router.replace(res.data.user.role === 'admin' ? '/admin' : '/');
    } catch (err: unknown) { setError(getErrorMessage(err, '登录失败，请检查用户名或密码')); }
    finally { setLoading(false); }
  };

  const verifyCode = async (event: React.FormEvent) => {
    event.preventDefault(); if (!challenge) return; setError(''); setLoading(true);
    try {
      const res = await authApi.verifyAdminEmail({ challenge_id: challenge.id, code });
      login(res.data.token, res.data.user); router.replace('/admin');
    } catch (err: unknown) { setError(getErrorMessage(err, '验证码错误或已过期')); }
    finally { setLoading(false); }
  };

  return <div className="flex min-h-[calc(100vh-12rem)] items-center justify-center">
    <div className="paper-card w-full max-w-md rounded-2xl p-8">
      <p className="text-center text-xs font-semibold uppercase tracking-[0.24em] text-[var(--accent-ink)]">统一身份入口</p>
      <h1 className="mt-2 text-center text-2xl font-bold">{challenge ? '管理员邮箱验证' : '登录 Agora BBS'}</h1>
      <p className="mt-2 text-center text-xs text-[var(--text-muted)]">{challenge ? `验证码已发送至 ${challenge.email}` : '普通用户进入社区，管理员验证邮箱后进入管理端'}</p>
      {error && <div className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-xs text-red-700">{error}</div>}
      {challenge?.devCode && <div className="mt-4 rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-800">开发环境演示验证码：<strong className="tracking-widest">{challenge.devCode}</strong></div>}
      {!challenge ? <form onSubmit={submitCredentials} className="mt-6 space-y-4">
        <Field label="用户名"><input required value={username} onChange={(e) => setUsername(e.target.value)} className="mt-1 w-full rounded-lg border p-2.5 text-sm" placeholder="请输入用户名" /></Field>
        <Field label="密码"><input type="password" required value={password} onChange={(e) => setPassword(e.target.value)} className="mt-1 w-full rounded-lg border p-2.5 text-sm" placeholder="请输入密码" /></Field>
        <button disabled={loading} className="paper-btn-primary w-full rounded-lg py-2.5 text-sm font-medium disabled:opacity-50">{loading ? '登录中…' : '登录'}</button>
      </form> : <form onSubmit={verifyCode} className="mt-6 space-y-4">
        <Field label="6 位验证码"><input inputMode="numeric" pattern="[0-9]{6}" maxLength={6} required value={code} onChange={(e) => setCode(e.target.value.replace(/\D/g, ''))} className="mt-1 w-full rounded-lg border p-3 text-center text-xl tracking-[0.5em]" autoFocus /></Field>
        <button disabled={loading || code.length !== 6} className="w-full rounded-lg bg-slate-900 py-2.5 text-sm font-medium text-white disabled:opacity-50">{loading ? '验证中…' : '验证并进入管理端'}</button>
        <button type="button" onClick={() => { setChallenge(null); setCode(''); setError(''); }} className="w-full text-xs text-[var(--text-muted)] hover:underline">返回重新登录</button>
      </form>}
      {!challenge && <p className="mt-6 text-center text-xs text-[var(--text-muted)]">还没有账号？ <Link href="/register" className="font-semibold text-[var(--accent-ink)] hover:underline">立即注册</Link></p>}
    </div>
  </div>;
}

function Field({ label, children }: { label: string; children: React.ReactNode }) { return <label className="block text-xs font-medium">{label}{children}</label>; }
