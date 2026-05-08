package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/drumilbhati/disscheduler/model"
	"github.com/drumilbhati/disscheduler/store"
)

type JobHandler struct {
	store *store.Store
}

func NewJobHandler(s *store.Store) *JobHandler {
	return &JobHandler{store: s}
}

func (j *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job model.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = j.store.CreateJob(&job)
	if err != nil {
		log.Printf("failed to create job: %v", err)
		http.Error(w, "Failed to create job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func (j *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := j.store.GetAllJobs()
	if err != nil {
		log.Printf("failed to fetch jobs: %v", err)
		http.Error(w, "Failed to fetch jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}
