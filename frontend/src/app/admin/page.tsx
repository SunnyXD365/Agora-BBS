'use client';

import { useEffect, useState } from 'react';
import { adminApi, getErrorMessage } from '@/services';
import { AdminOverview } from '@/types/api';
import { Notice, PageTitle } from '@/components/admin/AdminUI';

export default function AdminDashboard() {
  const [data, setData] = useState<AdminOverview | null>(null); const [days, setDays] = useState(7); const [error, setError] = useState(''); const [loading, setLoading] = useState(true);
  useEffect(() => { let cancelled=false; adminApi.overview(days).then(r => { if(!cancelled) setData(r.data); }).catch(e => { if(!cancelled) setError(getErrorMessage(e,'统计数据加载失败')); }).finally(() => { if(!cancelled) setLoading(false); }); return () => { cancelled=true; }; }, [days]);
  if (loading && !data) return <div className="h-56 animate-pulse rounded-2xl bg-white" />;
  const cards = data ? [['用户总数',data.users_total],['今日新增',data.new_users_today],['7 日活跃',data.active_users_7_days],['主题 / 回复',`${data.topics_total} / ${data.posts_total}`],['语境反馈',data.feedback_total],['收藏记录',data.bookmarks_total],['有效阅读',`${data.verified_read_hours.toFixed(1)} h`],['停用账号',data.suspended_users]] : [];
  const max = Math.max(1,...(data?.trend.map(x=>x.users+x.topics+x.posts+x.feedback)||[1]));
  return <div className="mx-auto max-w-7xl space-y-6"><PageTitle title="运营总览" subtitle="社区增长、参与质量、治理状态与模型运行情况" />
    {error && <Notice>{error}</Notice>}
    <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">{cards.map(([k,v])=><Metric key={k} label={String(k)} value={v}/>)}</section>
    {data && <><section className="rounded-2xl border border-slate-200 bg-white p-5"><div className="flex items-center justify-between"><div><h2 className="font-semibold">内容与互动趋势</h2><p className="text-xs text-slate-500">用户、主题、回复和反馈的每日新增总和</p></div><div className="flex gap-1">{[7,30,90].map(x=><button key={x} onClick={()=>setDays(x)} className={`rounded-lg px-3 py-1.5 text-xs ${days===x?'bg-slate-900 text-white':'bg-slate-100'}`}>{x} 天</button>)}</div></div><div className="mt-8 flex h-52 items-end gap-1 overflow-hidden">{data.trend.map(day=>{const total=day.users+day.topics+day.posts+day.feedback;return <div key={day.date} className="group flex min-w-2 flex-1 flex-col items-center justify-end" title={`${day.date}：${total} 次新增`}><div className="w-full max-w-8 rounded-t bg-cyan-500 transition hover:bg-cyan-400" style={{height:`${Math.max(4,total/max*170)}px`}}/><span className="mt-2 hidden text-[9px] text-slate-400 sm:block">{data.trend.length<=30?day.date.slice(5):''}</span></div>})}</div></section>
    <div className="grid gap-6 xl:grid-cols-2"><Breakdown title="用户等级分布" data={data.level_distribution}/><Breakdown title="信任分分布" data={data.trust_distribution}/><Breakdown title="反馈标签分布" data={data.feedback_distribution}/><section className="rounded-2xl border border-slate-200 bg-white p-5"><h2 className="font-semibold">分类内容规模</h2><div className="mt-4 space-y-3">{data.categories.map(x=><div key={x.name} className="flex items-center justify-between border-b border-slate-100 pb-2 text-sm"><span>{x.name}</span><span className="text-slate-500">{x.topics} 主题 · {x.posts} 回复</span></div>)}</div></section></div>
    <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4"><Metric label="盲审完成率" value={percent(data.review_completed,data.review_total)}/><Metric label="盲审超时" value={data.review_expired}/><Metric label="LLM 成功率" value={percent(data.llm_success,data.llm_calls)}/><Metric label="LLM 平均延迟 / Token" value={`${Math.round(data.llm_average_ms)} ms / ${data.llm_prompt_tokens+data.llm_output_tokens}`}/></section></>}
  </div>;
}

function Metric({label,value}:{label:string;value:React.ReactNode}){return <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"><p className="text-xs text-slate-500">{label}</p><p className="mt-2 text-2xl font-bold tracking-tight">{value}</p></div>}
function Breakdown({title,data}:{title:string;data:Record<string,number>}){const total=Math.max(1,Object.values(data).reduce((a,b)=>a+b,0));return <section className="rounded-2xl border border-slate-200 bg-white p-5"><h2 className="font-semibold">{title}</h2><div className="mt-4 space-y-3">{Object.entries(data).map(([k,v])=><div key={k}><div className="flex justify-between text-xs"><span>{k}</span><span>{v}</span></div><div className="mt-1 h-2 rounded bg-slate-100"><div className="h-2 rounded bg-cyan-500" style={{width:`${v/total*100}%`}}/></div></div>)}</div></section>}
function percent(a:number,b:number){return b?`${(a/b*100).toFixed(1)}%`:'—'}
