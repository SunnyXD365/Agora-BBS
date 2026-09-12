import { expect, test } from '@playwright/test';

type Envelope<T> = { code: number; data: T };
type Auth = { token: string; user: { id: number; role: string; unlock_level: number; onboarding_status: string } };
type LoginResult = Partial<Auth> & { requires_email_verification?: boolean; challenge_id?: string; development_verification_code?: string };

test('registration, cooling, feedback, blind review and administration', async ({ request, page }) => {
  const api = '/api/v1';
  const login = async (username: string, password = 'password') => {
    const response = await request.post(`${api}/auth/login`, { data: { username, password } });
    expect(response.ok()).toBeTruthy();
    const data = (await response.json() as Envelope<LoginResult>).data;
    if (data.requires_email_verification) {
      expect(data.challenge_id).toBeTruthy();
      expect(data.development_verification_code).toMatch(/^\d{6}$/);
      const verified = await request.post(`${api}/auth/admin/verify-email`, { data: { challenge_id: data.challenge_id, code: data.development_verification_code } });
      expect(verified.ok()).toBeTruthy();
      return (await verified.json() as Envelope<Auth>).data;
    }
    return data as Auth;
  };
  const authHeaders = (token: string) => ({ Authorization: `Bearer ${token}` });

  await page.goto('/');
  await expect(page.getByRole('heading', { name: '第一次来到 Agora' })).toBeVisible();
  await expect(page.getByRole('heading', { name: '社区概览' })).toBeVisible();
  await expect(page.getByRole('heading', { name: '快捷入口' })).toBeVisible();
  await expect(page.getByRole('navigation', { name: '分页导航' })).toBeVisible();
  await page.getByRole('link', { name: '查看完整论坛使用手册' }).click();
  await expect(page.getByRole('heading', { name: '论坛使用手册' })).toBeVisible();

  const suffix = Date.now();
  const registered = await request.post(`${api}/auth/register`, { data: { username: `e2e_${suffix}`, email: `e2e_${suffix}@example.test`, password: 'e2e-password' } });
  expect(registered.ok()).toBeTruthy();
  const candidate = (await registered.json() as Envelope<Auth>).data;
  expect(candidate.user.unlock_level).toBe(0);

  await page.goto('/login');
  await page.evaluate(({ token }) => localStorage.setItem('token', token), candidate);
  await page.goto('/profile');
  await expect(page.getByRole('heading', { name: '如何获得权限' })).toBeVisible();
  await page.getByPlaceholder('背景标签，例如：在校学生 / 软件工程').fill('education');
  await page.getByPlaceholder('介绍你希望如何参与社区讨论…').fill('我愿意基于事实参与讨论，说明依据，尊重不同经验，并在不确定时明确表达边界。');
  await page.getByRole('button', { name: '提交自述' }).click();
  await expect(page.getByText('自述已提交，成长状态已刷新。')).toBeVisible();
  await expect(page.getByText('匿名评审进行中')).toBeVisible();

  const topicsResponse = await request.get(`${api}/topics?page_size=20`);
  const topics = (await topicsResponse.json()).data.items;
  expect(topics.length).toBeGreaterThan(0);
  const topic = topics[0];

  const l2 = await login('demo_l2');
  const cooling = await request.post(`${api}/topics`, { headers: authHeaders(l2.token), data: { category_id: topic.category_id, title: `E2E 冷静期 ${suffix}`, structured_content: { claim: '这是一个用于验证冷静期状态机的完整观点。', evidence: '测试会立即撤回，因此不会污染公开主题列表。', uncertainty: '仅用于自动化验收。' } } });
  expect(cooling.ok()).toBeTruthy();
  const coolingData = (await cooling.json()).data;
  expect(coolingData.status).toBe('cooling');
  expect((await request.delete(`${api}/topics/${coolingData.topic_id}`, { headers: authHeaders(l2.token) })).ok()).toBeTruthy();

  const l1 = await login('demo_l1');
  const feedback = await request.post(`${api}/feedbacks`, { headers: authHeaders(l1.token), data: { target_type: 'topic', target_id: topic.id, stance: 'support', tag: 'logical', reason: '核心观点与给出的依据能够相互对应，讨论边界也表达清楚。' } });
  expect(feedback.ok()).toBeTruthy();
  expect((await request.get(`${api}/feedbacks/summary?target_type=topic&target_id=${topic.id}`, { headers: authHeaders(l1.token) })).ok()).toBeTruthy();

  const reviewerNames = ['demo_admin', 'demo_l3_a', 'demo_l3_b', 'demo_l3_c'];
  let submitted = 0;
  for (const name of reviewerNames) {
    const reviewer = await login(name);
    const tasksResponse = await request.get(`${api}/reviews/tasks`, { headers: authHeaders(reviewer.token) });
    const tasks = (await tasksResponse.json()).data as Array<{ id: number; subject: { statement?: string } }>;
    const task = tasks.find((item) => item.subject.statement?.includes('我愿意基于事实参与讨论'));
    if (task && submitted < 2) {
      const result = await request.post(`${api}/reviews/tasks/${task.id}`, { headers: authHeaders(reviewer.token), data: { appropriateness: true, sincerity: true, reason: '自述表达得体且承诺具体，能够体现真诚参与社区讨论的意愿。' } });
      expect(result.ok()).toBeTruthy();
      submitted += 1;
    }
  }
  // 在全新 Seed 数据库中会完成两票；已有数据库可能随机选中额外的历史 L3 账号。
  expect(submitted).toBeGreaterThan(0);
  const profile = await request.get(`${api}/users/me`, { headers: authHeaders(candidate.token) });
  expect(['pending_review', 'approved']).toContain((await profile.json()).data.onboarding_status);

  const admin = await login('demo_admin');
  const overview = await request.get(`${api}/admin/overview`, { headers: authHeaders(admin.token) });
  expect(overview.ok()).toBeTruthy();
  expect((await overview.json()).data.users_total).toBeGreaterThan(0);

  await page.goto('/login');
  await page.evaluate(({ token, user }) => { localStorage.setItem('token', token); localStorage.setItem('e2e_user', JSON.stringify(user)); }, admin);
  await page.goto('/admin');
  await expect(page.getByRole('heading', { name: '运营总览' })).toBeVisible();
  await expect(page.getByText('邮箱二次验证已完成')).toBeVisible();
  await page.getByRole('link', { name: /用户管理/ }).first().click();
  await expect(page.getByRole('heading', { name: '用户管理' })).toBeVisible();
  await expect(page.getByRole('navigation', { name: '分页导航' })).toBeVisible();
});

