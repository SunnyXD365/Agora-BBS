'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useAuth } from '@/context/AuthContext';

const navigation = [
  ['总览', '/admin', '数据与运行态势'], ['用户管理', '/admin/users', '账号、等级与状态'],
  ['内容治理', '/admin/content', '主题与回复可见性'], ['分类管理', '/admin/categories', '板块与高争议策略'],
  ['信任流水', '/admin/trust', '积分变化审计'], ['LLM 作业', '/admin/llm', '调用质量与失败重试'],
];

export default function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname(); const router = useRouter(); const { user, isLoading, logout } = useAuth();
  useEffect(() => { if (!isLoading && user?.role !== 'admin') router.replace('/login'); }, [isLoading, router, user]);
  useEffect(() => { const handle = () => { logout(); router.replace('/login'); }; window.addEventListener('agora:admin-reverification-required', handle); return () => window.removeEventListener('agora:admin-reverification-required', handle); }, [logout, router]);
  if (isLoading) return <div className="min-h-screen bg-slate-950 p-10 text-slate-300">正在验证管理身份…</div>;
  if (user?.role !== 'admin') return null;
  return <div className="min-h-screen bg-slate-100 font-sans text-slate-900">
    <aside className="fixed inset-y-0 left-0 z-20 hidden w-64 flex-col bg-slate-950 text-slate-100 lg:flex">
      <div className="border-b border-slate-800 px-6 py-6"><p className="text-xs uppercase tracking-[0.25em] text-cyan-400">Agora BBS</p><h1 className="mt-2 text-xl font-semibold">治理管理中心</h1></div>
      <nav className="flex-1 space-y-1 p-3">{navigation.map(([label, href, hint]) => { const active = href === '/admin' ? pathname === href : pathname.startsWith(href); return <Link key={href} href={href} className={`block rounded-lg px-4 py-3 ${active ? 'bg-cyan-500 text-slate-950' : 'text-slate-300 hover:bg-slate-900 hover:text-white'}`}><strong className="block text-sm">{label}</strong><span className={`text-[11px] ${active ? 'text-slate-800' : 'text-slate-500'}`}>{hint}</span></Link>; })}</nav>
      <div className="border-t border-slate-800 p-4 text-xs"><p className="text-slate-500">当前管理员</p><p className="mt-1 truncate font-semibold">{user.username}</p><p className="truncate text-slate-400">{user.email}</p></div>
    </aside>
    <div className="lg:pl-64"><header className="sticky top-0 z-10 flex min-h-16 items-center justify-between border-b border-slate-200 bg-white/95 px-4 backdrop-blur sm:px-8"><div><p className="text-sm font-semibold">独立管理端</p><p className="text-[11px] text-emerald-700">● 邮箱二次验证已完成</p></div><div className="flex items-center gap-3"><Link href="/" className="rounded-lg border border-slate-200 px-3 py-2 text-xs hover:bg-slate-50">查看用户端</Link><button onClick={() => { logout(); router.replace('/login'); }} className="rounded-lg bg-slate-900 px-3 py-2 text-xs text-white">退出</button></div></header>
      <div className="border-b border-slate-200 bg-white px-3 py-2 lg:hidden"><nav className="flex gap-2 overflow-x-auto">{navigation.map(([label, href]) => <Link key={href} href={href} className="shrink-0 rounded bg-slate-100 px-3 py-2 text-xs">{label}</Link>)}</nav></div>
      <main className="p-4 sm:p-8">{children}</main>
    </div>
  </div>;
}
