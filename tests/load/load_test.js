// k6 load test — WanderPlan services
// Runs against http://localhost:8080 (api-gateway) by default.
// Usage via Docker: docker run --rm -i --network host grafana/k6 run - < tests/load/load_test.js
// Or: k6 run tests/load/load_test.js
//
// Env vars (pass with -e):
//   BASE_URL    default: http://localhost:8080
//   AUTH_TOKEN  Bearer token; generate with: go run ./tests/load/gen-test-jwt/
//   DURATION    default: 30s
//   VUS         default: 10

import http from "k6/http";
import ws from "k6/ws";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

// ── Config ───────────────────────────────────────────────────────────────────
const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";
const TOKEN = __ENV.AUTH_TOKEN || "";
const DURATION = __ENV.DURATION || "30s";
const VUS = parseInt(__ENV.VUS || "10");

export const options = {
  scenarios: {
    // Ramp up, sustain, ramp down
    ramp: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "10s", target: VUS },       // ramp up
        { duration: DURATION, target: VUS },    // sustain
        { duration: "5s", target: 0 },          // ramp down
      ],
      gracefulRampDown: "5s",
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<500", "p(99)<1000"],
    http_req_failed: ["rate<0.01"],             // <1% error rate
    checks: ["rate>0.99"],
  },
};

// ── Custom metrics ────────────────────────────────────────────────────────────
const tripCreateDuration = new Trend("trip_create_duration");
const authErrors = new Rate("auth_errors");

// ── Shared auth header ────────────────────────────────────────────────────────
const authHeaders = TOKEN
  ? { Authorization: `Bearer ${TOKEN}`, "Content-Type": "application/json" }
  : { "Content-Type": "application/json" };

// ── Helpers ───────────────────────────────────────────────────────────────────
function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function randomTrip() {
  const destinations = ["Paris", "Tokyo", "New York", "London", "Sydney", "Bali", "Rome"];
  const dest = destinations[randomInt(0, destinations.length - 1)];
  const start = new Date(Date.now() + randomInt(7, 60) * 86400000).toISOString().split("T")[0];
  const end = new Date(Date.now() + randomInt(61, 100) * 86400000).toISOString().split("T")[0];
  return JSON.stringify({
    title: `${dest} Adventure ${randomInt(1000, 9999)}`,
    destination: dest,
    start_date: start,
    end_date: end,
    status: "planning",
    visibility: "private",
    budget_total: randomInt(500, 5000),
    currency: "USD",
  });
}

// ── Scenario: health checks (all 7 services direct) ──────────────────────────
export function healthChecks() {
  const services = [
    { name: "api-gateway", url: `${BASE_URL}/healthz` },
    { name: "auth-service", url: "http://localhost:8081/healthz" },
    { name: "trip-service", url: "http://localhost:8082/healthz" },
    { name: "user-service", url: "http://localhost:8083/healthz" },
    { name: "collab-service", url: "http://localhost:8084/healthz" },
    { name: "notif-service", url: "http://localhost:8085/healthz" },
    { name: "search-service", url: "http://localhost:8086/healthz" },
  ];
  for (const svc of services) {
    const res = http.get(svc.url, { tags: { service: svc.name } });
    check(res, {
      [`${svc.name} status 200`]: (r) => r.status === 200,
    });
  }
  sleep(0.5);
}

// ── Main VU function ──────────────────────────────────────────────────────────
export default function () {
  const scenario = randomInt(0, 4);

  if (scenario === 0) {
    // Health check via gateway
    const res = http.get(`${BASE_URL}/healthz`);
    check(res, { "healthz 200": (r) => r.status === 200 });
    sleep(0.1);
    return;
  }

  if (!TOKEN) {
    // Without a JWT, only test health + metrics
    const res = http.get(`${BASE_URL}/healthz`);
    check(res, { "healthz 200": (r) => r.status === 200 });
    sleep(0.5);
    return;
  }

  if (scenario === 1) {
    // List trips
    const res = http.get(`${BASE_URL}/api/v1/trips`, { headers: authHeaders });
    const ok = check(res, {
      "list trips 200": (r) => r.status === 200,
      "list trips has data": (r) => {
        try { return JSON.parse(r.body).data !== undefined; } catch { return false; }
      },
    });
    if (!ok) authErrors.add(1);
    sleep(0.2);
    return;
  }

  if (scenario === 2) {
    // Create trip
    const start = Date.now();
    const res = http.post(`${BASE_URL}/api/v1/trips`, randomTrip(), { headers: authHeaders });
    tripCreateDuration.add(Date.now() - start);
    const ok = check(res, {
      "create trip 201": (r) => r.status === 201,
    });
    if (!ok) authErrors.add(1);

    // If trip was created, fetch it
    if (res.status === 201) {
      try {
        const body = JSON.parse(res.body);
        const tripID = body.data && (body.data.id || body.data.trip_id);
        if (tripID) {
          const getRes = http.get(`${BASE_URL}/api/v1/trips/${tripID}`, {
            headers: authHeaders,
          });
          check(getRes, { "get trip 200": (r) => r.status === 200 });
        }
      } catch (_) {}
    }
    sleep(0.3);
    return;
  }

  if (scenario === 3) {
    // Get user profile
    const res = http.get(`${BASE_URL}/api/v1/users/me`, { headers: authHeaders });
    check(res, { "get me 200|404": (r) => r.status === 200 || r.status === 404 });
    sleep(0.2);
    return;
  }

  if (scenario === 4) {
    // Search trips
    const queries = ["Paris", "Tokyo", "beach", "mountain", "Europe"];
    const q = queries[randomInt(0, queries.length - 1)];
    const res = http.get(`${BASE_URL}/api/v1/search/trips?q=${encodeURIComponent(q)}`, {
      headers: authHeaders,
    });
    check(res, { "search trips 200|404": (r) => r.status === 200 || r.status === 404 });
    sleep(0.2);
    return;
  }

  sleep(0.1);
}
