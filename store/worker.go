package store

import (
	"context"
	"log"
	"time"
)

func (s *Store) StartWorkers(n int, ctx context.Context) {
	for i := range n {
		go s.worker(i, ctx)
	}
}

func (s *Store) worker(id int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := s.ClaimJob(ctx)

		if err != nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
			}
			log.Printf("worker=%d drain error: %v\n", id, err)
			continue
		}
		if job == nil {
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
			continue
		}
		s.processJob(job, ctx)
	}
}
