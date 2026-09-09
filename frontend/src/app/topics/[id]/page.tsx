'use client';

import { use, useEffect, useState } from 'react';
import { Topic, Post } from '@/types/api';
import { topicsApi, postsApi } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import Link from 'next/link';

export default function TopicDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const { user } = useAuth();

  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // 发回复表单状态
  const [replyContent, setReplyContent] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [replyError, setReplyError] = useState('');

  useEffect(() => {
    const fetchTopicDetail = async () => {
      try {
        const [topicRes, postsRes] = await Promise.all([
          topicsApi.getTopicDetail(id),
          postsApi.getPosts(id),
        ]);
        setTopic(topicRes.data);
        setPosts(postsRes.data);
      } catch (err: any) {
        setError(err.message || '获取帖子详情失败');
      } finally {
        setLoading(false);
      }
    };

    fetchTopicDetail();
  }, [id]);

  // 提交回复
  const handlePostSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    setReplyError('');
    setSubmitting(true);

    try {
      const res = await postsApi.createPost(id, {
        parent_id: null,
        content: replyContent,
      });

      // 重新获取回复列表，确保数据一致
      const postsRes = await postsApi.getPosts(id);
      setPosts(postsRes.data);
      setReplyContent('');
      if (topic) {
        setTopic({ ...topic, reply_count: topic.reply_count + 1 });
      }
    } catch (err: any) {
      setReplyError(err.message || '回复失败，请重试');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="py-12 text-center text-gray-500">帖子加载中...</div>;
  }

  if (error || !topic) {
    return <div className="py-12 text-center text-red-500">{error || '帖子不存在'}</div>;
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto">
      {/* 楼主贴内容卡片 */}
      <div className="bg-[#141414] p-6 border border-gray-800 rounded-lg space-y-4">
        <h1 className="text-2xl font-bold text-gray-200">{topic.title}</h1>
        
        <div className="flex items-center space-x-4 text-xs text-gray-500 border-b border-gray-800 pb-4">
          <span className="font-semibold text-gray-300">
            {topic.author_name || `用户 #${topic.author_id}`}
          </span>
          <span>•</span>
          <span>发布于 {new Date(topic.created_at).toLocaleString()}</span>
          <span>•</span>
          <span>浏览量: {topic.view_count}</span>
        </div>

        <div className="text-gray-300 whitespace-pre-wrap min-h-[100px] leading-relaxed">
          {topic.content || '（暂无正文）'}
        </div>
      </div>

      {/* 回复楼层列表 */}
      <div className="space-y-4">
        <h2 className="text-lg font-bold text-gray-200">全部回复 ({topic.reply_count})</h2>

        {posts.length === 0 ? (
          <div className="p-6 bg-[#141414] border border-gray-800 rounded-lg text-center text-gray-500 text-sm">
            暂无楼层回复，抢沙发吧！
          </div>
        ) : (
          posts.map((post, index) => (
            <div key={post.id || index} className="p-4 bg-[#141414] border border-gray-800 rounded-lg space-y-2">
              <div className="flex justify-between items-center text-xs text-gray-500 border-b border-gray-800 pb-2">
                <span>{post.author_name || `用户 #${post.author_id}`}</span>
                <span>#{index + 1} 楼 • {new Date(post.created_at).toLocaleString()}</span>
              </div>
              <p className="text-gray-300 text-sm whitespace-pre-wrap">{post.content}</p>
            </div>
          ))
        )}
      </div>

      {/* 发表回复框 */}
      <div className="bg-[#141414] p-6 border border-gray-800 rounded-lg space-y-4">
        <h3 className="text-md font-bold text-gray-200">发表回复</h3>

        {user ? (
          <form onSubmit={handlePostSubmit} className="space-y-3">
            {replyError && (
              <div className="p-2 bg-red-900/30 text-red-400 text-xs rounded border border-red-800">
                {replyError}
              </div>
            )}
            <textarea
              required
              rows={4}
              value={replyContent}
              onChange={(e) => setReplyContent(e.target.value)}
              className="w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-500 text-sm text-gray-200 placeholder-gray-500"
              placeholder="撰写你的回复..."
            />
            <div className="flex justify-end">
              <button
                type="submit"
                disabled={submitting || !replyContent.trim()}
                className="px-4 py-2 bg-gray-700 text-gray-100 rounded-md text-sm font-medium hover:bg-gray-600 transition disabled:bg-gray-800 disabled:text-gray-500"
              >
                {submitting ? '提交中...' : '提交回复'}
              </button>
            </div>
          </form>
        ) : (
          <div className="p-4 text-center bg-gray-900 border border-dashed border-gray-700 rounded-md text-sm text-gray-400">
            你需要{' '}
            <Link href="/login" className="text-gray-300 font-medium hover:text-gray-100 underline">
              登录
            </Link>{' '}
            后才能参与讨论。
          </div>
        )}
      </div>
    </div>
  );
}
