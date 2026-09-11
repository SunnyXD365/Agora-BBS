'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { Post, Topic } from '@/types/api';
import { likeApi, postApi, topicApi } from '@/services';
import { useAuth } from '@/context/AuthContext';

export default function TopicDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { user } = useAuth();
  
  const topicId = Number(params.id);

  const [topic, setTopic] = useState<Topic | null>(null);
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  // 互动与回复状态
  const [isLiked, setIsLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [replyContent, setReplyContent] = useState('');
  const [submitting, setSubmitting] = useState(false);

  // 加载帖子详情与楼层列表
  const fetchData = async () => {
    if (!topicId) return;
    setLoading(true);

    try {
      const [topicRes, postsRes] = await Promise.all([
        topicApi.getTopicDetail(topicId),
        postApi.getPosts(topicId, { page: 1, page_size: 50 }),
      ]);

      if (topicRes.code === 0 && topicRes.data) {
        setTopic(topicRes.data);
        setLikeCount(topicRes.data.like_count);
      }
      if (postsRes.code === 0) {
        setPosts(postsRes.data || []);
      }
    } catch (err) {
      console.error('Failed to load topic detail:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [topicId]);

  // 点赞/取消点赞主题帖
  const handleToggleTopicLike = async () => {
    if (!user) {
      alert('请先登录后再点赞');
      router.push('/login');
      return;
    }

    const nextState = !isLiked;
    setIsLiked(nextState);
    setLikeCount((prev) => (nextState ? prev + 1 : prev - 1));

    try {
      await likeApi.toggleLike('topic', topicId, nextState);
    } catch (err) {
      // 回滚状态
      setIsLiked(!nextState);
      setLikeCount((prev) => (nextState ? prev - 1 : prev + 1));
    }
  };

  // 发表回复
  const handleSubmitReply = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) {
      alert('请先登录后再回复');
      router.push('/login');
      return;
    }

    if (!replyContent.trim()) return;

    setSubmitting(true);
    try {
      const res = await postApi.createPost(topicId, { content: replyContent });
      if (res.code === 0) {
        setReplyContent('');
        // 重新拉取回复列表并更新帖子回复计数
        const postsRes = await postApi.getPosts(topicId, { page: 1, page_size: 50 });
        if (postsRes.code === 0) setPosts(postsRes.data || []);
        if (topic) setTopic({ ...topic, post_count: topic.post_count + 1 });
      }
    } catch (err: any) {
      alert(err.msg || '回复发表失败');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="mx-auto max-w-4xl space-y-4">
        <div className="h-48 animate-pulse rounded-xl bg-gray-200" />
        <div className="h-24 animate-pulse rounded-xl bg-gray-200" />
      </div>
    );
  }

  if (!topic) {
    return (
      <div className="rounded-xl border bg-white p-12 text-center text-gray-500">
        该主题帖不存在或已被删除。
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      {/* 帖子主卡片 */}
      <article className="rounded-xl border bg-white p-6 shadow-sm">
        <header className="border-b pb-4">
          <h1 className="text-2xl font-bold text-gray-900">{topic.title}</h1>
          <div className="mt-3 flex items-center justify-between text-xs text-gray-500">
            <div className="flex items-center gap-3">
              <span className="font-semibold text-gray-700">{topic.author_name}</span>
              <span>·</span>
              <span>{new Date(topic.created_at).toLocaleString()}</span>
            </div>
            <div className="flex items-center gap-4">
              <span>👀 {topic.view_count} 浏览</span>
              <span>💬 {topic.post_count} 回复</span>
            </div>
          </div>
        </header>

        {/* 正文内容 */}
        <div className="py-6 text-sm text-gray-800 leading-relaxed whitespace-pre-wrap">
          {topic.content}
        </div>

        {/* 交互操作区 */}
        <footer className="flex items-center justify-end border-t pt-4">
          <button
            onClick={handleToggleTopicLike}
            className={`flex items-center gap-1.5 rounded-full px-4 py-1.5 text-xs font-medium transition-colors ${
              isLiked
                ? 'bg-red-50 text-red-600 border border-red-200'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            {isLiked ? '❤️ 已赞' : '🤍 点赞'} ({likeCount})
          </button>
        </footer>
      </article>

      {/* 发表回复框 */}
      <section className="rounded-xl border bg-white p-6 shadow-sm">
        <h3 className="text-sm font-bold text-gray-900">发表回复</h3>
        <form onSubmit={handleSubmitReply} className="mt-3 space-y-3">
          <textarea
            rows={3}
            required
            value={replyContent}
            onChange={(e) => setReplyContent(e.target.value)}
            placeholder={user ? '撰写你的观点...' : '请先登录后再发表回复'}
            disabled={!user}
            className="w-full rounded-md border border-gray-300 p-3 text-sm focus:border-blue-500 focus:outline-none disabled:bg-gray-50"
          />
          <div className="flex justify-end">
            <button
              type="submit"
              disabled={submitting || !user}
              className="rounded-md bg-blue-600 px-5 py-2 text-xs font-medium text-white hover:bg-blue-500 disabled:opacity-50 transition-colors"
            >
              {submitting ? '提交中...' : '提交回复'}
            </button>
          </div>
        </form>
      </section>

      {/* 楼层回复列表 */}
      <section className="space-y-3">
        <h3 className="text-sm font-bold text-gray-900 px-1">
          全部回复 ({posts.length})
        </h3>
        {posts.length === 0 ? (
          <div className="rounded-xl border bg-white p-8 text-center text-xs text-gray-500">
            暂无回复，快来发表第一条观点吧！
          </div>
        ) : (
          posts.map((post, idx) => (
            <div key={post.id} className="rounded-xl border bg-white p-4 shadow-sm space-y-2">
              <div className="flex items-center justify-between text-xs text-gray-500">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-gray-700">{post.author_name}</span>
                  <span className="rounded bg-gray-100 px-1.5 py-0.5 text-[10px] text-gray-500">
                    #{idx + 1} 楼
                  </span>
                </div>
                <span>{new Date(post.created_at).toLocaleString()}</span>
              </div>
              <p className="text-sm text-gray-800 leading-normal">{post.content}</p>
            </div>
          ))
        )}
      </section>
    </div>
  );
}