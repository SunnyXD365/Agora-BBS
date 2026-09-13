'use client';

import Link from 'next/link';
import { FormEvent, Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Pagination from '@/components/Pagination';
import { contentApi, getErrorMessage } from '@/services';
import { SearchResult } from '@/types/api';

export default function SearchPage() {
  return <Suspense fallback={<div className="h-48 animate-pulse rounded-xl bg-stone-200" />}><SearchResults /></Suspense>;
}

function SearchResults() {
  const searchParams = useSearchParams();
  const query = (searchParams.get('q') || '').trim();
  return <SearchResultsForQuery key={query} query={query} />;
}

function SearchResultsForQuery({ query }: { query: string }) {
  const router = useRouter();
  const [input, setInput] = useState(query);
  const [type, setType] = useState<'' | 'topic' | 'post'>('');
  const [items, setItems] = useState<SearchResult[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(query.length >= 2);
  const [error, setError] = useState('');

  useEffect(() => {
    if (query.length < 2) return;
    let cancelled = false;
    contentApi.search({ q: query, type: type || undefined, page, page_size: pageSize })
      .then((result) => { if (!cancelled) { setItems(result.data.items); setTotal(result.data.total); } })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '搜索失败')); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, [page, pageSize, query, type]);

  const submit = (event: FormEvent) => {
    event.preventDefault();
    const value = input.trim();
    if (value.length < 2) { setError('请输入至少 2 个字进行搜索。'); return; }
    router.push(`/search?q=${encodeURIComponent(value)}`);
  };

  return <div className="mx-auto max-w-5xl space-y-5"><header><p className="text-xs font-semibold tracking-widest text-[var(--accent-ink)]">全站检索</p><h1 className="mt-1 text-2xl font-bold">搜索论坛内容</h1></header><form onSubmit={submit} className="paper-card flex rounded-xl p-3"><input aria-label="搜索关键词" minLength={2} maxLength={100} value={input} onChange={(event) => setInput(event.target.value)} placeholder="输入主题、观点或回复中的关键词" className="min-w-0 flex-1 rounded-l-lg border p-3" /><button className="paper-btn-primary rounded-r-lg px-6">搜索</button></form>{query.length >= 2 && <div className="flex flex-wrap items-center justify-between gap-3"><p className="text-sm text-[var(--text-muted)]">“{query}” 找到 {total} 条结果</p><div className="flex rounded-lg bg-stone-100 p-1">{([['', '全部'], ['topic', '主题'], ['post', '回复']] as const).map(([value, label]) => <button key={value} onClick={() => { setLoading(true); setError(''); setType(value); setPage(1); }} className={`rounded-md px-4 py-2 text-sm ${type === value ? 'bg-white font-bold shadow-sm' : ''}`}>{label}</button>)}</div></div>}{error && <p className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</p>}{loading ? <div className="h-48 animate-pulse rounded-xl bg-stone-200" /> : query.length < 2 ? <div className="paper-card rounded-xl p-12 text-center text-[var(--text-muted)]">输入至少 2 个字，可以模糊搜索主题和回复。</div> : items.length === 0 ? <div className="paper-card rounded-xl p-12 text-center"><p className="font-bold">没有找到相关内容</p><p className="mt-2 text-sm text-[var(--text-muted)]">可以缩短关键词，或换一个相近说法再试。</p></div> : <div className="space-y-3">{items.map((item) => <Link key={`${item.type}-${item.id}`} href={`/topics/${item.topic_id}`} className="paper-card block rounded-xl p-5 transition hover:-translate-y-0.5 hover:border-stone-400"><div className="flex items-center gap-2 text-xs text-[var(--text-muted)]"><span className="rounded bg-stone-100 px-2 py-1">{item.type === 'topic' ? '主题' : '回复'}</span><span>{item.author_name}</span><span>·</span><time>{new Date(item.created_at).toLocaleDateString()}</time></div><h2 className="mt-3 font-bold">{item.title}</h2><p className="mt-2 line-clamp-3 text-sm leading-6 text-[var(--text-muted)]">{item.excerpt}</p></Link>)}</div>}{!loading && <Pagination page={page} pageSize={pageSize} total={total} onPageChange={(nextPage) => { setLoading(true); setPage(nextPage); }} onPageSizeChange={(size) => { setLoading(true); setPage(1); setPageSize(size); }} itemLabel="条结果" />}</div>;
}
