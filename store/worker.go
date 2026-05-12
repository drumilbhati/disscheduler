package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/drumilbhati/disscheduler/model"
)

func (s *Store) StartWorkers(n int) {
	for i := range n {
		go s.worker(i)
	}
}

func (s *Store) worker(id int) {
	for {
		job, err := s.drainJob()
		if err != nil {
			log.Printf("worker=%d drain error: %v\n", id, err)
			time.Sleep(time.Second)
			continue
		}
		if job == nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		s.processJob(job)
	}
}

func (s *Store) drainJob() (*model.Job, error) {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var job model.Job

	err = tx.QueryRow(
		`SELECT id, type, payload, priority, run_at, idempotency_key, status, attempts, max_attempts, created_at, updated_at FROM job
		WHERE status = 'queued'
		AND run_at IS NULL OR run_at <= NOW()
		ORDER BY priority DESC, run_at ASC
		FOR UPDATE
		SKIP LOCKED
		LIMIT 1`,
	).Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Priority,
		&job.RunAt,
		&job.IdempotencyKey,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		`UPDATE job
		SET status = 'running', updated_at = NOW() WHERE id =$1`,
		job.ID,
	)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &job, nil
}
