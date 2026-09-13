'use client';

import { useCallback, useEffect, useState } from 'react';
import { AdminListToolbar, Notice, PageTitle, Panel } from '@/components/admin/AdminUI';
import { adminApi, getErrorMessage } from '@/services';
import { Category } from '@/types/api';
import { useDebouncedValue } from '@/hooks/useDebouncedValue';

export default function CategoriesPage() {
  const [items, setItems] = useState<Category[]>([]);
  const [draft, setDraft] = useState({ name: '', slug: '', description: '' });
  const [queryInput, setQueryInput] = useState('');
  const query = useDebouncedValue(queryInput);
  const [sort, setSort] = useState('sort_order');
  const [order, setOrder] = useState<'asc' | 'desc'>('asc');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    setError('');
    try { setItems((await adminApi.categories({ q: query || undefined, sort, order })).data); }
    catch (loadError) { setError(getErrorMessage(loadError, '分类加载失败')); }
  }, [order, query, sort]);
  useEffect(() => { void Promise.resolve().then(load); }, [load]);

  const update = async (item: Category, change: Partial<Category>) => {
    setBusy(true);
    try { await adminApi.updateCategory(item.id, { ...item, ...change }); await load(); }
    catch (updateError) { setError(getErrorMessage(updateError, '分类更新失败')); }
    finally { setBusy(false); }
  };

  const create = async () => {
    setBusy(true);
    setError('');
    try {
      const allCategories = (await adminApi.categories({ sort: 'sort_order', order: 'asc' })).data;
      const nextSortOrder = allCategories.reduce((maximum, item) => Math.max(maximum, item.sort_order), -1) + 1;
      await adminApi.createCategory({ ...draft, sort_order: nextSortOrder, is_active: true, requires_review: false });
      setDraft({ name: '', slug: '', description: '' });
      await load();
    } catch (createError) {
      setError(getErrorMessage(createError, '新增分类失败'));
    } finally {
      setBusy(false);
    }
  };

  return <div className="mx-auto max-w-5xl space-y-6">
    <PageTitle title="分类管理" subtitle="检索分类并按顺序、名称、状态或创建时间排列" />
    {error && <Notice>{error}</Notice>}
    <Panel>
      <form onSubmit={(event) => { event.preventDefault(); void create(); }} className="grid gap-3 md:grid-cols-4">
        <input required minLength={2} value={draft.name} onChange={(event) => setDraft({ ...draft, name: event.target.value })} placeholder="分类名称" className="rounded-lg border p-2.5" />
        <input required minLength={2} value={draft.slug} onChange={(event) => setDraft({ ...draft, slug: event.target.value })} placeholder="英文 slug" className="rounded-lg border p-2.5" />
        <input value={draft.description} onChange={(event) => setDraft({ ...draft, description: event.target.value })} placeholder="分类说明" className="rounded-lg border p-2.5" />
        <button disabled={busy} className="rounded-lg bg-cyan-600 px-4 py-2 text-white">新增分类</button>
      </form>
      <div className="mt-6">
        <AdminListToolbar query={queryInput} onQueryChange={setQueryInput} searchPlaceholder="搜索分类名称、slug 或说明" sort={sort} onSortChange={setSort} order={order} onOrderChange={setOrder} sortOptions={[{ value: 'sort_order', label: '按手工顺序' }, { value: 'name', label: '按名称' }, { value: 'status', label: '按启用状态' }, { value: 'created_at', label: '按创建时间' }]} />
      </div>
      <div className="space-y-3">
        {items.map((item) => <div key={item.id} className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-slate-200 p-4"><div><strong>{item.name}</strong><span className="ml-2 text-xs text-slate-400">/{item.slug}</span><p className="mt-1 text-xs text-slate-500">{item.description || '暂无说明'} · 排序 {item.sort_order}</p></div><div className="flex gap-2"><button disabled={busy} onClick={() => void update(item, { requires_review: !item.requires_review })} className={`rounded-lg px-3 py-2 text-xs ${item.requires_review ? 'bg-amber-100 text-amber-800' : 'bg-slate-100'}`}>{item.requires_review ? '高争议' : '普通分类'}</button><button disabled={busy} onClick={() => void update(item, { is_active: !item.is_active })} className={`rounded-lg px-3 py-2 text-xs ${item.is_active ? 'bg-emerald-100 text-emerald-800' : 'bg-red-100 text-red-700'}`}>{item.is_active ? '已启用' : '已停用'}</button></div></div>)}
        {items.length === 0 && <div className="rounded-xl border border-dashed p-10 text-center text-sm text-slate-500">没有符合搜索条件的分类</div>}
      </div>
    </Panel>
  </div>;
}
