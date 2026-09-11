'use client';

import Link from 'next/link';
import { Topic } from '@/types/api';

export default function TopicCard({ topic }: { topic: Topic }) {
  return (
    <div className="paper-card rounded-lg p-4 transition-colors hover:border-gray-400">
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 space-y-1.5">
          <Link
            href={`/topics/${topic.id}`}
            className="line-clamp-1 text-base font-semibold text-gray-900 transition-colors hover:text-[var(--accent-ink)]"
          >
            {topic.title}
          </Link>
          <p className="line-clamp-2 text-xs text-[var(--text-muted)]">{topic.content}</p>
        </div>
      </div>

      <div className="mt-3 flex items-center justify-between text-xs text-[var(--text-muted)]">
        <div className="flex items-center gap-3">
          <span className="font-medium text-gray-700">{topic.author_name}</span>
          <span>·</span>
          <span>{new Date(topic.created_at).toLocaleDateString()}</span>
        </div>
        <div className="flex items-center gap-4">
          <span>👀 {topic.view_count}</span>
          <span>💬 {topic.post_count}</span>
          <span>👍 {topic.like_count}</span>
        </div>
      </div>
    </div>
  );
}