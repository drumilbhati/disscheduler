package main

import (
	"log"
	"net/http"
	"os"

	"github.com/drumilbhati/disscheduler/controller"
	"github.com/drumilbhati/disscheduler/db"
	"github.com/drumilbhati/disscheduler/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Connect to database
	godotenv.Load()
	db, err := db.Connect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	port := os.Getenv("PORT")
	if err != nil {
		log.Fatal("Failed to connect to database, err:", err)
	}
	// Close database connection
	defer db.Close()

	// Create a new store to interact with the database
	s := store.NewStore(db)

	// Create handlers
	j := controller.NewJobHandler(s)

	// Create a new router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/job", j.GetAllJobs)
	r.Post("/job", j.CreateJob)

	log.Fatal(http.ListenAndServe(port, r))
}
