import { expect, test, type APIRequestContext } from '@playwright/test';

type Envelope<T> = { code: number; data: T };
type User = { id: number; unlock_level: number; verified_read_seconds: number; onboarding_status: string };
type Auth = { token: string; user: User };
type Topic = { id: number; title: string };
type Session = {
  id: string;
  reading_seconds: number;
  reply_dwell_seconds: number;
  bottom_reached: boolean;
  requires_reply_dwell: boolean;
  eligible: boolean;
  completed: boolean;
};

const enabled = process.env.RUN_GROWTH_JOURNEY === '1';

async function login(request: APIRequestContext, username: string): Promise<Auth> {
  const response = await request.post('/api/v1/auth/login', { data: { username, password: 'password' } });
  expect(response.ok()).toBeTruthy();
  return (await response.json() as Envelope<Auth>).data;
}

async function readFor(
  request: APIRequestContext,
  token: string,
  topicID: number,
  heartbeatCount: number,
  focusAtEnd = false,
): Promise<Session> {
  const headers = { Authorization: `Bearer ${token}` };
  const started = await request.post('/api/v1/reading-sessions', { headers, data: { topic_id: topicID } });
  expect(started.ok()).toBeTruthy();
  let session = (await started.json() as Envelope<Session>).data;
  for (let index = 0; index < heartbeatCount; index += 1) {
    await new Promise((resolve) => setTimeout(resolve, 5_100));
    const progress = Math.min(100, Math.round(((index + 1) / heartbeatCount) * 100));
    const replyFocused = focusAtEnd && index >= heartbeatCount - 3;
    const heartbeat = await request.patch(`/api/v1/reading-sessions/${session.id}/heartbeat`, {
      headers,
      data: { progress, reply_focused: replyFocused },
    });
    expect(heartbeat.ok()).toBeTruthy();
    session = (await heartbeat.json() as Envelope<Session>).data;
  }
  const completed = await request.post(`/api/v1/reading-sessions/${session.id}/complete`, {
    headers,
    data: { progress: 100, reply_focused: focusAtEnd },
  });
  expect(completed.ok()).toBeTruthy();
  return (await completed.json() as Envelope<Session>).data;
}

