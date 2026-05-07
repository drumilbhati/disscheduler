# Distributed Job Queue & Scheduler in Go

## 1. Problem Statement

Modern backend systems need to process work asynchronously: sending emails, generating reports, syncing data, processing media, and running recurring maintenance tasks.  
Synchronous handling of such work increases API latency, reduces reliability under load, and creates poor failure isolation.

The goal of this project is to design and implement a **production-grade distributed job queue and scheduler** in **Golang** that enables reliable asynchronous execution of one-time and recurring jobs across multiple worker instances.

The system should support the full job lifecycle from ingestion to completion while preserving correctness under failures, restarts, and horizontal scaling.

---

## 2. Background and Motivation

As systems grow, asynchronous processing becomes essential for:

- Decoupling user-facing request latency from long-running tasks
- Retrying transient failures without blocking clients
- Controlling throughput using queues and worker pools
- Running delayed and periodic jobs reliably
- Operating safely in distributed, multi-instance environments

This project models real-world backend infrastructure and directly builds skills in system design, distributed coordination, data consistency, and production operations.

---

## 3. Project Goal

Build a service that allows clients to submit jobs and schedules, then executes them reliably with observability and operational control.

At a high level, the system includes:

1. **Producer API**: accepts job creation and scheduling requests
2. **Persistent queue store**: durable job and schedule state
3. **Dispatcher/claimer**: selects eligible jobs for execution
4. **Workers**: process jobs concurrently with retry policies
5. **Scheduler**: materializes delayed and recurring runs
6. **Control/observability layer**: metrics, logs, health, and admin endpoints

---

## 4. Formal Scope

### In Scope

- Create and persist asynchronous jobs with payloads and metadata
- Execute jobs using a configurable worker pool
- Track job states end-to-end (`queued`, `running`, `succeeded`, `failed`, `dead_lettered`)
- Retry failed jobs with configurable exponential backoff
- Move exhausted jobs to a dead-letter queue (DLQ)
- Support delayed jobs (`run_at`)
- Support recurring jobs (cron-like schedules)
- Support idempotency keys / deduplication for safe retries
- Run multiple worker instances concurrently without unsafe double-processing
- Expose operational visibility (metrics/logs/health/status APIs)

### Out of Scope (initial version)

- Complex workflow orchestration / DAGs
- Exactly-once delivery guarantees
- Multi-region replication
- Arbitrary script execution sandboxing
- Rich web UI dashboard (CLI/basic APIs are acceptable)

---

## 5. Functional Requirements

### FR-1: Job Ingestion
- Clients can submit a job with:
  - `type`
  - `payload` (JSON)
  - `priority` (optional)
  - `run_at` (optional delayed execution)
  - `idempotency_key` (optional but recommended)
- The API returns a unique job ID and initial state.

### FR-2: Durable Persistence
- Accepted jobs must be persisted before acknowledgment.
- Job data survives process crash/restart.

### FR-3: Job Claiming and Execution
- Eligible queued jobs are claimed by workers safely.
- Claimed jobs transition to `running`.
- Worker executes handler based on job type.

### FR-4: Retry and Backoff
- Failed jobs are retried until `max_attempts` is reached.
- Retry delay follows configurable backoff strategy.

### FR-5: Dead-Letter Queue
- Jobs that exceed retry limits are moved to DLQ state with failure reason.
- Operators can inspect and optionally requeue DLQ jobs.

### FR-6: Scheduling
- Delayed jobs become executable at/after `run_at`.
- Recurring schedules generate executable jobs according to cron expression.

### FR-7: Status and Query APIs
- Fetch job by ID and view current state, attempts, timestamps, and error details.
- List/filter jobs by state, type, and time range.

### FR-8: Graceful Shutdown and Recovery
- On shutdown, in-flight jobs are handled predictably.
- On restart, stale running jobs are recovered/reclaimed according to lease policy.

---

## 6. Non-Functional Requirements

### NFR-1: Reliability
- No silently dropped accepted jobs.
- Deterministic state transitions and failure reporting.

### NFR-2: Scalability
- Support horizontal scaling of workers with safe concurrent claiming.

### NFR-3: Performance
- Bounded API latency for job submission under moderate load.
- Configurable throughput via worker concurrency and queue policies.

### NFR-4: Observability
- Structured logs with correlation/job IDs.
- Metrics for queue depth, job latency, success/failure/retry rates.
- Health/readiness endpoints.

### NFR-5: Maintainability
- Clear module boundaries (API, queue store, scheduler, workers, handlers).
- Config-driven behavior for retries, timeouts, and pool sizes.

---

## 7. Constraints and Assumptions

### Constraints
- Implementation language: **Go**
- Durable storage: **PostgreSQL** (or equivalent transactional DB)
- Containerized local environment (e.g., Docker Compose)

### Assumptions
- Delivery semantics target **at-least-once**, not exactly-once.
- Job handlers are designed to be idempotent or guarded by idempotency keys.
- Minor schedule jitter is acceptable within a defined tolerance window.

---

## 8. Failure Model to Handle

The system must remain correct under:

- Worker crash during job execution
- API process crash after request receipt but before response
- Database connection interruptions/transient errors
- Duplicate submission attempts from clients
- Multiple workers racing to claim the same job
- Scheduler restarts causing potential duplicate schedule evaluation

Expected behavior: failures are surfaced, retried or dead-lettered according to policy, and never silently ignored.

---

## 9. Success Criteria (Acceptance)

The project is successful if:

1. Accepted jobs persist durably and are eventually executed or dead-lettered.
2. Retry + backoff + DLQ behavior is deterministic and configurable.
3. Delayed and recurring jobs are executed close to scheduled time.
4. Multi-worker deployment processes jobs safely under load.
5. Operators can inspect health, queue depth, and job outcomes through logs/metrics/APIs.

---

## 10. Suggested Deliverables

- Source code for API, scheduler, worker, and persistence layer
- DB schema + migrations
- Config files/environment template
- Docker Compose for local distributed setup
- Seed/demo handlers (e.g., email/report/sync simulation)
- Basic load test script and reliability test scenarios
- Documentation for architecture, setup, and operational runbook

---

## 11. Why This Project Matters for System Design Growth

This project provides hands-on exposure to real design trade-offs:

- Consistency vs throughput
- Retry safety vs duplicate side effects
- Scheduling precision vs operational simplicity
- Lease/locking strategies in distributed processing
- Observability-driven debugging in asynchronous systems

Completing this project end-to-end builds practical skills directly transferable to backend and distributed systems roles.

