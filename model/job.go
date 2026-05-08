package model

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
)

type Job struct {
	ID             uuid.UUID       `json:"id"`
	Type           string          `json:"type"`
	Payload        json.RawMessage `json:"payload"`
	Priority       int             `json:"priority"`
	RunAt          *time.Time      `json:"run_at,omitempty"`
	IdempotencyKey *uuid.UUID      `json:"idempotency_key,omitempty"`
	Status         Status          `json:"status"`
	Attempts       int             `json:"attempts"`
	MaxAttempts    int             `json:"max_attempts"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}
