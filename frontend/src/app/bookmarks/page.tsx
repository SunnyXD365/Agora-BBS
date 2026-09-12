'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import TopicCard from '@/components/TopicCard';
import { bookmarkApi, getErrorMessage } from '@/services';
import { Bookmark } from '@/types/api';
import { useAuth } from '@/context/AuthContext';
import Pagination from '@/components/Pagination';

export default function BookmarksPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [items, setItems] = useState<Bookmark[]>([]);
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
    bookmarkApi.list({ page, page_size: pageSize }).then((result) => { if (!cancelled) { setItems(result.data.items); setTotal(result.data.total); } })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '收藏加载失败')); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, [page, pageSize, user]);

  if (isLoading) return <div className="h-40 animate-pulse rounded-xl bg-stone-200" />;
  if (!user) return null;
  return <div className="mx-auto max-w-4xl space-y-4"><h1 className="text-2xl font-bold">我的收藏</h1>{error && <p className="text-red-700">{error}</p>}{loading ? <div className="h-32 animate-pulse rounded-xl bg-stone-200" /> : items.length ? items.map((item) => <TopicCard key={item.id} topic={item.topic} />) : <div className="paper-card rounded-xl p-10 text-center">暂未收藏主题。</div>}{!loading && <Pagination page={page} pageSize={pageSize} total={total} onPageChange={(nextPage) => { setLoading(true); setPage(nextPage); }} onPageSizeChange={(size) => { setLoading(true); setPage(1); setPageSize(size); }} itemLabel="个收藏" />}</div>;
}
