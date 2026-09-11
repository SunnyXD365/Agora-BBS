import Link from 'next/link';

const sections = [
  ['quick-start', '第一次使用'],
  ['growth', '成长与权限'],
  ['reading', '有效阅读'],
  ['topics', '发起主题'],
  ['discussion', '回复与反馈'],
  ['cooling', '冷静期'],
  ['review', '匿名盲审'],
  ['status', '内容状态'],
  ['faq', '常见问题'],
];

function Section({ id, title, children }: { id: string; title: string; children: React.ReactNode }) {
  return <section id={id} className="paper-card scroll-mt-24 rounded-xl p-6"><h2 className="text-xl font-bold">{title}</h2><div className="mt-4 space-y-3 text-sm leading-7 text-[var(--text-muted)]">{children}</div></section>;
}

export default function GuidePage() {
  return (
    <div className="mx-auto grid max-w-6xl gap-8 lg:grid-cols-[220px_minmax(0,1fr)]">
      <aside className="h-fit lg:sticky lg:top-6">
        <p className="text-xs uppercase tracking-[0.2em] text-[var(--accent-ink)]">Agora BBS</p>
        <h1 className="mt-2 text-2xl font-bold">论坛使用手册</h1>
        <p className="mt-2 text-sm leading-6 text-[var(--text-muted)]">从第一次阅读，到参与讨论、发起主题和承担评审。</p>
        <nav className="mt-6 space-y-1 border-l border-[var(--border-paper)] pl-4">
          {sections.map(([id, title]) => <a key={id} href={`#${id}`} className="block py-1.5 text-sm hover:text-[var(--accent-ink)]">{title}</a>)}
        </nav>
        <Link href="/" className="mt-6 inline-block text-sm font-semibold text-[var(--accent-ink)] hover:underline">← 返回主题列表</Link>
      </aside>

      <article className="space-y-6">
        <div className="rounded-xl border border-amber-200 bg-amber-50 p-6">
          <h2 className="font-bold text-amber-950">一分钟了解这里</h2>
          <p className="mt-2 text-sm leading-7 text-amber-900">Agora 鼓励“先读、再回应、说明理由”。权限会随有效阅读和合规参与逐步开放；发出的内容先经过冷静期，重要长文还会接受匿名盲审。</p>
        </div>

        <Section id="quick-start" title="第一次使用">
          <ol className="list-decimal space-y-2 pl-5"><li>注册或登录账号；未登录也可以浏览公开主题。</li><li>进入一个感兴趣的主题，阅读观点、依据和作者声明的不确定性。</li><li>登录后可以收藏；到“成长中心”提交社区自述并查看当前等级。</li><li>完成有效阅读后逐步解锁回复、语境反馈、发帖和匿名评审。</li></ol>
        </Section>

        <Section id="growth" title="成长与权限">
          <div className="grid gap-3 sm:grid-cols-2">
            {[['L0 阅读者', '浏览主题、收藏内容'], ['L1 讨论参与者', '回复、语境反馈'], ['L2 主题发起者', '创建结构化主题'], ['L3 社区评审者', '参与匿名盲审']].map(([level, ability]) => <div key={level} className="rounded-lg border border-[var(--border-paper)] p-3"><strong className="text-[var(--text-main)]">{level}</strong><p>{ability}</p></div>)}
          </div>
          <p>升级综合考虑注册时间、有效阅读、合规互动和社区信任。普通用户看不到隐藏信任分，只会看到已解锁能力和下一阶段提示；请以成长中心显示为准。</p>
          <p>社区自述不会阻止注册或浏览。评审只判断表达是否得体、参与意愿是否真诚，不判断你的观点是否“正确”。</p>
        </Section>

        <Section id="reading" title="什么算有效阅读">
          <p>打开主题后，系统会建立阅读会话并定期记录进度。请正常浏览内容，不要反复刷新或直接拖动进度。</p>
          <p>长文或高争议分类在回复前需要滚动到底，并在回复区域停留一段时间。满足条件后，页面会提示可以参与回复；累计有效时间也会同步到成长中心。</p>
        </Section>

        <Section id="topics" title="如何发起一个好主题">
          <p>L2 用户可以点击首页“发新帖”。主题由三个部分组成：</p>
          <ul className="list-disc space-y-2 pl-5"><li><strong className="text-[var(--text-main)]">观点：</strong>你希望讨论的核心判断。</li><li><strong className="text-[var(--text-main)]">依据：</strong>事实、数据、经历或推理过程。</li><li><strong className="text-[var(--text-main)]">不确定性：</strong>你尚不确定、希望他人补充或可能出错的部分。</li></ul>
          <p>选择最匹配的分类。高争议分类和重要长文可能在冷静期后进入匿名盲审。</p>
        </Section>

        <Section id="discussion" title="回复与语境反馈">
          <p>回复前选择“讨论观点、补充证据、分享经验、表达感谢”之一，并可针对任意已有回复继续嵌套讨论。</p>
          <p>语境反馈不是简单点赞：先选择支持或质疑，再选择逻辑严密、新视角、引证充分、事实存疑等标签，并用 5–120 字说明理由。同一用户对同一内容只保留一条有效反馈，可以修改或撤回。</p>
          <p>收藏用于稍后阅读，可在顶部“收藏”页面统一查看。</p>
        </Section>

        <Section id="cooling" title="冷静期如何工作">
          <p>主题和回复提交后先进入冷静期。此时只有作者本人可见，可以继续编辑或无痕撤回；编辑会重新开始倒计时。</p>
          <p>冷静期结束后，普通内容自动公开，重要长文或高争议内容转为待盲审。请不要因为首页暂时看不到而重复提交。</p>
        </Section>

        <Section id="review" title="匿名盲审怎么做">
          <p>L3 用户会在“匿名盲审”看到分配给自己的任务。评审人看不到其他评审人的选择，应独立判断语言是否得体、表达是否真诚，并写出具体理由。</p>
          <p>通常由三名活跃评审者参与，达到通过票数后即可形成结果；人数不足或超时会由模型作兜底判断。恶意或敷衍评审会影响后续社区权限。</p>
        </Section>

        <Section id="status" title="常见内容状态">
          <dl className="grid gap-2 sm:grid-cols-[140px_1fr]"><dt className="font-semibold text-[var(--text-main)]">cooling</dt><dd>冷静期，仅作者可见。</dd><dt className="font-semibold text-[var(--text-main)]">pending_review</dt><dd>等待匿名盲审。</dd><dt className="font-semibold text-[var(--text-main)]">published</dt><dd>已公开发布。</dd><dt className="font-semibold text-[var(--text-main)]">rejected</dt><dd>评审未通过。</dd><dt className="font-semibold text-[var(--text-main)]">recalled</dt><dd>作者在冷静期撤回。</dd><dt className="font-semibold text-[var(--text-main)]">hidden</dt><dd>因治理需要暂时隐藏。</dd></dl>
        </Section>

        <Section id="faq" title="常见问题">
          <div><strong className="text-[var(--text-main)]">为什么不能回复或发帖？</strong><p>先查看成长中心的等级和下一阶段条件；长文还需要完成本主题的阅读要求。</p></div>
          <div><strong className="text-[var(--text-main)]">为什么提交后首页没有内容？</strong><p>内容正在冷静期或等待盲审，可在原页面查看作者可见状态，不要重复发布。</p></div>
          <div><strong className="text-[var(--text-main)]">自述提交后多久更新？</strong><p>成长中心会立即刷新为评审中，并每 5 秒同步一次；最终时间取决于评审任务完成情况。</p></div>
          <div><strong className="text-[var(--text-main)]">模型会直接决定所有内容吗？</strong><p>不会。常规治理优先由规则和匿名评审完成，模型主要用于抽检、超时兜底与讨论聚类。</p></div>
        </Section>
      </article>
    </div>
  );
}
