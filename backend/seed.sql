-- 仅用于本地开发/课堂演示。执行前 API 应已完成全部迁移。
-- PowerShell: Get-Content backend/seed.sql -Raw | docker exec -i agora-postgres psql -U agora_user -d agora_db

INSERT INTO categories(name,slug,description,sort_order,is_active,requires_review) VALUES
  ('综合讨论','general','默认交流与综合话题讨论区',0,TRUE,FALSE),
  ('技术实践','technology','软件工程与技术实践',10,TRUE,FALSE),
  ('公共议题','public-issues','需要更审慎表达的高争议讨论区',20,TRUE,TRUE)
ON CONFLICT(slug) DO UPDATE SET name=EXCLUDED.name,description=EXCLUDED.description,sort_order=EXCLUDED.sort_order,is_active=TRUE,requires_review=EXCLUDED.requires_review;

-- admin 的密码为 admin123；其余 demo_* 账号密码均为 password。生产环境禁止执行本文件。
INSERT INTO users(username,email,password_hash,role,status) VALUES
  ('admin','2935515304@qq.com','$2a$10$DxdiQ6MMLk9/PQgBNCfiPubVQjpJgZhUtep9p1Xyr9/mkjjSG1oum','admin','active'),
  ('demo_admin','demo-admin@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','admin','active'),
  ('demo_l0','demo-l0@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l1','demo-l1@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l2','demo-l2@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l3_a','demo-l3-a@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l3_b','demo-l3-b@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active'),
  ('demo_l3_c','demo-l3-c@agora.local','$2a$10$CukE6vTZfw8gYajkRiCyk.hWDK0Z9N4DD2dlvvXF1nRS4jciYqeNu','user','active')
ON CONFLICT(username) DO UPDATE SET email=EXCLUDED.email,password_hash=EXCLUDED.password_hash,role=EXCLUDED.role,status='active',updated_at=CURRENT_TIMESTAMP;

-- 将成长旅程演示账号恢复为可重复验收的干净状态；不会影响真实用户数据。
DELETE FROM reading_sessions WHERE user_id=(SELECT id FROM users WHERE username='demo_l0');
DELETE FROM trust_logs WHERE user_id=(SELECT id FROM users WHERE username='demo_l0');
DELETE FROM blind_review_batches WHERE subject_type='user' AND author_id=(SELECT id FROM users WHERE username='demo_l0');

UPDATE user_trust_profiles tp SET
  unlock_level=CASE u.username WHEN 'demo_l0' THEN 0 WHEN 'demo_l1' THEN 1 WHEN 'demo_l2' THEN 2 ELSE 3 END,
  trust_score=CASE WHEN u.username='demo_l0' THEN 15 ELSE 20 END,
  verified_read_seconds=CASE u.username WHEN 'demo_l0' THEN 0 WHEN 'demo_l1' THEN 60 WHEN 'demo_l2' THEN 180 ELSE 300 END,
  compliant_interactions=CASE WHEN u.username IN ('admin','demo_admin','demo_l3_a','demo_l3_b','demo_l3_c') THEN 5 WHEN u.username='demo_l2' THEN 2 ELSE 0 END,
  audit_probability=0.05,
  updated_at=CURRENT_TIMESTAMP
FROM users u WHERE tp.user_id=u.id AND u.username IN ('admin','demo_admin','demo_l0','demo_l1','demo_l2','demo_l3_a','demo_l3_b','demo_l3_c');

UPDATE user_profiles p SET
  onboarding_statement='本账号用于本地课堂演示，承诺基于事实、尊重他人并说明判断依据。',
  background_tag=CASE u.username WHEN 'demo_l3_a' THEN 'engineering' WHEN 'demo_l3_b' THEN 'humanities' WHEN 'demo_l3_c' THEN 'design' ELSE 'education' END,
  onboarding_status=CASE WHEN u.username='demo_l0' THEN 'not_submitted' ELSE 'approved' END,
  updated_at=CURRENT_TIMESTAMP
FROM users u WHERE p.user_id=u.id AND u.username IN ('admin','demo_admin','demo_l0','demo_l1','demo_l2','demo_l3_a','demo_l3_b','demo_l3_c');

INSERT INTO topics(category_id,user_id,title,content,structured_content,status)
SELECT c.id,u.id,'欢迎来到 Agora-BBS','这是用于本地演示的公开主题。',jsonb_build_object('claim','这是用于本地演示的公开主题。','evidence','项目通过阅读感知、冷静期和匿名盲审改善讨论过程。','uncertainty','欢迎通过结构化回复补充不同视角。'),'published'
FROM categories c CROSS JOIN users u
WHERE c.slug='general' AND u.username='demo_admin'
  AND NOT EXISTS(SELECT 1 FROM topics WHERE title='欢迎来到 Agora-BBS');

