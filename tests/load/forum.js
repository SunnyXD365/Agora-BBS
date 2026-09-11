import http from 'k6/http';
import { check } from 'k6';

const full = __ENV.PROFILE === 'full';
// One iteration makes three HTTP requests. 334 iterations/s is therefore
// approximately the 1000 HTTP requests/s target documented for this project.
const rate = Number(__ENV.RATE || (full ? 334 : 20));
const warmup = full ? '2m' : '5s';
const duration = full ? '5m' : '20s';

export const options = {
  scenarios: {
    warmup: {
      executor: 'ramping-arrival-rate',
      startRate: Math.max(1, Math.floor(rate / 10)),
      timeUnit: '1s',
      preAllocatedVUs: full ? 200 : 10,
      maxVUs: full ? 2000 : 50,
      stages: [{ target: rate, duration: warmup }],
    },
    sustained: {
      executor: 'constant-arrival-rate',
      startTime: warmup,
      rate,
      timeUnit: '1s',
      duration,
      preAllocatedVUs: full ? 500 : 20,
      maxVUs: full ? 3000 : 100,
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<500'],
  },
};

const baseURL = __ENV.BASE_URL || 'http://agora-nginx';

export default function () {
  const responses = http.batch([
    ['GET', `${baseURL}/api/v1/topics?page=1&page_size=20`],
    ['GET', `${baseURL}/api/v1/categories`],
    ['GET', `${baseURL}/api/v1/governance/policy`],
  ]);
  for (const response of responses) {
    check(response, { 'status is 200': (item) => item.status === 200 });
  }
}

export function handleSummary(data) {
  return {
    '/scripts/results/latest-summary.json': JSON.stringify(data, null, 2),
    stdout: `\nAgora-BBS k6: requests=${data.metrics.http_reqs?.values.count || 0}, failed=${data.metrics.http_req_failed?.values.rate || 0}, p95_ms=${data.metrics.http_req_duration?.values['p(95)'] || 0}\n`,
  };
}
