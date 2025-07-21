import http from 'k6/http';
import { check } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

export let options = {
    stages: [
        { duration: '10s', target: 50 },   // Ramp up to 50 VUs in 10s
        { duration: '50s', target: 100 },  // Ramp up to 100 VUs in 50s
        { duration: '10s', target: 0 },    // Ramp down to 0 VUs in 10s
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95th percentile should be below 200ms
        errors: ['rate<0.1'], // Error rate should be below 10%
    },
};

export default function () {
    // Sample data untuk register - sesuaikan dengan API requirements
    const payload = JSON.stringify({
        name: `user_${__VU}_${__ITER}_queue`, // Unique username per VU and iteration
        email: `user_${__VU}_${__ITER}_queue_${Math.random()}@example.com`,
        password: 'testPassword123',
        // Tambahkan field lain sesuai kebutuhan API Anda
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Send POST request to register endpoint
    const response = http.post('http://localhost:8000/auth/register/queue', payload, params);

    // Checks
    const result = check(response, {
        'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
        'response time < 200ms': (r) => r.timings.duration < 200,
        'response has body': (r) => r.body.length > 0,
    });

    // Track errors
    errorRate.add(!result);

    // Optional: Add sleep between requests if needed
    // sleep(1);
}