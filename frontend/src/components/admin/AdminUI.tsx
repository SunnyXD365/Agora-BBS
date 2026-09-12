export function PageTitle({title,subtitle}:{title:string;subtitle:string}){return <header><p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-700">Administration</p><h1 className="mt-1 text-2xl font-bold">{title}</h1><p className="mt-1 text-sm text-slate-500">{subtitle}</p></header>}
export function Notice({children}:{children:React.ReactNode}){return <div className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700">{children}</div>}
export function Panel({children}:{children:React.ReactNode}){return <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">{children}</section>}
export function Th({children}:{children:React.ReactNode}){return <th className="whitespace-nowrap px-3 py-2 text-left text-xs font-semibold text-slate-500">{children}</th>}
export function Td({children}:{children:React.ReactNode}){return <td className="px-3 py-3 align-top text-sm">{children}</td>}
