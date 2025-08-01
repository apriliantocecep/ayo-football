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
        title: `user_${__VU}_${__ITER} CONTOH TITLE`,
        content: `user_${__VU}_${__ITER} CONTOH CONTENT`,
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNlY2VwYXByaWxpYW50b0BnbWFpbC5jb20iLCJpc3MiOiJwb3NmaW4iLCJzdWIiOiI5NzAxMzlkYS03MzYyLTQ0M2EtYjAyOC01ODA0NWQ1ZTIyNmEiLCJleHAiOjE3NTMxMzE3MzYsIm5iZiI6MTc1MzA0NTMzNiwiaWF0IjoxNzUzMDQ1MzM2LCJqdGkiOiJlMjAxNjBjYS00MjFiLTQ4ODEtYWJlYS01Y2Y4OWM3MGI0NWQifQ._HCllX_1uFjEo-5c4T4n5bYEXoM6vXZEtFgLcR9v-oE',
        },
    };

    // Send POST request to register endpoint
    const response = http.post('http://localhost:8000/articles', payload, params);

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