test('demo_l0 follows the real development growth path to L3', async ({ request, page }) => {
  test.skip(!enabled, 'Set RUN_GROWTH_JOURNEY=1 and reset demo data with backend/seed.sql to run this five-minute journey.');
  test.setTimeout(8 * 60_000);

  const candidate = await login(request, 'demo_l0');
  const headers = { Authorization: `Bearer ${candidate.token}` };
  expect(candidate.user.unlock_level).toBe(0);

  const topicsResponse = await request.get('/api/v1/topics?page=1&page_size=50');
  expect(topicsResponse.ok()).toBeTruthy();
  const topics = (await topicsResponse.json() as Envelope<{ items: Topic[] }>).data.items;
  const byTitle = (fragment: string) => {
    const topic = topics.find((item) => item.title.includes(fragment));
    expect(topic, `missing seeded topic: ${fragment}`).toBeTruthy();
    return topic!;
  };
  const shortTopic = byTitle('Go 语言：从语法到工程思维');
  const mediumTopic = byTitle('Go 并发模型');
  const longTopic = byTitle('高并发系统');
  const sensitiveTopic = byTitle('公共议题讨论中的事实与边界');

  const denied = await request.post(`/api/v1/topics/${shortTopic.id}/posts`, {
    headers,
    data: { content: 'L0 阶段不应允许直接回复。', post_type: 'experience' },
  });
  expect(denied.status()).toBe(403);

  const statementMarker = `成长旅程 ${Date.now()}`;
  const onboarding = await request.put('/api/v1/users/me/onboarding', {
    headers,
    data: {
      background_tag: 'quality-engineering',
      statement: `${statementMarker}：我会先完整阅读，再区分事实、经验与判断，给出具体依据并尊重不同背景的参与者。`,
    },
  });
  expect(onboarding.ok()).toBeTruthy();

  let reviewVotes = 0;
  for (let attempt = 0; attempt < 10 && reviewVotes < 2; attempt += 1) {
    for (const reviewerName of ['demo_l3_a', 'demo_l3_b', 'demo_l3_c']) {
      const reviewer = await login(request, reviewerName);
      const reviewerHeaders = { Authorization: `Bearer ${reviewer.token}` };
      const tasksResponse = await request.get('/api/v1/reviews/tasks', { headers: reviewerHeaders });
      expect(tasksResponse.ok()).toBeTruthy();
      const tasks = (await tasksResponse.json() as Envelope<Array<{ id: number; subject: { statement?: string } }>>).data;
      const task = tasks.find((item) => item.subject.statement?.includes(statementMarker));
      if (!task) continue;
      const submitted = await request.post(`/api/v1/reviews/tasks/${task.id}`, {
        headers: reviewerHeaders,
        data: { appropriateness: true, sincerity: true, reason: '自述具体说明了阅读、证据与尊重边界，表达得体且态度真诚。' },
      });
      if (submitted.ok()) reviewVotes += 1;
      if (reviewVotes >= 2) break;
    }
    if (reviewVotes < 2) await new Promise((resolve) => setTimeout(resolve, 1_000));
  }
  expect(reviewVotes).toBeGreaterThanOrEqual(2);

  const shortReading = await readFor(request, candidate.token, shortTopic.id, 13);
  expect(shortReading.requires_reply_dwell).toBeFalsy();
  expect(shortReading.completed).toBeTruthy();
  expect(shortReading.eligible).toBeTruthy();

  let profileResponse = await request.get('/api/v1/users/me', { headers });
  let profile = (await profileResponse.json() as Envelope<User>).data;
  expect(profile.unlock_level).toBe(1);
  expect(profile.verified_read_seconds).toBeGreaterThanOrEqual(60);

  for (let index = 0; index < 5; index += 1) {
    const post = await request.post(`/api/v1/topics/${shortTopic.id}/posts`, {
      headers,
      data: {
        content: `成长旅程合规回复 ${index + 1}：这篇文章把语言特性与工程边界联系起来，我会在后续实践中用测试验证理解。`,
        post_type: index % 2 === 0 ? 'evidence' : 'experience',
      },
    });
    expect(post.ok()).toBeTruthy();
    expect((await post.json() as Envelope<{ status: string }>).data.status).toBe('cooling');
  }

  const mediumReading = await readFor(request, candidate.token, mediumTopic.id, 25);
  expect(mediumReading.requires_reply_dwell).toBeFalsy();
  expect(mediumReading.completed).toBeTruthy();
  profileResponse = await request.get('/api/v1/users/me', { headers });
  profile = (await profileResponse.json() as Envelope<User>).data;
  expect(profile.unlock_level).toBe(2);
  expect(profile.verified_read_seconds).toBeGreaterThanOrEqual(180);
  expect(profile.onboarding_status).toBe('approved');

  const longReading = await readFor(request, candidate.token, longTopic.id, 24, true);
  expect(longReading.requires_reply_dwell).toBeTruthy();
  expect(longReading.reply_dwell_seconds).toBeGreaterThanOrEqual(10);
  expect(longReading.completed).toBeTruthy();
  profileResponse = await request.get('/api/v1/users/me', { headers });
  profile = (await profileResponse.json() as Envelope<User>).data;
  expect(profile.verified_read_seconds).toBeGreaterThanOrEqual(300);
  expect(profile.unlock_level).toBe(3);

  const sensitiveStart = await request.post('/api/v1/reading-sessions', { headers, data: { topic_id: sensitiveTopic.id } });
  let sensitive = (await sensitiveStart.json() as Envelope<Session>).data;
  expect(sensitive.requires_reply_dwell).toBeTruthy();
  await new Promise((resolve) => setTimeout(resolve, 5_100));
  const earlyComplete = await request.post(`/api/v1/reading-sessions/${sensitive.id}/complete`, {
    headers,
    data: { progress: 100, reply_focused: false },
  });
  sensitive = (await earlyComplete.json() as Envelope<Session>).data;
  expect(sensitive.completed).toBeFalsy();
  expect(sensitive.eligible).toBeFalsy();

  await page.goto('/login');
  await page.evaluate(({ token }) => localStorage.setItem('token', token), candidate);
  await page.goto('/profile');
  await expect(page.getByRole('heading', { name: /L3 · 社区评审者/ })).toBeVisible();
  await expect(page.getByText(/已验证阅读 [1-9]\d* 分钟/)).toBeVisible();
});