test('homepage pagination requests and renders the selected page', async ({ page }) => {
  await page.route('**/api/v1/topics?**', async (route) => {
    const requestedPage = Number(new URL(route.request().url()).searchParams.get('page') || 1);
    const topic = {
      id: 9000 + requestedPage, category_id: 1, user_id: 1, author_name: 'pagination_test',
      title: `分页测试主题 · 第 ${requestedPage} 页`, content: '用于验证分页请求。',
      structured_content: { claim: '用于验证分页请求。', evidence: '', uncertainty: '' }, status: 'published',
      view_count: requestedPage, post_count: 0, like_count: 0,
      created_at: '2026-09-11T00:00:00Z', updated_at: '2026-09-11T00:00:00Z',
    };
    await route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify({ code: 0, msg: 'success', data: { items: [topic], total: 21, page: requestedPage, page_size: 10 }, request_id: 'e2e-pagination' }) });
  });
  await page.goto('/');
  await expect(page.getByRole('link', { name: '分页测试主题 · 第 1 页', exact: true })).toBeVisible();
  await page.getByRole('button', { name: '下一页' }).click();
  await expect(page.getByRole('link', { name: '分页测试主题 · 第 2 页', exact: true })).toBeVisible();
  await expect(page.getByText('第 2/3 页')).toBeVisible();
});
