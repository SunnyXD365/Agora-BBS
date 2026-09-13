export function PageTitle({title,subtitle}:{title:string;subtitle:string}){return <header><p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-700">Administration</p><h1 className="mt-1 text-2xl font-bold">{title}</h1><p className="mt-1 text-sm text-slate-500">{subtitle}</p></header>}
export function Notice({children}:{children:React.ReactNode}){return <div className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700">{children}</div>}
export function Panel({children}:{children:React.ReactNode}){return <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">{children}</section>}
export function Th({children}:{children:React.ReactNode}){return <th className="whitespace-nowrap px-3 py-2 text-left text-xs font-semibold text-slate-500">{children}</th>}
export function Td({children}:{children:React.ReactNode}){return <td className="px-3 py-3 align-top text-sm">{children}</td>}

export type AdminSortOption = { value: string; label: string };

export function AdminListToolbar({
  query,
  onQueryChange,
  searchPlaceholder,
  sort,
  onSortChange,
  order,
  onOrderChange,
  sortOptions,
  children,
}: {
  query: string;
  onQueryChange: (value: string) => void;
  searchPlaceholder: string;
  sort: string;
  onSortChange: (value: string) => void;
  order: 'asc' | 'desc';
  onOrderChange: (value: 'asc' | 'desc') => void;
  sortOptions: AdminSortOption[];
  children?: React.ReactNode;
}) {
  return (
    <div className="mb-5 flex flex-wrap items-center gap-2 rounded-xl bg-slate-50 p-3">
      <input
        type="search"
        value={query}
        onChange={(event) => onQueryChange(event.target.value)}
        maxLength={100}
        aria-label="搜索"
        placeholder={searchPlaceholder}
        className="min-w-56 flex-1 rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm"
      />
      {children}
      <select aria-label="排序字段" value={sort} onChange={(event) => onSortChange(event.target.value)} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm">
        {sortOptions.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
      </select>
      <select aria-label="排序方向" value={order} onChange={(event) => onOrderChange(event.target.value as 'asc' | 'desc')} className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm">
        <option value="desc">降序</option>
        <option value="asc">升序</option>
      </select>
    </div>
  );
}
