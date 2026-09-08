'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { topicsApi } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';

export default function NewTopicPage() {
  const router = useRouter();
  const { user } = useAuth();
  
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [categoryId] = useState(1); // MVP 阶段写死分类 1
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) {
      setError('请先登录后再发布帖子');
      return;
    }

    setError('');
    setSubmitting(true);

    try {
      const res = await topicsApi.createTopic({
        category_id: categoryId,
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
    <div className="max-w-3xl mx-auto bg-white p-6 border rounded-lg shadow-sm">
      <h1 className="text-2xl font-bold mb-6 text-gray-900">发布新主题帖</h1>

      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-600 text-sm rounded-md border border-red-200">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">标题</label>
          <input
            type="text"
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="概括主题的主要内容..."
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">正文内容</label>
          <textarea
            required
            rows={8}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 resize-y"
            placeholder="请输入正文内容..."
          />
        </div>

        <div className="flex justify-end space-x-3">
          <button
            type="button"
            onClick={() => router.back()}
            className="px-4 py-2 border text-gray-600 rounded-md hover:bg-gray-50 transition"
          >
            取消
          </button>
          <button
            type="submit"
            disabled={submitting}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium transition disabled:bg-gray-400"
          >
            {submitting ? '发布中...' : '发布帖子'}
          </button>
        </div>
      </form>
    </div>
  );
}
