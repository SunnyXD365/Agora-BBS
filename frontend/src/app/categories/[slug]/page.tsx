'use client';

import { use, useEffect, useState } from 'react';
import Link from 'next/link';
import { Category, Topic } from '@/types/api';
import { categoriesApi, topicsApi } from '@/lib/api';

export default function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);
  const [category, setCategory] = useState<Category | null>(null);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchCategory = async () => {
      try {
        const catRes = await categoriesApi.getCategory(slug);
        setCategory(catRes.data);

        const topicRes = await topicsApi.getTopics({ category_id: catRes.data.id, page: 1, page_size: 20 });
        setTopics(topicRes.data.list);
      } catch (err: any) {
        setError(err.message || '加载失败');
      } finally {
        setLoading(false);
      }
    };

    fetchCategory();
  }, [slug]);

  if (loading) {
    return <div className="py-12 text-center text-gray-500">加载中...</div>;
  }

  if (error || !category) {
    return <div className="py-12 text-center text-red-500">{error || '板块不存在'}</div>;
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      <div className="bg-[#141414] p-6 border border-gray-800 rounded-lg">
        <h1 className="text-2xl font-bold text-gray-200">{category.name}</h1>
        {category.description && (
          <p className="mt-2 text-sm text-gray-500">{category.description}</p>
        )}
      </div>

      <div className="flex items-center justify-between">
        <h2 className="text-lg font-bold text-gray-200">主题列表</h2>
        <Link
          href={`/topics/new?category_id=${category.id}`}
          className="px-4 py-2 bg-gray-700 text-gray-100 rounded-md text-sm font-medium hover:bg-gray-600 transition"
        >
          发布新主题
        </Link>
      </div>

      <div className="space-y-3">
        {topics.length === 0 ? (
          <div className="p-6 bg-[#141414] border border-gray-800 rounded-lg text-center text-gray-500 text-sm">
            该板块暂无主题，快来发布第一个吧！
          </div>
        ) : (
          topics.map((topic) => (
            <Link
              key={topic.id}
              href={`/topics/${topic.id}`}
              className="block p-4 bg-[#141414] border border-gray-800 rounded-lg hover:border-gray-600 transition"
            >
              <div className="flex items-center justify-between">
                <h3 className="text-base font-semibold text-gray-200 truncate">{topic.title}</h3>
                <span className="text-xs text-gray-500 shrink-0 ml-4">
                  {new Date(topic.created_at).toLocaleDateString()}
                </span>
              </div>
              <div className="mt-2 flex items-center text-xs text-gray-500 space-x-4">
                <span>作者：{topic.author_name || `用户 #${topic.author_id}`}</span>
                <span>回复：{topic.reply_count}</span>
                <span>浏览：{topic.view_count}</span>
              </div>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