-- 分层阅读样本：用于验证普通短文、中篇、长文以及高争议分类的差异化阅读规则。
WITH samples(category_slug,title,claim,evidence,uncertainty) AS (VALUES
  ('technology','成长测试·Go 语言：从语法到工程思维',
   'Go 的核心吸引力不是语法技巧多，而是选择少且边界清晰。变量声明、结构体、接口和错误值共同形成一套直接的工程表达：数据是什么、能力是什么、失败如何向上传递，都可以在代码中明确看见。对初学者来说，先掌握包、函数、切片、映射和结构体，再理解接口由使用方定义，通常比背诵复杂语法更有效。',
   'Go 倾向组合而非继承，并通过 gofmt 统一格式。显式错误处理虽然会增加几行代码，却使调用链上的失败路径易于定位。编写第一个服务时，可以从一个只包含 handler、service、dao 三层的小程序开始，配合 table-driven test 验证输入边界。',
   '简洁不等于所有项目都应采用同一种目录结构；团队规模、部署方式和领域复杂度仍会影响最终组织方式。'),
  ('technology','成长测试·Go 并发模型：Goroutine、Channel 与 Context',
   'Go 并发编程的重点不是尽可能多地启动 goroutine，而是明确每个并发单元的所有权、退出条件和错误传播。goroutine 很轻量，但它依然持有栈、引用和调度成本；如果调用方不知道由谁关闭它，泄漏就会从偶发问题变成长时间运行服务中的稳定故障。',
   'Channel 适合表达协作和所有权转移，互斥锁适合保护共享状态，两者并非互相替代。生产者通常负责关闭 channel，消费者通过 range 感知结束。Context 则沿调用链传递取消、截止时间和请求级数据。一个可靠的并发流程应做到：入口创建或继承 context；所有阻塞点同时监听工作信号与 ctx.Done；多个子任务使用 errgroup 汇总首个错误；关闭顺序由拥有资源的一方控制；测试中使用超时防止死锁永久挂起。限流可以使用固定大小的工作池或带缓冲 channel，但队列长度必须结合延迟目标设置，不能把内存当成无限缓冲区。',
   'Channel 的缓冲大小、工作池并发数和超时时间都依赖负载特征，应以压测和线上指标为依据，而不是照搬固定数字。'),
  ('technology','成长测试·高并发系统：从限流到可观测性的完整链路',
   '高并发系统不是把单个接口优化到极致，而是在流量进入、任务排队、数据访问、故障恢复和观测反馈之间建立一条可控制的链路。设计时应先明确服务等级目标：允许多大延迟、能够接受多少错误、峰值持续多久、哪些请求必须成功。只有目标清楚，缓存、异步化、分片和降级才有判断依据。入口层可以使用令牌桶限制瞬时速率，并按用户、租户与接口分别设定配额；应用层通过有界队列施加背压，队列满时尽快拒绝或降级，避免请求无限堆积；数据库层要控制连接池大小，让连接数与数据库实际执行能力匹配。',
   '缓存解决的是重复计算与重复读取，但必须同时设计失效策略、穿透保护和热点隔离。对于可容忍短暂旧数据的列表，可以采用较短 TTL 与主动失效结合；对于不存在的键可使用短时空值缓存；对极热键可在进程内增加一级缓存。写路径应优先保证数据库事实，再通过事件或事务后操作清理缓存。异步任务需要幂等键、有限重试、死信记录和人工重放入口，否则消息系统只会把一次失败变成多次副作用。数据库查询要用真实执行计划验证索引，分页较深时考虑基于稳定游标而非大 OFFSET。连接池不是越大越好，过多并行查询会让数据库在上下文切换和锁竞争中变慢。

过载时最重要的是保持系统可恢复。超时必须逐层递减，上游截止时间应大于下游执行预算并留出网络余量；重试只能用于可安全重放的操作，并加入指数退避和随机抖动；熔断器根据近期失败率短暂阻止请求，把资源留给可能成功的流量。日志应携带 request_id，指标至少覆盖吞吐、错误率、延迟分位数、队列深度、缓存命中率和连接池等待，链路追踪用于解释一次慢请求跨越了哪些边界。告警应围绕用户影响和错误预算，而不是每次 CPU 波动。

最后，性能结论必须可复现。压测要记录硬件、容器资源、数据规模、脚本、预热时间和持续时间；先建立基线，再一次只调整一个主要变量。达到每秒一千请求并不自动代表系统健康，如果 P95 延迟持续上升、错误被客户端重试掩盖、或者测试数据全部命中缓存，这个数字就没有解释力。真正可靠的高并发设计，是在压力超过预期时仍能有界失败、快速恢复，并让维护者知道瓶颈在哪里。',
   '具体限流阈值、缓存时长与连接池参数不能脱离部署硬件和业务分布确定；本文给出的是验证框架，而不是通用常数。'),
  ('general','成长测试·Agora 社区：为什么先阅读再回应',
   'Agora 希望把讨论中的摩擦放在有价值的位置：不是让注册和发言变麻烦，而是在回应之前给理解留出一点时间。有效阅读、结构化主题、冷静期和语境反馈分别约束不同阶段，目的不是判断观点正确与否，而是让参与者更容易说明自己在支持什么、质疑什么，以及依据来自哪里。',
   '等级并不是身份优越感，而是能力逐步开放。L0 可以阅读与收藏；完成一定有效阅读后，L1 可以回复和提供语境反馈；积累合规互动后，L2 可以发起主题；达到阅读、互动、信任和自述条件后，L3 才参与匿名盲审。长文要求滚动到底并在回复区停留，是为了降低只看标题就回应的概率。内容提交后的冷静期允许作者编辑或无痕撤回。匿名盲审只评价表达是否得体、真诚，不裁定立场是否正确。管理员可以隐藏违规内容和查看统计，但不能直接改写用户正文。',
   '任何治理机制都可能被形式化使用，因此平台仍需观察误伤率、参与成本和用户反馈，并公开调整规则的理由。'),
  ('technology','成长测试·Coding Agent 让开发者走向全栈',
   'Coding Agent 正在改变全栈开发的含义。过去的全栈往往意味着一个人记住多套框架和工具链；现在更重要的能力，是能否把目标拆成可验证的约束，让代理阅读现有系统、做出局部修改、运行测试并留下清晰证据。开发者仍然负责产品判断、架构边界和风险授权，Agent 则擅长在明确范围内完成搜索、机械修改、测试组合与文档同步。',
   '一次高质量协作通常从上下文开始：说明用户真正遇到的问题、相关文件、不可触碰的数据以及完成标准。Agent 应先阅读仓库约定和现有实现，再选择最小闭环。修改 API 时同时检查调用方、类型、错误码、迁移和测试；修改界面时同时验证加载、空状态、权限和移动端；涉及数据库时只追加迁移并考虑旧数据。工具输出不是结论，测试通过也不等于体验正确，因此仍需要浏览器中的真实流程和数据核对。

这种协作会让后端开发者更容易跨到前端，也让前端开发者能理解数据库与部署，但它没有消除专业性。安全边界、并发正确性、事务隔离、可访问性和运维恢复仍需要深入知识。Agent 可能自信地延续错误假设，也可能为了让测试通过而偏离业务目标，所以代码评审应关注为什么改、失败时怎样、是否保留用户数据。更有效的做法是把大任务分为若干可交付阶段，每阶段都有静态检查、自动化测试和可人工观察的结果。

对团队而言，Agent 还带来知识沉淀机会。过去散落在聊天和个人经验中的启动方式、故障处理、接口约定，可以在实施时同步写入仓库文档和脚本。新人不必从零猜测环境，资深成员也能把注意力放到关键决策。不过，凭据必须只保存在忽略文件或密钥系统中，自动化操作必须遵守最小权限，推送、部署、删除数据等外部动作需要明确授权。全栈化最终不是让一个人承担所有工作，而是让每个人更顺畅地穿过技术边界，同时仍对结果负责。',
   'Agent 能力、模型成本和工具生态变化很快；团队应定期复盘真实收益，并保留在没有 Agent 时仍可运行的测试与部署流程。'),
  ('public-issues','成长测试·公共议题讨论中的事实与边界',
   '公共议题往往同时包含事实判断、价值选择和个人经验。发言前先区分这三类内容，给可核查事实提供来源，给价值判断说明取舍，给个人经验标明适用范围，可以显著减少彼此误解。',
   '高争议分类即使正文较短，也执行与长文相同的阅读停留门槛；这是分类风险策略，而不是按字数判断观点价值。',
   '分类本身只能提供提醒，不能替代具体语境中的善意解释与证据核验。')
)
INSERT INTO topics(category_id,user_id,title,content,structured_content,status)
SELECT c.id,u.id,s.title,s.claim,jsonb_build_object('claim',s.claim,'evidence',s.evidence,'uncertainty',s.uncertainty),'published'
FROM samples s
JOIN categories c ON c.slug=s.category_slug
JOIN users u ON u.username='demo_admin'
WHERE NOT EXISTS(SELECT 1 FROM topics t WHERE t.title=s.title);
