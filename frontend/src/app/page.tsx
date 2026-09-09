'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Topic } from '@/types/api';
import { topicsApi } from '@/lib/api';
import CategoryNav from '@/components/CategoryNav';

export default function HomePage() {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTopics = async () => {
      try {
        const res = await topicsApi.getTopics({ page: 1, page_size: 10 });
        setTopics(res.data.list || []);
      } catch (err: any) {
        setError(err.message || '获取帖子列表失败');
      } finally {
        setLoading(false);
      }
    };

    fetchTopics();
  }, []);

  return (
    <div className="space-y-8">
      {/* 板块导航区域 */}
      <section className="bg-[#141414] p-6 border border-gray-800 rounded-lg">
        <h2 className="text-lg font-bold text-gray-200 mb-4">板块导航</h2>
        <CategoryNav />
      </section>

      {/* 最新主题区域 */}
      <section>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-bold text-gray-200">最新讨论</h2>
          <Link href="/topics" className="text-sm text-gray-400 hover:text-gray-200 transition">
            查看全部 →
          </Link>
        </div>

        {loading ? (
          <div className="py-12 text-center text-gray-500">帖子加载中...</div>
        ) : error ? (
          <div className="py-12 text-center text-red-500">{error}</div>
        ) : topics.length === 0 ? (
          <div className="p-8 text-center bg-[#141414] border border-gray-800 rounded-lg text-gray-500">
            暂无主题帖，赶快去发布第一个帖子吧！
          </div>
        ) : (
          <div className="space-y-3">
            {topics.map((topic) => (
              <div
                key={topic.id}
                className="p-4 bg-[#141414] border border-gray-800 rounded-lg hover:border-gray-600 transition flex justify-between items-center"
              >
                <div className="space-y-1">
                  <Link
                    href={`/topics/${topic.id}`}
                    className="text-lg font-semibold text-gray-200 hover:text-gray-100 transition"
                  >
                    {topic.title}
                  </Link>
                  <div className="flex items-center space-x-3 text-xs text-gray-500">
                    <span>作者：{topic.author_name || `用户 #${topic.author_id}`}</span>
                    <span>•</span>
                    <span>{new Date(topic.created_at).toLocaleString()}</span>
                  </div>
                </div>

                <div className="flex items-center space-x-4 text-xs text-gray-500 shrink-0">
                  <span className="bg-gray-800 px-2 py-1 rounded">浏览 {topic.view_count}</span>
                  <span className="bg-gray-800 text-gray-300 px-2 py-1 rounded">回复 {topic.reply_count}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
