'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Category, Topic } from '@/types/api';
import { categoriesApi, topicsApi } from '@/lib/api';

export default function TopicsPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>(undefined);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const pageSize = 10;

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await categoriesApi.getCategories();
        setCategories(res.data);
      } catch (err: any) {
        console.error('获取板块失败', err);
      }
    };
    fetchCategories();
  }, []);

  useEffect(() => {
    const fetchTopics = async () => {
      setLoading(true);
      setError('');
      try {
        const res = await topicsApi.getTopics({
          page,
          page_size: pageSize,
          category_id: selectedCategory,
        });
        setTopics(res.data.list);
        setTotal(res.data.pagination.total);
      } catch (err: any) {
        setError(err.message || '加载失败');
      } finally {
        setLoading(false);
      }
    };

    fetchTopics();
  }, [page, selectedCategory]);

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">全部主题</h1>
        <Link
          href="/topics/new"
          className="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-medium hover:bg-blue-700 transition"
        >
          发布新主题
        </Link>
      </div>

      {/* 板块筛选 */}
      <div className="flex flex-wrap gap-2">
        <button
          onClick={() => {
            setSelectedCategory(undefined);
            setPage(1);
          }}
          className={`px-3 py-1 rounded-full text-sm ${
            selectedCategory === undefined
              ? 'bg-blue-600 text-white'
              : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
          }`}
        >
          全部
        </button>
        {categories.map((cat) => (
          <button
            key={cat.id}
            onClick={() => {
              setSelectedCategory(cat.id);
              setPage(1);
            }}
            className={`px-3 py-1 rounded-full text-sm ${
              selectedCategory === cat.id
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            {cat.name}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="py-12 text-center text-gray-500">加载中...</div>
      ) : error ? (
        <div className="py-12 text-center text-red-500">{error}</div>
      ) : topics.length === 0 ? (
        <div className="p-6 bg-white border rounded-lg text-center text-gray-400 text-sm">
          暂无主题，快来发布第一个吧！
        </div>
      ) : (
        <div className="space-y-3">
          {topics.map((topic) => (
            <Link
              key={topic.id}
              href={`/topics/${topic.id}`}
              className="block p-4 bg-white border rounded-lg shadow-sm hover:shadow-md transition"
            >
              <div className="flex items-center justify-between">
                <h3 className="text-base font-semibold text-gray-900 truncate">{topic.title}</h3>
                <span className="text-xs text-gray-400 shrink-0 ml-4">
                  {new Date(topic.created_at).toLocaleDateString()}
                </span>
              </div>
              <div className="mt-2 flex items-center text-xs text-gray-500 space-x-4">
                <span>作者：{topic.author_name || `用户 #${topic.author_id}`}</span>
                <span>回复：{topic.reply_count}</span>
                <span>浏览：{topic.view_count}</span>
              </div>
            </Link>
          ))}
        </div>
      )}

      {/* 分页 */}
      {totalPages > 1 && (
        <div className="flex justify-center space-x-2">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-3 py-1 bg-gray-200 rounded-md text-sm disabled:opacity-50"
          >
            上一页
          </button>
          <span className="px-3 py-1 text-sm text-gray-700">
            {page} / {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page === totalPages}
            className="px-3 py-1 bg-gray-200 rounded-md text-sm disabled:opacity-50"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  );
}
