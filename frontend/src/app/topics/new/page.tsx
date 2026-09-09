'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { topicsApi, categoriesApi } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { Category } from '@/types/api';

export default function NewTopicPage() {
  const router = useRouter();
  const { user } = useAuth();
  
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [categoryId, setCategoryId] = useState<number | ''>('');
  const [categories, setCategories] = useState<Category[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  // 获取板块列表
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await categoriesApi.getCategories();
        setCategories(res.data);
        // 如果有板块，默认选中第一个
        if (res.data.length > 0) {
          setCategoryId(res.data[0].id);
        }
      } catch (err: any) {
        console.error('获取板块列表失败:', err);
      }
    };

    fetchCategories();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) {
      setError('请先登录后再发布帖子');
      return;
    }

    if (!categoryId) {
      setError('请选择板块');
      return;
    }

    setError('');
    setSubmitting(true);

    try {
      const res = await topicsApi.createTopic({
        category_id: Number(categoryId),
        title,
        content,
      });
      // 发布成功，跳转至新帖子详情页
      router.push(`/topics/${res.data.id}`);
    } catch (err: any) {
      setError(err.message || '发布失败，请稍后重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto bg-[#141414] p-6 border border-gray-800 rounded-lg">
      <h1 className="text-2xl font-bold mb-6 text-gray-200">发布新主题帖</h1>

      {error && (
        <div className="mb-4 p-3 bg-red-900/30 text-red-400 text-sm rounded-md border border-red-800">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">选择板块</label>
          <select
            value={categoryId}
            onChange={(e) => setCategoryId(e.target.value ? Number(e.target.value) : '')}
            className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-gray-200"
          >
            <option value="">请选择板块</option>
            {categories.map((cat) => (
              <option key={cat.id} value={cat.id}>
                {cat.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">标题</label>
          <input
            type="text"
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-gray-200 placeholder-gray-500"
            placeholder="概括主题的主要内容..."
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">正文内容</label>
          <textarea
            required
            rows={8}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-gray-200 placeholder-gray-500 resize-y"
            placeholder="请输入正文内容..."
          />
        </div>

        <div className="flex justify-end space-x-3">
          <button
            type="button"
            onClick={() => router.back()}
            className="px-4 py-2 border border-gray-700 text-gray-400 rounded-md hover:bg-gray-800 transition"
          >
            取消
          </button>
          <button
            type="submit"
            disabled={submitting}
            className="px-4 py-2 bg-gray-700 text-gray-100 rounded-md hover:bg-gray-600 font-medium transition disabled:bg-gray-800 disabled:text-gray-500"
          >
            {submitting ? '发布中...' : '发布帖子'}
          </button>
        </div>
      </form>
    </div>
  );
}
