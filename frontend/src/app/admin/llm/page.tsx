'use client';

import { useCallback, useEffect, useState } from 'react';
import Pagination from '@/components/Pagination';
import { AdminListToolbar, Notice, PageTitle, Panel, Td, Th } from '@/components/admin/AdminUI';
import { adminApi, getErrorMessage } from '@/services';
import { AdminLLMJob } from '@/types/api';
import { useDebouncedValue } from '@/hooks/useDebouncedValue';

export default function LLMPage() {
  const [items, setItems] = useState<AdminLLMJob[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState('');
  const [queryInput, setQueryInput] = useState('');
  const query = useDebouncedValue(queryInput);
  const [sort, setSort] = useState('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const load = useCallback(async () => {
    setBusy(true); setError('');
    try {
      const result = await adminApi.llmJobs({ q: query || undefined, page, page_size: pageSize, status: status || undefined, sort, order });
      setItems(result.data.items); setTotal(result.data.total);
    } catch (loadError) { setError(getErrorMessage(loadError, 'LLM 作业加载失败')); }
    finally { setBusy(false); }
  }, [order, page, pageSize, query, sort, status]);
  useEffect(() => { void Promise.resolve().then(load); }, [load]);

  return <div className="mx-auto max-w-7xl space-y-6">
    <PageTitle title="LLM 作业" subtitle="按类型、对象、模型、状态、延迟和 Token 检索调用记录" />
    {error && <Notice>{error}</Notice>}
    <Panel>
      <AdminListToolbar query={queryInput} onQueryChange={(value) => { setQueryInput(value); setPage(1); }} searchPlaceholder="搜索类型、对象 ID、模型或错误" sort={sort} onSortChange={(value) => { setSort(value); setPage(1); }} order={order} onOrderChange={(value) => { setOrder(value); setPage(1); }} sortOptions={[{ value: 'created_at', label: '按创建时间' }, { value: 'latency_ms', label: '按延迟' }, { value: 'tokens', label: '按 Token' }, { value: 'attempts', label: '按尝试次数' }, { value: 'status', label: '按状态' }]}>
        <select aria-label="作业状态筛选" value={status} onChange={(event) => { setStatus(event.target.value); setPage(1); }} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm"><option value="">全部状态</option>{['pending', 'running', 'processing', 'completed', 'failed'].map((value) => <option key={value}>{value}</option>)}</select>
      </AdminListToolbar>
      <div className="overflow-x-auto"><table className="w-full"><thead><tr><Th>类型 / 对象</Th><Th>状态</Th><Th>模型</Th><Th>尝试</Th><Th>延迟</Th><Th>Token</Th><Th>错误 / 操作</Th></tr></thead><tbody>
        {items.map((item) => <tr key={item.id} className="border-t border-slate-100"><Td>{item.job_type}<div className="text-xs text-slate-400">{item.aggregate_type} #{item.aggregate_id}</div></Td><Td>{item.status}</Td><Td>{item.model || '—'}</Td><Td>{item.attempts}</Td><Td>{item.latency_ms} ms</Td><Td>{item.prompt_tokens + item.completion_tokens}</Td><Td><span className="text-xs text-red-700">{item.error_message}</span>{item.status === 'failed' && <button disabled={busy} onClick={async () => { setBusy(true); try { await adminApi.retryLLMJob(item.id); await load(); } catch (actionError) { setError(getErrorMessage(actionError, '重试失败')); } finally { setBusy(false); } }} className="ml-2 font-medium text-cyan-700 hover:underline">安全重试</button>}</Td></tr>)}
        {!busy && items.length === 0 && <tr><td colSpan={7} className="border-t px-3 py-10 text-center text-sm text-slate-500">没有符合筛选条件的 LLM 作业</td></tr>}
      </tbody></table></div>
      <div className="mt-4"><Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPage(1); setPageSize(size); }} itemLabel="个作业" disabled={busy} /></div>
    </Panel>
  </div>;
}
