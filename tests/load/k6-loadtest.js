import http from 'k6/http';
import { check, group, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },  // Ramp-up to 50 VUs
    { duration: '1m', target: 200 },  // Stay at 200 VUs
    { duration: '30s', target: 500 },  // Peak spike to 500 VUs
    { duration: '30s', target: 0 },    // Ramp-down to 0 VUs
  ],
  thresholds: {
    http_req_duration: ['p(95)<300', 'p(99)<800'], // 95% of requests under 300ms
    http_req_failed: ['rate<0.01'],                 // Error rate < 1%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  group('1. Health and Gateway Liveness', () => {
    const res = http.get(`${BASE_URL}/healthz`);
    check(res, {
      'status is 200': (r) => r.status === 200,
      'healthy response': (r) => r.json('status') === 'healthy',
    });
  });

  group('2. Public Feed Read Throughput', () => {
    const res = http.get(`${BASE_URL}/api/v1/feed?limit=20`);
    check(res, {
      'status is 200 or 401 (auth enforced)': (r) => r.status === 200 || r.status === 401,
    });
  });

  group('3. Ad Auction Delivery Latency', () => {
    const payload = JSON.stringify({
      user_id: 'vu_' + __VU,
      age: 26,
      gender: 'MALE',
      location: 'KR',
      interests: ['gaming', 'tech'],
    });

    const params = {
      headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post(`${BASE_URL}/api/v1/ads/auction`, payload, params);
    check(res, {
      'auction handled in < 50ms': (r) => r.timings.duration < 50,
    });
  });

  sleep(0.5);
}
