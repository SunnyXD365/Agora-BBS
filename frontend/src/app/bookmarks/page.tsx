'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import TopicCard from '@/components/TopicCard';
import { bookmarkApi, getErrorMessage } from '@/services';
import { Bookmark } from '@/types/api';
import { useAuth } from '@/context/AuthContext';

export default function BookmarksPage() {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [items, setItems] = useState<Bookmark[]>([]);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isLoading && !user) router.replace('/login');
  }, [isLoading, router, user]);

  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    bookmarkApi.list().then((result) => { if (!cancelled) setItems(result.data.items); })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '收藏加载失败')); });
    return () => { cancelled = true; };
  }, [user]);

  if (isLoading) return <div className="h-40 animate-pulse rounded-xl bg-stone-200" />;
  if (!user) return null;
  return <div className="mx-auto max-w-4xl space-y-4"><h1 className="text-2xl font-bold">我的收藏</h1>{error && <p className="text-red-700">{error}</p>}{items.length ? items.map((item) => <TopicCard key={item.id} topic={item.topic} />) : <div className="paper-card rounded-xl p-10 text-center">暂未收藏主题。</div>}</div>;
}
