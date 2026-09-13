'use client';

import { useCallback, useEffect, useState } from 'react';
import Pagination from '@/components/Pagination';
import { AdminListToolbar, Notice, PageTitle, Panel, Td, Th } from '@/components/admin/AdminUI';
import { adminApi, getErrorMessage } from '@/services';
import { AdminTrustLog } from '@/types/api';
import { useDebouncedValue } from '@/hooks/useDebouncedValue';

export default function TrustPage() {
  const [items, setItems] = useState<AdminTrustLog[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [total, setTotal] = useState(0);
  const [userId, setUserId] = useState('');
  const [queryInput, setQueryInput] = useState('');
  const query = useDebouncedValue(queryInput);
  const [sort, setSort] = useState('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const load = useCallback(async () => {
    setBusy(true); setError('');
    try {
      const result = await adminApi.trustLogs({ q: query || undefined, page, page_size: pageSize, user_id: userId ? Number(userId) : undefined, sort, order });
      setItems(result.data.items); setTotal(result.data.total);
    } catch (loadError) { setError(getErrorMessage(loadError, '信任流水加载失败')); }
    finally { setBusy(false); }
  }, [order, page, pageSize, query, sort, userId]);
  useEffect(() => { void Promise.resolve().then(load); }, [load]);

  return <div className="mx-auto max-w-7xl space-y-6">
    <PageTitle title="信任流水" subtitle="按用户、事件、分值与时间追溯不可覆盖的信任变化" />
    {error && <Notice>{error}</Notice>}
    <Panel>
      <AdminListToolbar query={queryInput} onQueryChange={(value) => { setQueryInput(value); setPage(1); }} searchPlaceholder="搜索用户名、事件、原因或关联类型" sort={sort} onSortChange={(value) => { setSort(value); setPage(1); }} order={order} onOrderChange={(value) => { setOrder(value); setPage(1); }} sortOptions={[{ value: 'created_at', label: '按发生时间' }, { value: 'score_delta', label: '按分值变化' }, { value: 'username', label: '按用户名' }, { value: 'event_type', label: '按事件类型' }]}>
        <input aria-label="用户 ID 筛选" type="number" min="1" value={userId} onChange={(event) => { setUserId(event.target.value); setPage(1); }} placeholder="用户 ID" className="w-28 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm" />
      </AdminListToolbar>
      <div className="overflow-x-auto"><table className="w-full"><thead><tr><Th>用户</Th><Th>事件</Th><Th>变化</Th><Th>原因</Th><Th>关联对象</Th><Th>时间</Th></tr></thead><tbody>
        {items.map((item) => <tr key={item.id} className="border-t border-slate-100"><Td>{item.username}<div className="text-xs text-slate-400">#{item.user_id}</div></Td><Td>{item.event_type}</Td><Td><strong className={item.score_delta >= 0 ? 'text-emerald-700' : 'text-red-700'}>{item.score_delta >= 0 ? '+' : ''}{item.score_delta}</strong></Td><Td>{item.reason}</Td><Td>{item.reference_type} {item.reference_id || ''}</Td><Td>{new Date(item.created_at).toLocaleString()}</Td></tr>)}
        {!busy && items.length === 0 && <tr><td colSpan={6} className="border-t px-3 py-10 text-center text-sm text-slate-500">没有符合筛选条件的信任流水</td></tr>}
      </tbody></table></div>
      <div className="mt-4"><Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPage(1); setPageSize(size); }} itemLabel="条流水" disabled={busy} /></div>
    </Panel>
  </div>;
}
