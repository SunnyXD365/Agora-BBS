'use client';

import { useCallback, useEffect, useState } from 'react';
import Pagination from '@/components/Pagination';
import { AdminListToolbar, Notice, PageTitle, Panel } from '@/components/admin/AdminUI';
import { adminApi, getErrorMessage } from '@/services';
import { AdminContent } from '@/types/api';
import { useDebouncedValue } from '@/hooks/useDebouncedValue';

export default function ContentPage() {
  const [type, setType] = useState<'topic' | 'post'>('topic');
  const [status, setStatus] = useState('');
  const [items, setItems] = useState<AdminContent[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [queryInput, setQueryInput] = useState('');
  const query = useDebouncedValue(queryInput);
  const [sort, setSort] = useState('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setBusy(true); setError('');
    try {
      const result = await adminApi.contents(type, { q: query || undefined, page, page_size: pageSize, status: status || undefined, sort, order });
      setItems(result.data.items); setTotal(result.data.total);
    } catch (loadError) { setError(getErrorMessage(loadError, '内容加载失败')); }
    finally { setBusy(false); }
  }, [order, page, pageSize, query, sort, status, type]);
  useEffect(() => { void Promise.resolve().then(load); }, [load]);

  return <div className="mx-auto max-w-7xl space-y-6">
    <PageTitle title="内容治理" subtitle="检索主题与回复；只能隐藏或恢复，不能篡改正文" />
    {error && <Notice>{error}</Notice>}
    <Panel>
      <AdminListToolbar query={queryInput} onQueryChange={(value) => { setQueryInput(value); setPage(1); }} searchPlaceholder="搜索标题、正文或作者" sort={sort} onSortChange={(value) => { setSort(value); setPage(1); }} order={order} onOrderChange={(value) => { setOrder(value); setPage(1); }} sortOptions={[{ value: 'created_at', label: '按创建时间' }, { value: 'title', label: type === 'topic' ? '按标题' : '按回复正文' }, { value: 'author', label: '按作者' }, { value: 'status', label: '按状态' }]}>
        <select aria-label="内容类型筛选" value={type} onChange={(event) => { setType(event.target.value as 'topic' | 'post'); setPage(1); }} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm"><option value="topic">主题</option><option value="post">回复</option></select>
        <select aria-label="内容状态筛选" value={status} onChange={(event) => { setStatus(event.target.value); setPage(1); }} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm"><option value="">全部状态</option>{['cooling', 'pending_review', 'published', 'rejected', 'recalled', 'hidden'].map((value) => <option key={value}>{value}</option>)}</select>
      </AdminListToolbar>
      <div className="space-y-3">
        {items.map((item) => <article key={`${item.type}-${item.id}`} className="flex items-start justify-between gap-4 rounded-xl border border-slate-200 p-4"><div><strong>{item.title}</strong><p className="mt-1 text-xs text-slate-500">{item.author_name} · {item.status} · {new Date(item.created_at).toLocaleString()}</p><p className="mt-2 line-clamp-2 text-sm text-slate-600">{item.excerpt}</p></div>{['published', 'hidden'].includes(item.status) && <button disabled={busy} onClick={async () => { setBusy(true); try { await adminApi.setContentVisibility(item.type, item.id, item.status !== 'hidden'); await load(); } catch (actionError) { setError(getErrorMessage(actionError, '操作失败')); } finally { setBusy(false); } }} className="shrink-0 rounded-lg bg-slate-900 px-3 py-2 text-xs text-white">{item.status === 'hidden' ? '恢复' : '隐藏'}</button>}</article>)}
        {!busy && items.length === 0 && <div className="rounded-xl border border-dashed p-10 text-center text-sm text-slate-500">没有符合筛选条件的内容</div>}
      </div>
      <div className="mt-4"><Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPage(1); setPageSize(size); }} itemLabel={type === 'topic' ? '个主题' : '条回复'} disabled={busy} /></div>
    </Panel>
  </div>;
}
