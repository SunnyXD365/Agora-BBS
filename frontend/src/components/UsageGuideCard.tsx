import Link from 'next/link';

type UsageGuideCardProps = {
  context?: 'home' | 'growth';
};

const homeSteps = [
  ['先阅读', '打开主题并认真读完，长文会记录滚动与停留进度。'],
  ['再参与', 'L1 可回复和提交语境反馈，反馈需要说明理由。'],
  ['后发起', '达到 L2 后，用观点、依据、不确定性创建结构化主题。'],
];

const growthSteps = [
  ['L0 阅读者', '可浏览和收藏，通过有效阅读积累成长进度。'],
  ['L1 参与者', '解锁回复与语境反馈。'],
  ['L2 发起者', '解锁结构化主题发布。'],
  ['L3 评审者', '解锁匿名盲审，帮助维护讨论质量。'],
];

export default function UsageGuideCard({ context = 'home' }: UsageGuideCardProps) {
  const steps = context === 'growth' ? growthSteps : homeSteps;
  return (
    <aside className="paper-card rounded-xl p-4 text-sm lg:sticky lg:top-6">
      <p className="text-xs uppercase tracking-[0.18em] text-[var(--accent-ink)]">新手使用说明</p>
      <h2 className="mt-2 text-lg font-bold">{context === 'growth' ? '如何获得权限' : '第一次来到 Agora'}</h2>
      <div className="mt-4 space-y-4">
        {steps.map(([title, description], index) => (
          <div key={title} className="flex gap-3">
            <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-stone-100 text-xs font-bold text-[var(--accent-ink)]">{index + 1}</span>
            <div><h3 className="font-semibold">{title}</h3><p className="mt-1 text-xs leading-5 text-[var(--text-muted)]">{description}</p></div>
          </div>
        ))}
      </div>
      <div className="mt-5 border-t border-[var(--border-paper)] pt-4">
        <Link href="/guide" className="font-semibold text-[var(--accent-ink)] hover:underline">查看完整论坛使用手册 →</Link>
      </div>
    </aside>
  );
}
