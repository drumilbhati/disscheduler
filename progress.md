# Disscheduler Progress Tracker

Last updated: 2026-05-18

This file tracks implementation progress against the problem statement in `README.md`.

## Functional Requirements (FR)

| Requirement | Status | Current implementation |
| --- | --- | --- |
| FR-1 Job Ingestion | **Partial** | `POST /job` exists and persists jobs with generated UUID and initial `queued` state. Input validation is still minimal. |
| FR-2 Durable Persistence | **Partial** | Jobs are written to PostgreSQL before response. Persistence durability depends on DB availability and current schema setup. |
| FR-3 Claiming and Execution | **Partial** | Workers claim jobs using transaction + `FOR UPDATE SKIP LOCKED`, mark `running`, then process. Type-specific handler dispatch is not implemented yet. |
| FR-4 Retry and Backoff | **Partial** | Failed processing increments attempts and requeues with jittered delay until `max_attempts`; then marks `failed`. Policy is basic and not configurable yet. |
| FR-5 Dead-Letter Queue | **Not started** | No `dead_lettered` state, failure reason storage, or requeue endpoint yet. |
| FR-6 Scheduling | **Partial** | Delayed execution via `run_at` works for queued jobs. Recurring cron-like schedules are not implemented. |
| FR-7 Status and Query APIs | **Partial** | `GET /job` exists. `GET /job/{id}`, filtering, pagination, and richer querying are missing. |
| FR-8 Graceful Shutdown and Recovery | **Partial** | Graceful shutdown added: signal-aware context, cancellable workers, and HTTP server shutdown timeout. Stale-running recovery/lease reclaim logic is not implemented yet. |

## Non-Functional Requirements (NFR)

| Requirement | Status | Current implementation |
| --- | --- | --- |
| NFR-1 Reliability | **Partial** | Core state transitions exist (`queued` → `running` → `succeeded/failed`). More robust failure semantics and DLQ handling are pending. |
| NFR-2 Scalability | **Partial** | Multi-worker model exists and supports safe concurrent claiming through row locking. Horizontal scaling validation is still pending. |
| NFR-3 Performance | **Not measured** | Worker count is fixed in code; no throughput/latency benchmarks yet. |
| NFR-4 Observability | **Not started** | Basic request logging only. No metrics, health/readiness endpoints, or job-correlated structured logs. |
| NFR-5 Maintainability | **Partial** | Basic package boundaries (`controller`, `store`, `db`, `model`) are present. Config-driven runtime controls are still limited. |

## Milestones Completed

1. PostgreSQL schema and job persistence implemented.
2. Job create/list API endpoints implemented.
3. Concurrent worker loop with transactional job claiming implemented.
4. Context-based graceful shutdown for workers and HTTP server implemented.
5. Context-aware DB operations added in worker execution path.

## Next Focus

1. Implement DLQ (`dead_lettered` state + failure reason + requeue path).
2. Add job detail/filter APIs (`GET /job/{id}`, filtering/pagination).
3. Add health/readiness endpoints and baseline metrics.
4. Add stale-running recovery/lease logic.
5. Add integration tests for claim/retry/shutdown behavior.
