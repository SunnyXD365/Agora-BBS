type PaginationProps = {
  page: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  itemLabel?: string;
  disabled?: boolean;
};

export default function Pagination({ page, pageSize, total, onPageChange, itemLabel = '条记录', disabled = false }: PaginationProps) {
  if (total <= 0) return null;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const current = Math.min(Math.max(1, page), totalPages);
  const start = Math.max(1, Math.min(current - 2, totalPages - 4));
  const end = Math.min(totalPages, start + 4);
  const pages = Array.from({ length: end - start + 1 }, (_, index) => start + index);

  return (
    <nav aria-label="分页导航" className="paper-card flex flex-wrap items-center justify-between gap-3 rounded-xl px-4 py-3 text-sm">
      <span className="text-xs text-[var(--text-muted)]">共 {total} {itemLabel} · 第 {current}/{totalPages} 页</span>
      <div className="flex items-center gap-1">
        <button type="button" disabled={disabled || current === 1} onClick={() => onPageChange(current - 1)} className="rounded border px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40">上一页</button>
        {start > 1 && <span className="px-1 text-[var(--text-muted)]">…</span>}
        {pages.map((item) => <button type="button" key={item} disabled={disabled} aria-current={item === current ? 'page' : undefined} onClick={() => onPageChange(item)} className={`min-w-8 rounded px-2 py-1.5 ${item === current ? 'bg-[var(--accent-ink)] text-white' : 'border hover:bg-stone-100'}`}>{item}</button>)}
        {end < totalPages && <span className="px-1 text-[var(--text-muted)]">…</span>}
        <button type="button" disabled={disabled || current === totalPages} onClick={() => onPageChange(current + 1)} className="rounded border px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40">下一页</button>
      </div>
    </nav>
  );
}
