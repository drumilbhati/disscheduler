package store

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/drumilbhati/disscheduler/model"
	"github.com/gofrs/uuid"
)

func (s *Store) CreateJob(job *model.Job) error {
	jobID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	job.ID = jobID
	job.Status = model.StatusQueued
	job.Attempts = 0
	job.UpdatedAt = now

	err = s.db.QueryRow(
		`INSERT INTO job(id, type, payload, priority, run_at, idempotency_key, status, attempts, max_attempts, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`,
		job.ID,
		job.Type, job.Payload, job.Priority, job.RunAt, job.IdempotencyKey, job.Status, job.Attempts, job.MaxAttempts, job.UpdatedAt,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetAllJobs() ([]model.Job, error) {
	var jobs []model.Job

	rows, err := s.db.Query(
		`SELECT id, type, payload, priority, run_at, idempotency_key, status, attempts, max_attempts, created_at, updated_at FROM job`,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var job model.Job
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *Store) processJob(job *model.Job, ctx context.Context) {
	// Simulate job processing
	select {
	case <-ctx.Done():
		return
	case <-time.After(2 * time.Second):
	}

	_, err := s.db.ExecContext(
		ctx,
		`UPDATE job SET status = $1, updated_at = NOW() WHERE id = $2`,
		model.StatusSucceeded,
		job.ID,
	)
	if err == nil {
		return
	}

	job.Attempts++
	if job.Attempts >= job.MaxAttempts {
		if _, updateErr := s.db.ExecContext(
			ctx,
			`UPDATE job SET status = $1, attempts = $2, updated_at = NOW() WHERE id = $3`,
			model.StatusFailed,
			job.Attempts,
			job.ID,
		); updateErr != nil {
			log.Printf("error updating failed status for job %s: %v", job.ID, updateErr)
		}
		return
	}

	maxJitterSeconds := 1 << job.Attempts
	jitterSeconds := rand.Intn(maxJitterSeconds)
	retryAt := time.Now().UTC().Add(time.Second * time.Duration(jitterSeconds))
	job.RunAt = &retryAt
	if _, updateErr := s.db.ExecContext(
		ctx,
		`UPDATE job SET status = $1, attempts = $2, run_at = $3, updated_at = NOW() WHERE id = $4`,
		model.StatusQueued,
		job.Attempts,
		job.RunAt,
		job.ID,
	); updateErr != nil {
		log.Printf("error requeueing job %s: %v", job.ID, updateErr)
	}
}
