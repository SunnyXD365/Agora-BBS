'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Category, GovernancePolicy, Topic } from '@/types/api';
import { getErrorMessage, governanceApi, topicApi } from '@/services';
import { useAuth } from '@/context/AuthContext';
import TopicCard from '@/components/TopicCard';
import UsageGuideCard from '@/components/UsageGuideCard';
import CommunitySidebar from '@/components/CommunitySidebar';
import Pagination from '@/components/Pagination';

export default function HomePage() {
  const { user } = useAuth();
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<number | undefined>(undefined);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [topicTotal, setTopicTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [policy, setPolicy] = useState<GovernancePolicy | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // 初始化加载板块列表
  useEffect(() => {
    let cancelled = false;
    Promise.all([topicApi.getCategories(), governanceApi.policy()])
      .then(([categoryRes, policyRes]) => {
        if (cancelled) return;
        if (categoryRes.code === 0) setCategories(categoryRes.data);
        if (policyRes.code === 0) setPolicy(policyRes.data);
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
      .getTopics({ category_id: selectedCategory, page, page_size: pageSize })
      .then((res) => {
        if (!cancelled && res.code === 0) { setTopics(res.data.items); setTopicTotal(res.data.total); }
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(getErrorMessage(err, '主题加载失败'));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [page, pageSize, selectedCategory]);

  const changePage = (nextPage: number) => {
    setLoading(true); setPage(nextPage);
    document.getElementById('topic-list')?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-[210px_minmax(0,1fr)] xl:grid-cols-[210px_minmax(0,1fr)_220px]">
      {/* 左侧：新用户导航 */}
      <div className="order-2 lg:order-1">
        <UsageGuideCard />
      </div>

      {/* 主栏：板块过滤与帖子列表 */}
      <div id="topic-list" className="order-1 min-w-0 scroll-mt-6 space-y-4 lg:order-2">
        {/* 板块 Selector 标栏 */}
        <div className="flex items-center justify-between rounded-lg border bg-white p-3 shadow-sm">
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => { setLoading(true); setPage(1); setSelectedCategory(undefined); }}
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
                onClick={() => { setLoading(true); setPage(1); setSelectedCategory(c.id); }}
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

          <Link
            href={user ? '/topics/new' : '/login'}
            className="rounded-md bg-blue-600 px-3.5 py-1.5 text-xs font-medium text-white hover:bg-blue-500 transition-colors"
          >
            + 发新帖
          </Link>
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
        {!loading && <Pagination page={page} pageSize={pageSize} total={topicTotal} onPageChange={changePage} onPageSizeChange={(size) => { setLoading(true); setPage(1); setPageSize(size); }} itemLabel="个主题" />}
      </div>

      {/* 右侧：成长、热议与社区规则 */}
      <div className="order-3 hidden xl:block">
        <CommunitySidebar user={user} topics={topics} topicTotal={topicTotal} categoryCount={categories.length} policy={policy} />
      </div>
    </div>
  );
}
