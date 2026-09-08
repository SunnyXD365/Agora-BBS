'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Topic } from '@/types/api';
import { topicsApi } from '@/lib/api';

export default function HomePage() {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchTopics = async () => {
      try {
        const res = await topicsApi.getTopics({ page: 1, page_size: 20 });
        setTopics(res.data.list || []);
      } catch (err: any) {
        setError(err.message || '获取帖子列表失败');
      } finally {
        setLoading(false);
      }
    };

    fetchTopics();
  }, []);

  if (loading) {
    return <div className="py-12 text-center text-gray-500">帖子加载中...</div>;
  }

  if (error) {
    return <div className="py-12 text-center text-red-500">{error}</div>;
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">最新讨论</h1>
      </div>

      {topics.length === 0 ? (
        <div className="p-8 text-center bg-white border rounded-lg text-gray-500">
          暂无主题帖，赶快去发布第一个帖子吧！
        </div>
      ) : (
        <div className="space-y-3">
          {topics.map((topic) => (
            <div
              key={topic.id}
              className="p-4 bg-white border rounded-lg hover:shadow-md transition flex justify-between items-center"
            >
              <div className="space-y-1">
                <Link
                  href={`/topics/${topic.id}`}
                  className="text-lg font-semibold text-gray-900 hover:text-blue-600 transition"
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
                <span className="bg-gray-100 px-2 py-1 rounded">浏览 {topic.view_count}</span>
                <span className="bg-blue-50 text-blue-600 px-2 py-1 rounded">回复 {topic.reply_count}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
