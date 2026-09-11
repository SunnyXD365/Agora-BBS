'use client';

import { useEffect, useState } from 'react';
import { Category, Topic } from '@/types/api';
import { getErrorMessage, topicApi } from '@/services';
import { useAuth } from '@/context/AuthContext';
import TopicCard from '@/components/TopicCard';
import CreateTopicModal from '@/components/CreateTopicModal';

export default function HomePage() {
  const { user } = useAuth();
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>(undefined);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [error, setError] = useState('');
  const [refreshKey, setRefreshKey] = useState(0);

  // 初始化加载板块列表
  useEffect(() => {
    let cancelled = false;
    topicApi.getCategories()
      .then((res) => {
        if (!cancelled && res.code === 0) setCategories(res.data);
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(getErrorMessage(err, '板块加载失败'));
      });
    return () => { cancelled = true; };
  }, []);

  // 当选中的板块变化或发帖成功时重新拉取帖子列表
  useEffect(() => {
    let cancelled = false;
    topicApi
      .getTopics({ category_id: selectedCategory, page: 1, page_size: 20 })
      .then((res) => {
        if (!cancelled && res.code === 0) setTopics(res.data.items);
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(getErrorMessage(err, '主题加载失败'));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [selectedCategory, refreshKey]);

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-4">
      {/* 主栏：板块过滤与帖子列表 */}
      <div className="lg:col-span-3 space-y-4">
        {/* 板块 Selector 标栏 */}
        <div className="flex items-center justify-between rounded-lg border bg-white p-3 shadow-sm">
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setSelectedCategory(undefined)}
              className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
                selectedCategory === undefined
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              全部
            </button>
            {categories.map((c) => (
              <button
                key={c.id}
                onClick={() => { setLoading(true); setSelectedCategory(c.id); }}
                className={`rounded-md px-3 py-1.5 text-xs font-medium transition-colors ${
                  selectedCategory === c.id
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {c.name}
              </button>
            ))}
          </div>

          <button
            onClick={() => {
              if (!user) {
                alert('请先登录后再发帖');
                return;
              }
              if (!user.capabilities.includes('create_topic')) {
                alert('发起主题尚未解锁，请先在成长中心查看阅读进度。');
                return;
              }
              setIsModalOpen(true);
            }}
            className="rounded-md bg-blue-600 px-3.5 py-1.5 text-xs font-medium text-white hover:bg-blue-500 transition-colors"
          >
            + 发新帖
          </button>
        </div>

        {/* 帖子列表渲染 */}
        {error && <div className="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</div>}
        {loading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-24 animate-pulse rounded-lg bg-gray-200" />
            ))}
          </div>
        ) : topics.length === 0 ? (
          <div className="rounded-lg border bg-white p-12 text-center text-sm text-gray-500">
            该板块下暂无帖子，快来抢沙发吧！
          </div>
        ) : (
          <div className="space-y-3">
            {topics.map((topic) => (
              <TopicCard key={topic.id} topic={topic} />
            ))}
          </div>
        )}
      </div>

      {/* 侧边栏：社区信息 */}
      <div className="space-y-4">
        <div className="rounded-lg border bg-white p-4 shadow-sm">
          <h3 className="font-semibold text-gray-900 text-sm">关于 Agora BBS</h3>
          <p className="mt-2 text-xs text-gray-500 leading-relaxed">
            基于 Go (Gin) + Next.js 构建的高性能论坛社区，支持高并发与轻量治理模式。
          </p>
        </div>
      </div>

      {/* 发帖弹窗 */}
      <CreateTopicModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        categories={categories}
        onSuccess={() => { setLoading(true); setRefreshKey((key) => key + 1); }}
      />
    </div>
  );
}
