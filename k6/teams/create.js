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
        "name": "Persib",
        "logo": "https://upload.wikimedia.org/wikipedia/id/thumb/0/0d/Logo_Persib_Bandung.png/250px-Logo_Persib_Bandung.png",
        "founded_at": 1933,
        "address": "Jl. sulanjana",
        "city": "Bandung"
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImNlY2VwYXByaWxpYW50b0BnbWFpbC5jb20iLCJpc3MiOiJheW9mb290YmFsbCIsInN1YiI6IjQ0NTJlYjcxLWIzY2ItNDUwZS1iZDRkLThjNzYyYmQ0MWI3MiIsImV4cCI6MTc1NDA2NzA4MSwibmJmIjoxNzUzOTgwNjgxLCJpYXQiOjE3NTM5ODA2ODEsImp0aSI6ImRkNTViZjQyLWEzOGYtNDg0Zi1iZDkzLTJkNjg4NzliN2NkMCJ9.b1WUPLizVbGbmiI6zkt_g478d2PSavQEYMvCxakPmL0',
        },
    };

    // Send POST request to register endpoint
    const response = http.post('http://localhost:8000/teams', payload, params);

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