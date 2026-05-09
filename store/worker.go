package store

import (
	"context"
	"database/sql"

	"github.com/drumilbhati/disscheduler/model"
)

func (s *Store) drainJob() error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var job model.Job

	err = tx.QueryRow(
		`SELECT id, type, payload, priority, run_at, idempotency_key, status, attempts, max_attempts, created_at, updated_at FROM job
		WHERE status = 'queued' AND run_at <= NOW()
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
		return nil
	} else if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE job
		SET status = 'running', updated_at = NOW() WHERE id =$1`,
		job.ID,
	)

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	go s.processJob(job)

	return nil
}

func (s *Store) processJob(job model.Job) {
	var err error
	switch job.Type {
	case "email":
		// Process email jobi
	case "report":
		// Process report job
	case "webhook":
		// Process webhook job
	case "cleanup":
		// Process cleanup job
	default:
		// Handle unknown job type
	}

	if err != nil {
		// mark the job as failed and log the error
	}

	// mark the job as completed
}
