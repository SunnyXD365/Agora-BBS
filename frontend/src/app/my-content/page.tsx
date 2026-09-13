'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Suspense } from 'react';
import Pagination from '@/components/Pagination';
import { useAuth } from '@/context/AuthContext';
import { contentApi, getErrorMessage, topicApi } from '@/services';
import { MyContentItem } from '@/types/api';

const statusLabels: Record<string, string> = {
  draft: '草稿', cooling: '冷静期', pending_review: '待盲审', published: '已发布', rejected: '未通过', recalled: '已撤回', hidden: '已隐藏',
};

export default function MyContentPage() {
  return <Suspense fallback={<div className="h-48 animate-pulse rounded-xl bg-stone-200" />}><MyContent /></Suspense>;
}

function MyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { user, isLoading } = useAuth();
  const [contentType, setContentType] = useState<'topic' | 'post'>('topic');
  const [status, setStatus] = useState(searchParams.get('status') || '');
  const [items, setItems] = useState<MyContentItem[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isLoading && !user) router.replace('/login');
  }, [isLoading, router, user]);

  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    contentApi.mine({ type: contentType, status: status || undefined, page, page_size: pageSize })
      .then((result) => { if (!cancelled) { setItems(result.data.items); setTotal(result.data.total); } })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '内容加载失败')); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, [contentType, page, pageSize, status, user]);

  const changeType = (type: 'topic' | 'post') => { setLoading(true); setError(''); setContentType(type); setStatus(''); setPage(1); };
  const deleteDraft = async (item: MyContentItem) => {
    if (!window.confirm(`确定删除草稿“${item.title || '无标题草稿'}”吗？删除后无法恢复。`)) return;
    try {
      await topicApi.deleteDraft(item.id);
      if (items.length === 1 && page > 1) setPage(page - 1);
      else {
        const result = await contentApi.mine({ type: contentType, status: status || undefined, page, page_size: pageSize });
        setItems(result.data.items); setTotal(result.data.total);
      }
    } catch (err: unknown) { setError(getErrorMessage(err, '草稿删除失败')); }
  };

  if (isLoading) return <div className="h-48 animate-pulse rounded-xl bg-stone-200" />;
  if (!user) return null;
  const statusOptions = contentType === 'topic'
    ? ['', 'draft', 'cooling', 'pending_review', 'published', 'rejected', 'recalled', 'hidden']
    : ['', 'cooling', 'published', 'recalled', 'hidden'];

  return (
    <div className="mx-auto max-w-5xl space-y-5">
      <header className="flex flex-wrap items-end justify-between gap-3"><div><p className="text-xs font-semibold tracking-widest text-[var(--accent-ink)]">个人创作台</p><h1 className="mt-1 text-2xl font-bold">我的内容</h1><p className="mt-2 text-sm text-[var(--text-muted)]">查看主题、回复和草稿，跟踪冷静期与盲审状态。</p></div><Link href="/topics/new" className="paper-btn-primary rounded-lg px-5 py-2 text-sm">写新主题</Link></header>
      <div className="paper-card flex flex-wrap items-center gap-3 rounded-xl p-4">
        <div className="flex rounded-lg bg-stone-100 p-1"><button onClick={() => changeType('topic')} className={`rounded-md px-4 py-2 text-sm ${contentType === 'topic' ? 'bg-white font-bold shadow-sm' : ''}`}>我的主题</button><button onClick={() => changeType('post')} className={`rounded-md px-4 py-2 text-sm ${contentType === 'post' ? 'bg-white font-bold shadow-sm' : ''}`}>我的回复</button></div>
        <label className="ml-auto flex items-center gap-2 text-sm text-[var(--text-muted)]">状态<select aria-label="内容状态" value={status} onChange={(event) => { setLoading(true); setError(''); setStatus(event.target.value); setPage(1); }} className="rounded-lg border px-3 py-2"><option value="">全部状态</option>{statusOptions.filter(Boolean).map((value) => <option key={value} value={value}>{statusLabels[value]}</option>)}</select></label>
      </div>
      {error && <p className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</p>}
      {loading ? <div className="h-48 animate-pulse rounded-xl bg-stone-200" /> : items.length === 0 ? <div className="paper-card rounded-xl p-12 text-center"><p className="font-bold">这里暂时没有内容</p><p className="mt-2 text-sm text-[var(--text-muted)]">保存的草稿和发表过的内容都会集中显示在这里。</p></div> : <div className="space-y-3">{items.map((item) => <ContentCard key={`${item.type}-${item.id}`} item={item} onDeleteDraft={deleteDraft} />)}</div>}
      {!loading && <Pagination page={page} pageSize={pageSize} total={total} onPageChange={(nextPage) => { setLoading(true); setPage(nextPage); }} onPageSizeChange={(size) => { setLoading(true); setPage(1); setPageSize(size); }} itemLabel={contentType === 'topic' ? '个主题' : '条回复'} />}
    </div>
  );
}

function ContentCard({ item, onDeleteDraft }: { item: MyContentItem; onDeleteDraft: (item: MyContentItem) => void }) {
  const visibleTopic = !['recalled', 'hidden'].includes(item.status);
  return <article className="paper-card rounded-xl p-5"><div className="flex flex-wrap items-start justify-between gap-3"><div className="min-w-0 flex-1"><div className="flex flex-wrap items-center gap-2"><span className={`rounded-full px-2.5 py-1 text-xs ${item.status === 'published' ? 'bg-emerald-100 text-emerald-800' : item.status === 'draft' ? 'bg-stone-200 text-stone-700' : 'bg-amber-100 text-amber-800'}`}>{statusLabels[item.status] || item.status}</span>{item.type === 'post' && <span className="text-xs text-[var(--text-muted)]">回复于《{item.title}》</span>}</div><h2 className="mt-3 font-bold">{item.type === 'topic' ? (item.title || '无标题草稿') : item.excerpt}</h2>{item.type === 'topic' && <p className="mt-2 line-clamp-2 text-sm leading-6 text-[var(--text-muted)]">{item.excerpt || '尚未填写核心观点'}</p>}<p className="mt-3 text-xs text-[var(--text-muted)]">最后更新 {new Date(item.updated_at).toLocaleString()}</p></div><div className="flex shrink-0 gap-3 text-sm">{item.status === 'draft' ? <><Link href={`/topics/new?draft=${item.id}`} className="font-bold text-[var(--accent-ink)] hover:underline">继续编辑</Link><button onClick={() => onDeleteDraft(item)} className="text-red-700 hover:underline">删除</button></> : visibleTopic && <Link href={`/topics/${item.topic_id}`} className="font-bold text-[var(--accent-ink)] hover:underline">查看</Link>}</div></div></article>;
}
