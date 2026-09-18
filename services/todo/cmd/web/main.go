package main

import (
	"context"
	dbrepo "droplets_mini_todoservice/internal/db_repo"
	"droplets_mini_todoservice/internal/handlers"
	"droplets_mini_todoservice/internal/services"
	"droplets_mini_todoservice/internal/utils"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load("services/todo/.env")

	// -------- postgres -------- //
	dbServer := utils.NewPostgresDB(os.Getenv("DSN"))
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Fatalf("[Todo] unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Printf("[Todo] connected to database\n")

	// ----- migrate ---------- //
	log.Printf("[Todo] migrating database models")
	pm := utils.NewPostgresMigrate("file://migrations", os.Getenv("DSN"))
	if err := pm.Run(); err != nil {
		log.Fatalf("[Todo] error while migrating: %v", err)
	}

	log.Printf("[Todo] successfully migrated")

	router := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(
		services.NewHealthService(),
	)
	router.Get("/health", healthHandler.GetHealthHandler)

	todoHandler := handlers.NewTodoHandler(
		services.NewTaskService(
			dbrepo.NewPostgresTaskDBRepo(db),
		),
	)
	router.Post("/tasks", todoHandler.PostCreateTaskHandler)
	router.Get("/tasks", todoHandler.GetAllTaskHandler)
	router.Put("/tasks/{pid}", todoHandler.PutUpdateTaskHandler)
	router.Delete("/tasks/{pid}", todoHandler.DeleteTaskHandler)

	log.Print("[Todo] service up and running...\n")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router); err != nil {
		log.Fatalf("[Todo] error: %v", err)
	}
}
