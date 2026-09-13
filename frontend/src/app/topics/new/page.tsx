'use client';

import Link from 'next/link';
import { Suspense, useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { getErrorMessage, topicApi } from '@/services';
import { Category, StructuredContent } from '@/types/api';

type EditorData = { category_id: number; title: string; structured_content: StructuredContent };

export default function NewTopicPage() {
  return <Suspense fallback={<div className="mx-auto h-72 max-w-3xl animate-pulse rounded-xl bg-stone-200" />}><NewTopicEditor /></Suspense>;
}

function NewTopicEditor() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const requestedDraftID = Number(searchParams.get('draft') || 0);
  const { user, isLoading } = useAuth();
  const [categories, setCategories] = useState<Category[]>([]);
  const [categoryID, setCategoryID] = useState(0);
  const [title, setTitle] = useState('');
  const [claim, setClaim] = useState('');
  const [evidence, setEvidence] = useState('');
  const [uncertainty, setUncertainty] = useState('');
  const [draftID, setDraftID] = useState(requestedDraftID);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [publishing, setPublishing] = useState(false);
  const [notice, setNotice] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isLoading && !user) router.replace('/login');
  }, [isLoading, router, user]);

  useEffect(() => {
    if (!user) return;
    let cancelled = false;
    Promise.all([
      topicApi.getCategories(),
      requestedDraftID ? topicApi.getTopicDetail(requestedDraftID) : Promise.resolve(null),
    ]).then(([categoryResult, draftResult]) => {
      if (cancelled) return;
      setCategories(categoryResult.data);
      setCategoryID(draftResult?.data.category_id || categoryResult.data[0]?.id || 0);
      if (draftResult) {
        if (draftResult.data.status !== 'draft' || draftResult.data.user_id !== user.id) throw new Error('该草稿不可编辑');
        setTitle(draftResult.data.title);
        setClaim(draftResult.data.structured_content?.claim || draftResult.data.content || '');
        setEvidence(draftResult.data.structured_content?.evidence || '');
        setUncertainty(draftResult.data.structured_content?.uncertainty || '');
        setDraftID(draftResult.data.id);
      }
    }).catch((err: unknown) => setError(getErrorMessage(err, err instanceof Error ? err.message : '编辑器加载失败')))
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, [requestedDraftID, user]);

  const data = (): EditorData => ({
    category_id: categoryID,
    title: title.trim(),
    structured_content: { claim: claim.trim(), evidence: evidence.trim(), uncertainty: uncertainty.trim() },
  });

  const saveDraft = async () => {
    if (!categoryID) { setError('请先选择所属板块。'); return 0; }
    setSaving(true); setError(''); setNotice('');
    try {
      const result = draftID ? await topicApi.updateDraft(draftID, data()) : await topicApi.createDraft(data());
      setDraftID(result.data.id);
      setNotice('草稿已安全保存，可以稍后在“我的内容”中继续编辑。');
      if (!draftID) router.replace(`/topics/new?draft=${result.data.id}`);
      return result.data.id;
    } catch (err: unknown) {
      setError(getErrorMessage(err, '草稿保存失败'));
      return 0;
    } finally { setSaving(false); }
  };

  const publish = async (event: React.FormEvent) => {
    event.preventDefault(); setError(''); setNotice('');
    if (title.trim().length < 3) { setError('标题至少需要 3 个字。'); return; }
    if (claim.trim().length < 5) { setError('核心观点至少需要 5 个字。'); return; }
    if (!user?.capabilities.includes('create_topic')) { setError('发布主题需要达到 L2；你仍然可以先保存草稿。'); return; }
    setPublishing(true);
    try {
      let result;
      if (draftID) {
        await topicApi.updateDraft(draftID, data());
        result = await topicApi.publishDraft(draftID);
      } else {
        result = await topicApi.createTopic(data());
      }
      router.push(`/topics/${result.data.topic_id}`);
    } catch (err: unknown) { setError(getErrorMessage(err, '主题提交失败')); }
    finally { setPublishing(false); }
  };

  if (isLoading || loading) return <div className="mx-auto h-72 max-w-3xl animate-pulse rounded-xl bg-stone-200" />;
  if (!user) return null;

  return (
    <div className="grid gap-6 xl:grid-cols-[220px_minmax(0,1fr)_220px]">
      <aside className="space-y-4">
        <TipCard title="先搭好结构" items={['观点：你真正主张什么', '依据：事实、数据或推理', '不确定：主动说明边界']} />
        <TipCard title="让别人容易回应" items={['标题具体，避免只有情绪', '一个主题聚焦一个问题', '引用资料时注明来源']} />
      </aside>

      <main className="paper-card rounded-2xl p-6 sm:p-8">
        <div className="flex flex-wrap items-start justify-between gap-3 border-b border-[var(--border-paper)] pb-5">
          <div><p className="text-xs font-semibold tracking-widest text-[var(--accent-ink)]">结构化发帖</p><h1 className="mt-1 text-2xl font-bold">{draftID ? `继续编辑草稿 #${draftID}` : '发起一个新主题'}</h1></div>
          <Link href="/my-content?status=draft" className="text-sm text-[var(--accent-ink)] hover:underline">查看我的草稿</Link>
        </div>
        {notice && <p className="mt-4 rounded-lg border border-emerald-200 bg-emerald-50 p-3 text-sm text-emerald-800">{notice}</p>}
        {error && <p className="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-700">{error}</p>}
        <form onSubmit={publish} className="mt-6 space-y-5">
          <Field label="所属板块"><select value={categoryID} onChange={(event) => setCategoryID(Number(event.target.value))} className="mt-2 w-full rounded-lg border p-3 text-sm">{categories.map((category) => <option key={category.id} value={category.id}>{category.name}{category.requires_review ? ' · 高争议/需盲审' : ''}</option>)}</select></Field>
          <Field label="标题" hint={`${title.length}/128`}><input required maxLength={128} value={title} onChange={(event) => setTitle(event.target.value)} placeholder="用一句具体的话概括你想讨论的问题" className="mt-2 w-full rounded-lg border p-3" /></Field>
          <Field label="我的核心观点" hint={`${claim.length} 字`}><textarea required rows={8} value={claim} onChange={(event) => setClaim(event.target.value)} placeholder="清晰说明你的判断、主张或希望解决的问题……" className="mt-2 w-full rounded-lg border p-3 text-sm leading-7" /></Field>
          <Field label="支持依据（推荐填写）"><textarea rows={5} value={evidence} onChange={(event) => setEvidence(event.target.value)} placeholder="事实、数据、出处、案例或推理过程……" className="mt-2 w-full rounded-lg border p-3 text-sm leading-7" /></Field>
          <Field label="我还不确定的地方"><textarea rows={4} value={uncertainty} onChange={(event) => setUncertainty(event.target.value)} placeholder="哪些条件可能改变结论？你希望其他人补充什么？" className="mt-2 w-full rounded-lg border p-3 text-sm leading-7" /></Field>
          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-[var(--border-paper)] pt-5">
            <span className="text-xs text-[var(--text-muted)]">草稿不会公开；正式提交后进入冷静期。</span>
            <div className="flex gap-3"><button type="button" disabled={saving || publishing || categories.length === 0} onClick={saveDraft} className="rounded-lg border px-5 py-2 text-sm disabled:opacity-50">{saving ? '保存中…' : '保存草稿'}</button><button type="submit" disabled={publishing || saving || categories.length === 0} className="paper-btn-primary rounded-lg px-5 py-2 text-sm disabled:opacity-50">{publishing ? '提交中…' : '提交并进入冷静期'}</button></div>
          </div>
        </form>
      </main>

      <aside className="space-y-4">
        <TipCard title="发布前检查" items={['区分事实、判断与个人经历', '避免泄露隐私和敏感信息', '不同意观点，也尊重表达者']} />
        <section className="paper-card rounded-xl p-4 text-sm leading-6"><h2 className="font-bold">发布流程</h2><p className="mt-2 text-[var(--text-muted)]">提交后先进入冷静期，你可以编辑或无痕撤回。长文及高争议板块在冷静期后还会进入匿名盲审。</p><Link href="/guide#topics" className="mt-3 inline-block text-[var(--accent-ink)] hover:underline">查看完整规则 →</Link></section>
      </aside>
    </div>
  );
}

function TipCard({ title, items }: { title: string; items: string[] }) {
  return <section className="paper-card rounded-xl p-4"><h2 className="font-bold">{title}</h2><ul className="mt-3 space-y-2 text-sm leading-6 text-[var(--text-muted)]">{items.map((item) => <li key={item} className="flex gap-2"><span className="text-[var(--accent-ink)]">•</span><span>{item}</span></li>)}</ul></section>;
}

function Field({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return <label className="block"><span className="flex items-center justify-between text-sm font-bold"><span>{label}</span>{hint && <span className="font-normal text-[var(--text-muted)]">{hint}</span>}</span>{children}</label>;
}
