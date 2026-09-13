import { expect, test } from '@playwright/test';

type Envelope<T> = { code: number; data: T };
type Auth = { token: string; user: { verified_read_seconds: number } };
type Session = { id: string; completed: boolean; eligible: boolean };

test('ordinary reading settles on page leave, displays seconds and remains idempotent', async ({ request, page }) => {
  const suffix = Date.now();
  const registered = await request.post('/api/v1/auth/register', {
    data: { username: `reader_${suffix}`, email: `reader_${suffix}@example.test`, password: 'reader-password' },
  });
  expect(registered.ok()).toBeTruthy();
  const auth = (await registered.json() as Envelope<Auth>).data;
  const headers = { Authorization: `Bearer ${auth.token}` };

  const topicsResponse = await request.get('/api/v1/topics?page=1&page_size=50');
  const topics = (await topicsResponse.json() as Envelope<{ items: Array<{ id: number; title: string }> }>).data.items;
  const ordinary = topics.find((topic) => topic.title === '欢迎来到 Agora-BBS');
  expect(ordinary).toBeTruthy();

  await page.goto('/login');
  await page.evaluate(({ token }) => localStorage.setItem('token', token), auth);
  const startedResponse = page.waitForResponse((response) => response.url().includes('/api/v1/reading-sessions') && response.request().method() === 'POST');
  await page.goto(`/topics/${ordinary!.id}`);
  const started = (await (await startedResponse).json() as Envelope<Session>).data;
  await page.evaluate(() => window.scrollTo(0, document.documentElement.scrollHeight));
  await page.waitForTimeout(2_100);

  // A full navigation triggers pagehide. The keepalive completion request must
  // carry the last progress even though the regular 5-second heartbeat did not run.
  await page.goto('/profile');
  await expect.poll(async () => {
    const profile = await request.get('/api/v1/users/me', { headers });
    return (await profile.json() as Envelope<{ verified_read_seconds: number }>).data.verified_read_seconds;
  }).toBeGreaterThan(0);

  const beforeDuplicate = await request.get('/api/v1/users/me', { headers });
  const beforeSeconds = (await beforeDuplicate.json() as Envelope<{ verified_read_seconds: number }>).data.verified_read_seconds;
  const duplicate = await request.post(`/api/v1/reading-sessions/${started.id}/complete`, {
    headers,
    data: { progress: 100, reply_focused: false },
  });
  expect(duplicate.ok()).toBeTruthy();
  const duplicateSession = (await duplicate.json() as Envelope<Session>).data;
  expect(duplicateSession.completed).toBeTruthy();
  expect(duplicateSession.eligible).toBeTruthy();
  const afterDuplicate = await request.get('/api/v1/users/me', { headers });
  expect((await afterDuplicate.json() as Envelope<{ verified_read_seconds: number }>).data.verified_read_seconds).toBe(beforeSeconds);

  await page.reload();
  await expect(page.getByText(/已验证阅读 [1-9]\d* 秒/)).toBeVisible();
});

test('forum guide contributes verified reading time after reaching the bottom', async ({ request, page }) => {
  const suffix = Date.now();
  const registered = await request.post('/api/v1/auth/register', {
    data: { username: `guide_${suffix}`, email: `guide_${suffix}@example.test`, password: 'guide-password' },
  });
  expect(registered.ok()).toBeTruthy();
  const auth = (await registered.json() as Envelope<Auth>).data;
  const headers = { Authorization: `Bearer ${auth.token}` };

  await page.goto('/login');
  await page.evaluate(({ token }) => localStorage.setItem('token', token), auth);
  const startedResponse = page.waitForResponse(async (response) => response.url().includes('/api/v1/reading-sessions')
    && response.request().method() === 'POST'
    && (response.request().postDataJSON() as { resource?: string }).resource === 'forum-guide');
  await page.goto('/guide');
  const started = (await (await startedResponse).json() as Envelope<Session & { resource_type: string; resource_key: string }>).data;
  expect(started).toMatchObject({ resource_type: 'guide', resource_key: 'forum-guide' });
  await expect(page.getByText(/手册阅读进度/)).toBeVisible();

  await page.evaluate(() => window.scrollTo(0, document.documentElement.scrollHeight));
  await expect(page.getByText(/本次手册阅读已完成，\d+ 秒已计入成长中心/)).toBeVisible({ timeout: 12_000 });
  await expect.poll(async () => {
    const profile = await request.get('/api/v1/users/me', { headers });
    return (await profile.json() as Envelope<{ verified_read_seconds: number }>).data.verified_read_seconds;
  }).toBeGreaterThan(0);

  const invalid = await request.post('/api/v1/reading-sessions', { headers, data: { resource: 'unknown-guide' } });
  expect(invalid.status()).toBe(400);
});
