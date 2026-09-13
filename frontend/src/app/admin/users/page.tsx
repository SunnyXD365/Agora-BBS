'use client';

import { useCallback, useEffect, useState } from 'react';
import Pagination from '@/components/Pagination';
import { AdminListToolbar, Notice, PageTitle, Panel, Td, Th } from '@/components/admin/AdminUI';
import { adminApi, getErrorMessage } from '@/services';
import { AdminUser } from '@/types/api';
import { useAuth } from '@/context/AuthContext';
import { formatReadingDuration } from '@/lib/format';
import { useDebouncedValue } from '@/hooks/useDebouncedValue';

type EditDraft = {
  username: string;
  email: string;
  role: 'user' | 'admin';
  status: 'active' | 'suspended';
  unlock_level: number;
  trust_score: number;
  reason: string;
};

export default function UsersPage() {
  const { user } = useAuth();
  const [items, setItems] = useState<AdminUser[]>([]);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [queryInput, setQueryInput] = useState('');
  const query = useDebouncedValue(queryInput);
  const [level, setLevel] = useState('');
  const [role, setRole] = useState('');
  const [status, setStatus] = useState('');
  const [sort, setSort] = useState('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [busy, setBusy] = useState(false);
  const [editing, setEditing] = useState<AdminUser | null>(null);
  const [editDraft, setEditDraft] = useState<EditDraft | null>(null);

  const load = useCallback(async () => {
    setBusy(true);
    setError('');
    try {
      const result = await adminApi.users({
        q: query || undefined,
        level: level === '' ? undefined : Number(level),
        role: role || undefined,
        status: status || undefined,
        sort,
        order,
        page,
        page_size: pageSize,
      });
      setItems(result.data.items);
      setTotal(result.data.total);
    } catch (loadError) {
      setError(getErrorMessage(loadError, '用户数据加载失败'));
    } finally {
      setBusy(false);
    }
  }, [level, order, page, pageSize, query, role, sort, status]);

  useEffect(() => { void Promise.resolve().then(load); }, [load]);
  const resetPage = () => setPage(1);

  const beginEdit = (item: AdminUser) => {
    setError('');
    setSuccess('');
    setEditing(item);
    setEditDraft({
      username: item.username,
      email: item.email || '',
      role: item.role as 'user' | 'admin',
      status: item.status as 'active' | 'suspended',
      unlock_level: item.unlock_level,
      trust_score: item.trust_score,
      reason: '',
    });
  };

  const saveEdit = async () => {
    if (!editing || !editDraft) return;
    setBusy(true);
    setError('');
    setSuccess('');
    try {
      await adminApi.updateUser(editing.id, editDraft);
      setEditing(null);
      setEditDraft(null);
      setSuccess(`用户 ${editDraft.username} 已更新，调整记录已写入信任流水。`);
      await load();
    } catch (updateError) {
      setError(getErrorMessage(updateError, '用户更新失败'));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="mx-auto max-w-7xl space-y-6">
      <PageTitle title="用户管理" subtitle="按账号、等级、信任、阅读与注册时间检索和排序" />
      {error && <Notice>{error}</Notice>}
      {success && <div className="rounded-xl border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-700">{success}</div>}
      <Panel>
        <AdminListToolbar
          query={queryInput}
          onQueryChange={(value) => { setQueryInput(value); resetPage(); }}
          searchPlaceholder="搜索用户名或邮箱"
          sort={sort}
          onSortChange={(value) => { setSort(value); resetPage(); }}
          order={order}
          onOrderChange={(value) => { setOrder(value); resetPage(); }}
          sortOptions={[
            { value: 'created_at', label: '按注册时间' },
            { value: 'level', label: '按等级' },
            { value: 'trust_score', label: '按信任分' },
            { value: 'reading_seconds', label: '按阅读时长' },
            { value: 'username', label: '按用户名' },
          ]}
        >
          <select aria-label="等级筛选" value={level} onChange={(event) => { setLevel(event.target.value); resetPage(); }} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm">
            <option value="">全部等级</option>
            {[0, 1, 2, 3].map((item) => <option key={item} value={item}>L{item}</option>)}
          </select>
          <select aria-label="身份筛选" value={role} onChange={(event) => { setRole(event.target.value); resetPage(); }} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm">
            <option value="">全部身份</option><option value="user">普通用户</option><option value="admin">管理员</option>
          </select>
          <select aria-label="账号状态筛选" value={status} onChange={(event) => { setStatus(event.target.value); resetPage(); }} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm">
            <option value="">全部状态</option><option value="active">正常</option><option value="suspended">已停用</option>
          </select>
        </AdminListToolbar>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead><tr><Th>用户</Th><Th>身份</Th><Th>等级</Th><Th>信任分</Th><Th>阅读</Th><Th>自述</Th><Th>注册时间</Th><Th>操作</Th></tr></thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id} className="border-t border-slate-100">
                  <Td><strong>{item.username}</strong><div className="text-xs text-slate-500">{item.email || '未填写邮箱'}</div></Td>
                  <Td>{item.role} / {item.status}</Td><Td>L{item.unlock_level}</Td><Td>{item.trust_score}</Td>
                  <Td>{formatReadingDuration(item.verified_read_seconds)}</Td><Td>{item.onboarding_status}</Td><Td>{new Date(item.created_at).toLocaleDateString()}</Td>
                  <Td><div className="flex gap-3"><button disabled={busy} onClick={() => beginEdit(item)} className="font-medium text-cyan-700 hover:underline">编辑</button>{item.id !== user?.id && <button disabled={busy} onClick={async () => { setBusy(true); setSuccess(''); try { await adminApi.setUserStatus(item.id, item.status === 'active' ? 'suspended' : 'active'); await load(); } catch (actionError) { setError(getErrorMessage(actionError, '操作失败')); } finally { setBusy(false); } }} className="font-medium text-slate-600 hover:underline">{item.status === 'active' ? '停用' : '恢复'}</button>}</div></Td>
                </tr>
              ))}
              {!busy && items.length === 0 && <tr><td colSpan={8} className="border-t px-3 py-10 text-center text-sm text-slate-500">没有符合筛选条件的用户</td></tr>}
            </tbody>
          </table>
        </div>
        <div className="mt-4"><Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} onPageSizeChange={(size) => { setPage(1); setPageSize(size); }} itemLabel="个用户" disabled={busy} /></div>
      </Panel>
      {editing && editDraft && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/45 p-4" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget && !busy) { setEditing(null); setEditDraft(null); } }}>
          <form role="dialog" aria-modal="true" aria-labelledby="edit-user-title" onSubmit={(event) => { event.preventDefault(); void saveEdit(); }} className="max-h-[90vh] w-full max-w-2xl overflow-y-auto rounded-2xl bg-white p-6 shadow-2xl">
            <div className="flex items-start justify-between gap-4">
              <div><h2 id="edit-user-title" className="text-xl font-bold">编辑用户</h2><p className="mt-1 text-sm text-slate-500">修改身份、等级和信任分；每次保存都会留下管理员调整流水。</p></div>
              <button type="button" disabled={busy} onClick={() => { setEditing(null); setEditDraft(null); }} aria-label="关闭编辑窗口" className="rounded-lg px-3 py-1.5 text-slate-500 hover:bg-slate-100">关闭</button>
            </div>
            <div className="mt-5 grid gap-4 sm:grid-cols-2">
              <label className="text-sm font-medium text-slate-700">用户名<input aria-label="编辑用户名" required minLength={3} maxLength={32} value={editDraft.username} onChange={(event) => setEditDraft({ ...editDraft, username: event.target.value })} className="mt-1.5 w-full rounded-lg border border-slate-200 px-3 py-2.5" /></label>
              <label className="text-sm font-medium text-slate-700">邮箱<input aria-label="编辑邮箱" type="email" maxLength={128} value={editDraft.email} onChange={(event) => setEditDraft({ ...editDraft, email: event.target.value })} className="mt-1.5 w-full rounded-lg border border-slate-200 px-3 py-2.5" /></label>
              <label className="text-sm font-medium text-slate-700">身份<select aria-label="编辑身份" disabled={editing.id === user?.id} value={editDraft.role} onChange={(event) => setEditDraft({ ...editDraft, role: event.target.value as EditDraft['role'] })} className="mt-1.5 w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5"><option value="user">普通用户</option><option value="admin">管理员</option></select></label>
              <label className="text-sm font-medium text-slate-700">账号状态<select aria-label="编辑账号状态" disabled={editing.id === user?.id} value={editDraft.status} onChange={(event) => setEditDraft({ ...editDraft, status: event.target.value as EditDraft['status'] })} className="mt-1.5 w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5"><option value="active">正常</option><option value="suspended">已停用</option></select></label>
              <label className="text-sm font-medium text-slate-700">等级<select aria-label="编辑等级" value={editDraft.unlock_level} onChange={(event) => setEditDraft({ ...editDraft, unlock_level: Number(event.target.value) })} className="mt-1.5 w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5">{[0, 1, 2, 3].map((levelValue) => <option key={levelValue} value={levelValue}>L{levelValue}</option>)}</select></label>
              <label className="text-sm font-medium text-slate-700">信任分<input aria-label="编辑信任分" required type="number" min={-1000} max={10000} value={editDraft.trust_score} onChange={(event) => setEditDraft({ ...editDraft, trust_score: Number(event.target.value) })} className="mt-1.5 w-full rounded-lg border border-slate-200 px-3 py-2.5" /></label>
            </div>
            <label className="mt-4 block text-sm font-medium text-slate-700">调整原因<textarea aria-label="调整原因" required minLength={5} maxLength={200} value={editDraft.reason} onChange={(event) => setEditDraft({ ...editDraft, reason: event.target.value })} placeholder="说明调整依据，5–200 字" className="mt-1.5 min-h-24 w-full rounded-lg border border-slate-200 px-3 py-2.5" /></label>
            <p className="mt-3 text-xs leading-5 text-slate-500">有效阅读时长来自服务端阅读会话，不能在这里修改。被提升为管理员的账号下次登录时必须完成邮箱验证码校验。</p>
            <div className="mt-5 flex justify-end gap-3"><button type="button" disabled={busy} onClick={() => { setEditing(null); setEditDraft(null); }} className="rounded-lg border border-slate-200 px-4 py-2 text-sm">取消</button><button disabled={busy} className="rounded-lg bg-cyan-700 px-4 py-2 text-sm font-medium text-white disabled:opacity-50">{busy ? '保存中…' : '保存修改'}</button></div>
          </form>
        </div>
      )}
    </div>
  );
}
