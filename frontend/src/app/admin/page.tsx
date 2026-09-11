'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/context/AuthContext';
import { adminApi, getErrorMessage } from '@/services';
import { AdminContent, AdminLLMJob, AdminOverview, AdminTrustLog, AdminUser, Category } from '@/types/api';

export default function AdminPage() {
  const { user, isLoading } = useAuth();
  const [overview, setOverview] = useState<AdminOverview | null>(null);
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [contents, setContents] = useState<AdminContent[]>([]);
  const [jobs, setJobs] = useState<AdminLLMJob[]>([]);
  const [logs, setLogs] = useState<AdminTrustLog[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const [newCategory, setNewCategory] = useState({ name: '', slug: '', description: '' });

  const load = async () => {
    setBusy(true); setError('');
    try {
      const [overviewRes, usersRes, topicRes, postRes, jobsRes, logsRes, categoriesRes] = await Promise.all([
        adminApi.overview(), adminApi.users(), adminApi.contents('topic'), adminApi.contents('post'), adminApi.llmJobs(), adminApi.trustLogs(), adminApi.categories(),
      ]);
      setOverview(overviewRes.data); setUsers(usersRes.data.items); setContents([...topicRes.data.items, ...postRes.data.items]);
      setJobs(jobsRes.data.items); setLogs(logsRes.data.items); setCategories(categoriesRes.data);
    } catch (err: unknown) { setError(getErrorMessage(err, '管理数据加载失败')); }
    finally { setBusy(false); }
  };

  useEffect(() => {
    if (isLoading || user?.role !== 'admin') return;
    let cancelled = false;
    Promise.all([adminApi.overview(), adminApi.users(), adminApi.contents('topic'), adminApi.contents('post'), adminApi.llmJobs(), adminApi.trustLogs(), adminApi.categories()])
      .then(([overviewRes, usersRes, topicRes, postRes, jobsRes, logsRes, categoriesRes]) => {
        if (cancelled) return;
        setOverview(overviewRes.data); setUsers(usersRes.data.items); setContents([...topicRes.data.items, ...postRes.data.items]);
        setJobs(jobsRes.data.items); setLogs(logsRes.data.items); setCategories(categoriesRes.data);
      })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '管理数据加载失败')); });
    return () => { cancelled = true; };
  }, [isLoading, user]);

  const run = async (action: () => Promise<unknown>) => { setBusy(true); setError(''); try { await action(); await load(); } catch (err: unknown) { setError(getErrorMessage(err, '治理操作失败')); setBusy(false); } };
  const saveCategory = (category: Category, changes: Partial<Category>) => run(() => adminApi.updateCategory(category.id, { name: category.name, slug: category.slug, description: category.description, sort_order: category.sort_order, is_active: category.is_active, requires_review: category.requires_review, ...changes }));
  const createCategory = async (event: React.FormEvent) => {
    event.preventDefault();
    await run(() => adminApi.createCategory({ ...newCategory, sort_order: categories.length, is_active: true, requires_review: false }));
    setNewCategory({ name: '', slug: '', description: '' });
  };

  if (isLoading) return <div className="mx-auto h-48 max-w-7xl animate-pulse rounded-xl bg-stone-200" />;
  if (user?.role !== 'admin') return <div className="paper-card mx-auto max-w-3xl rounded-xl p-10 text-center">此页面仅管理员可访问。</div>;

  const cards = overview ? [
    ['用户', overview.users_total], ['主题', overview.topics_total], ['回复', overview.posts_total], ['7 日活跃', overview.active_users_7_days],
    ['盲审完成', `${overview.review_completed}/${overview.review_total}`], ['LLM 成功', `${overview.llm_success}/${overview.llm_calls}`], ['平均延迟', `${Math.round(overview.llm_average_ms)} ms`], ['Token', overview.llm_prompt_tokens + overview.llm_output_tokens],
  ] : [];

  return (
    <div className="mx-auto max-w-7xl space-y-8 text-sm">
      <header className="flex items-center justify-between"><div><h1 className="text-2xl font-bold">社区治理后台</h1><p className="mt-1 text-[var(--text-muted)]">统计观察与轻量治理；管理员不能修改用户正文。</p></div><button disabled={busy} onClick={() => void load()} className="rounded border px-3 py-2 disabled:opacity-50">刷新</button></header>
      {error && <div className="rounded border border-red-200 bg-red-50 p-3 text-red-700">{error}</div>}
      <section className="grid grid-cols-2 gap-3 md:grid-cols-4">{cards.map(([label, value]) => <div key={label} className="paper-card rounded-xl p-4"><p className="text-xs text-[var(--text-muted)]">{label}</p><p className="mt-2 text-2xl font-bold">{value}</p></div>)}</section>

      {overview && <section className="paper-card rounded-xl p-5"><h2 className="font-bold">近 7 日新增趋势</h2><div className="mt-4 grid grid-cols-7 items-end gap-2">{overview.trend.map((day) => { const total = day.users + day.topics + day.posts + day.feedback; return <div key={day.date} className="text-center"><div title={`用户 ${day.users} / 主题 ${day.topics} / 回复 ${day.posts} / 反馈 ${day.feedback}`} style={{ height: `${Math.max(8, Math.min(120, total * 8))}px` }} className="mx-auto w-8 rounded-t bg-[var(--accent-ink)] opacity-75" /><p className="mt-1 text-[10px] text-[var(--text-muted)]">{day.date.slice(5)}</p></div>; })}</div><div className="mt-4 flex flex-wrap gap-4 text-xs"><span>内容状态：{Object.entries(overview.content_status).map(([key, value]) => `${key} ${value}`).join(' · ') || '无'}</span><span>信任分布：{Object.entries(overview.trust_distribution).map(([key, value]) => `${key} ${value}`).join(' · ') || '无'}</span></div></section>}

      <AdminSection title="用户与信任">
        <div className="overflow-x-auto"><table className="w-full text-left"><thead><tr><Th>用户</Th><Th>角色/状态</Th><Th>等级</Th><Th>信任分</Th><Th>有效阅读</Th><Th>操作</Th></tr></thead><tbody>{users.map((item) => <tr key={item.id} className="border-t"><Td>{item.username}<div className="text-xs text-[var(--text-muted)]">{item.email}</div></Td><Td>{item.role} / {item.status}</Td><Td>L{item.unlock_level}</Td><Td>{item.trust_score}</Td><Td>{item.verified_read_seconds}s</Td><Td>{item.id !== user.id && <button disabled={busy} onClick={() => void run(() => adminApi.setUserStatus(item.id, item.status === 'active' ? 'suspended' : 'active'))} className="underline">{item.status === 'active' ? '停用' : '恢复'}</button>}</Td></tr>)}</tbody></table></div>
        <h3 className="mt-6 font-semibold">最近信任流水</h3><div className="mt-2 space-y-2">{logs.slice(0, 20).map((log) => <div key={log.id} className="rounded bg-stone-50 p-2"><strong>{log.username}</strong> <span className={log.score_delta >= 0 ? 'text-emerald-700' : 'text-red-700'}>{log.score_delta >= 0 ? '+' : ''}{log.score_delta}</span> · {log.event_type}<span className="ml-2 text-xs text-[var(--text-muted)]">{log.reason}</span></div>)}</div>
      </AdminSection>

      <AdminSection title="内容状态与可见性">
        <div className="space-y-2">{contents.map((item) => <div key={`${item.type}-${item.id}`} className="flex items-center justify-between gap-4 rounded border p-3"><div><strong>{item.type === 'topic' ? '主题' : '回复'} · {item.title}</strong><p className="mt-1 line-clamp-1 text-xs text-[var(--text-muted)]">{item.author_name} · {item.status} · {item.excerpt}</p></div>{['published', 'hidden'].includes(item.status) && <button disabled={busy} onClick={() => void run(() => adminApi.setContentVisibility(item.type, item.id, item.status !== 'hidden'))} className="shrink-0 underline">{item.status === 'hidden' ? '恢复' : '隐藏'}</button>}</div>)}</div>
      </AdminSection>

      <AdminSection title="分类管理">
        <form onSubmit={createCategory} className="mb-4 grid gap-2 md:grid-cols-4"><input required minLength={2} value={newCategory.name} onChange={(event) => setNewCategory({ ...newCategory, name: event.target.value })} placeholder="分类名称" className="rounded border p-2" /><input required minLength={2} value={newCategory.slug} onChange={(event) => setNewCategory({ ...newCategory, slug: event.target.value })} placeholder="slug" className="rounded border p-2" /><input value={newCategory.description} onChange={(event) => setNewCategory({ ...newCategory, description: event.target.value })} placeholder="说明" className="rounded border p-2" /><button disabled={busy} className="paper-btn-primary rounded p-2">新增分类</button></form>
        <div className="space-y-2">{categories.map((category) => <div key={category.id} className="flex flex-wrap items-center justify-between gap-2 rounded border p-3"><span><strong>{category.name}</strong> / {category.slug} · 排序 {category.sort_order}</span><div className="space-x-3"><button disabled={busy} onClick={() => void saveCategory(category, { requires_review: !category.requires_review })} className="underline">{category.requires_review ? '取消高争议' : '设为高争议'}</button><button disabled={busy} onClick={() => void saveCategory(category, { is_active: !category.is_active })} className="underline">{category.is_active ? '停用' : '启用'}</button></div></div>)}</div>
      </AdminSection>

      <AdminSection title="LLM 作业与失败重试">
        <div className="overflow-x-auto"><table className="w-full text-left"><thead><tr><Th>类型</Th><Th>状态</Th><Th>模型</Th><Th>延迟</Th><Th>Token</Th><Th>错误/操作</Th></tr></thead><tbody>{jobs.map((job) => <tr key={job.id} className="border-t"><Td>{job.job_type}</Td><Td>{job.status}</Td><Td>{job.model}</Td><Td>{job.latency_ms}ms</Td><Td>{job.prompt_tokens + job.completion_tokens}</Td><Td><span className="text-xs text-red-700">{job.error_message}</span>{job.status === 'failed' && <button disabled={busy} onClick={() => void run(() => adminApi.retryLLMJob(job.id))} className="ml-2 underline">安全重试</button>}</Td></tr>)}</tbody></table></div>
      </AdminSection>
    </div>
  );
}

function AdminSection({ title, children }: { title: string; children: React.ReactNode }) { return <section className="paper-card rounded-xl p-5"><h2 className="mb-4 font-bold">{title}</h2>{children}</section>; }
function Th({ children }: { children: React.ReactNode }) { return <th className="p-2 text-xs text-[var(--text-muted)]">{children}</th>; }
function Td({ children }: { children: React.ReactNode }) { return <td className="p-2 align-top">{children}</td>; }
