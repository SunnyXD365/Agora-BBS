'use client';

import { useEffect, useState } from 'react';
import { feedbackApi, getErrorMessage } from '@/services';
import { FeedbackStance, FeedbackSummary, FeedbackTag } from '@/types/api';

const tagLabels: Record<FeedbackTag, string> = {
  logical: '逻辑严密',
  new_perspective: '新视角',
  well_sourced: '引证充分',
  empathetic: '真诚共情',
  factual_concern: '事实存疑',
  reasoning_gap: '推理缺口',
  inappropriate: '表达不当',
  legacy_support: '历史支持',
};

const stanceTags: Record<FeedbackStance, FeedbackTag[]> = {
  support: ['logical', 'new_perspective', 'well_sourced', 'empathetic'],
  challenge: ['factual_concern', 'reasoning_gap', 'inappropriate'],
};

export default function ContextualFeedback({ targetType, targetId, enabled }: { targetType: 'topic' | 'post'; targetId: number; enabled: boolean }) {
  const [summary, setSummary] = useState<FeedbackSummary>({ support: 0, challenge: 0, score: 0, tags: {} });
  const [open, setOpen] = useState(false);
  const [stance, setStance] = useState<FeedbackStance>('support');
  const [tag, setTag] = useState<FeedbackTag>('logical');
  const [reason, setReason] = useState('');
  const [error, setError] = useState('');
  const [saving, setSaving] = useState(false);

  const load = async () => {
    try {
      const result = await feedbackApi.summary(targetType, targetId);
      setSummary(result.data);
      if (result.data.mine) {
        setStance(result.data.mine.stance);
        setTag(result.data.mine.tag);
        setReason(result.data.mine.reason);
      }
    } catch (err: unknown) {
      setError(getErrorMessage(err, '反馈统计加载失败'));
    }
  };

  useEffect(() => {
    let cancelled = false;
    feedbackApi.summary(targetType, targetId)
      .then((result) => {
        if (cancelled) return;
        setSummary(result.data);
        if (result.data.mine) {
          setStance(result.data.mine.stance);
          setTag(result.data.mine.tag);
          setReason(result.data.mine.reason);
        }
      })
      .catch((err: unknown) => { if (!cancelled) setError(getErrorMessage(err, '反馈统计加载失败')); });
    return () => { cancelled = true; };
  }, [targetId, targetType]);

  const chooseStance = (value: FeedbackStance) => {
    setStance(value);
    setTag(stanceTags[value][0]);
  };

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    setSaving(true); setError('');
    try {
      await feedbackApi.upsert({ target_type: targetType, target_id: targetId, stance, tag, reason: reason.trim() });
      await load(); setOpen(false);
    } catch (err: unknown) {
      setError(getErrorMessage(err, '语境反馈提交失败'));
    } finally { setSaving(false); }
  };

  const withdraw = async () => {
    if (!summary.mine) return;
    setSaving(true); setError('');
    try {
      await feedbackApi.withdraw(summary.mine.id);
      setReason(''); await load(); setOpen(false);
    } catch (err: unknown) {
      setError(getErrorMessage(err, '反馈撤回失败'));
    } finally { setSaving(false); }
  };

  return (
    <div className="mt-3 text-xs">
      <div className="flex flex-wrap items-center gap-2 text-[var(--text-muted)]">
        <span className="rounded-full bg-emerald-50 px-2 py-1 text-emerald-800">支持 {summary.support}</span>
        <span className="rounded-full bg-orange-50 px-2 py-1 text-orange-800">质疑 {summary.challenge}</span>
        {Object.entries(summary.tags).filter(([, count]) => Boolean(count)).map(([name, count]) => (
          <span key={name} className="rounded-full border px-2 py-1">{tagLabels[name as FeedbackTag]} {count}</span>
        ))}
        {enabled && <button onClick={() => setOpen((value) => !value)} className="font-semibold text-[var(--accent-ink)] hover:underline">{summary.mine ? '修改我的反馈' : '给出语境反馈'}</button>}
      </div>
      {open && enabled && (
        <form onSubmit={submit} className="mt-3 space-y-2 rounded-lg border border-[var(--border-paper)] bg-stone-50 p-3">
          <div className="flex gap-2">
            <button type="button" onClick={() => chooseStance('support')} className={`rounded px-3 py-1 ${stance === 'support' ? 'bg-emerald-700 text-white' : 'bg-white'}`}>支持</button>
            <button type="button" onClick={() => chooseStance('challenge')} className={`rounded px-3 py-1 ${stance === 'challenge' ? 'bg-orange-700 text-white' : 'bg-white'}`}>质疑</button>
            <select value={tag} onChange={(event) => setTag(event.target.value as FeedbackTag)} className="rounded border bg-white px-2">
              {stanceTags[stance].map((value) => <option key={value} value={value}>{tagLabels[value]}</option>)}
            </select>
          </div>
          <textarea minLength={5} maxLength={120} required value={reason} onChange={(event) => setReason(event.target.value)} placeholder="请用 5–120 字说明理由，帮助对方理解你的判断。" className="w-full rounded border p-2" />
          <div className="flex justify-between">
            <span className="text-[var(--text-muted)]">{reason.trim().length}/120 · 反馈可能接受抽检</span>
            <div className="space-x-2">{summary.mine && <button type="button" disabled={saving} onClick={withdraw} className="text-red-700">撤回</button>}<button disabled={saving || reason.trim().length < 5} className="rounded bg-[var(--accent-ink)] px-3 py-1 text-white disabled:opacity-50">{saving ? '提交中…' : '确认反馈'}</button></div>
          </div>
        </form>
      )}
      {error && <p className="mt-2 text-red-700">{error}</p>}
    </div>
  );
}
