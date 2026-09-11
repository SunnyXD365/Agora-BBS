'use client';

import { useState } from 'react';
import { Category } from '@/types/api';
import { getErrorMessage, topicApi } from '@/services';

interface Props {
  isOpen: boolean;
  onClose: () => void;
  categories: Category[];
  onSuccess: () => void;
}

export default function CreateTopicModal({ isOpen, onClose, categories, onSuccess }: Props) {
  const [categoryId, setCategoryId] = useState<number | null>(null);
  const [title, setTitle] = useState('');
  const [claim, setClaim] = useState('');
  const [evidence, setEvidence] = useState('');
  const [uncertainty, setUncertainty] = useState('');
  const [submitting, setSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const selectedCategoryId = categoryId ?? categories[0]?.id;
    if (!selectedCategoryId || !title.trim() || !claim.trim()) return;

    setSubmitting(true);
    try {
      const res = await topicApi.createTopic({
        category_id: selectedCategoryId,
        title,
        structured_content: { claim, evidence, uncertainty },
      });
      if (res.code === 0) {
        setTitle('');
        setClaim('');
        setEvidence('');
        setUncertainty('');
        onSuccess();
        onClose();
      }
    } catch (err: unknown) {
      alert(getErrorMessage(err, '发帖失败'));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm p-4">
      <div className="w-full max-w-lg rounded-xl bg-white p-6 shadow-xl">
        <h2 className="text-lg font-bold text-gray-900">发布新帖子</h2>
        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="block text-xs font-medium text-gray-700">所属板块</label>
            <select
              value={categoryId ?? categories[0]?.id ?? ''}
              onChange={(e) => setCategoryId(Number(e.target.value))}
              className="mt-1 w-full rounded-md border border-gray-300 p-2 text-sm focus:border-blue-500 focus:outline-none"
            >
              {categories.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700">标题</label>
            <input
              type="text"
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="mt-1 w-full rounded-md border border-gray-300 p-2 text-sm focus:border-blue-500 focus:outline-none"
              placeholder="简要概括主题..."
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700">我的核心观点</label>
            <textarea
              required
              rows={5}
              value={claim}
              onChange={(e) => setClaim(e.target.value)}
              className="mt-1 w-full rounded-md border border-gray-300 p-2 text-sm focus:border-blue-500 focus:outline-none"
              placeholder="清晰写出你希望讨论的核心观点..."
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700">我支持它的依据</label>
            <textarea rows={3} value={evidence} onChange={(e) => setEvidence(e.target.value)} className="mt-1 w-full rounded-md border p-2 text-sm" placeholder="事实、数据、引用或推理依据..." />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700">我还不确定的地方</label>
            <textarea rows={2} value={uncertainty} onChange={(e) => setUncertainty(e.target.value)} className="mt-1 w-full rounded-md border p-2 text-sm" placeholder="主动说明疑问和可能的局限..." />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="rounded-md border px-4 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-50"
            >
              取消
            </button>
            <button
              type="submit"
              disabled={submitting || categories.length === 0}
              className="rounded-md bg-blue-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-blue-500 disabled:opacity-50"
            >
              {submitting ? '发布中...' : '提交'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
