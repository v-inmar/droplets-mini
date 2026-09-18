package main

import (
	"context"
	"droplets_mini/pkg/database"
	dbrepo "droplets_mini/services/todo/internal/db_repo"
	"droplets_mini/services/todo/internal/handlers"
	"droplets_mini/services/todo/internal/services"
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
	dbServer := database.NewPostgresDB(os.Getenv("DSN"))
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Fatalf("[Todo] unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Printf("[Todo] connected to database\n")

	// ----- migrate ---------- //
	log.Printf("[Todo] migrating database models")
	migrateSrc := fmt.Sprintf("file://%s", os.Getenv("MIGRATION_SRC"))
	pm := database.NewPostgresMigrate(migrateSrc, os.Getenv("DSN"))
	if err := pm.Run(); err != nil {
		log.Fatalf("[Todo] error while migrating: %v", err)
	}

	log.Printf("[Todo] successfully migrated")

	router := chi.NewRouter()

	// handlers
	// handler := handler.NewHandlerRepo(db, service.NewTaskService(db, ctx))

	// router.Get("/health", handler.HealthCheckHandler)
	// router.Post("/tasks", handler.PostHandler)
	// router.Get("/tasks", handler.GetAllHandler)
	// router.Put("/tasks/{pid}", handler.PutHandler)
	// router.Delete("/tasks/{pid}", handler.DeleteHandler)

	healthHandler := handlers.NewHealthHandler(
		services.NewHealthService(),
	)
	router.Get("/health", healthHandler.GetHealthHandler)

	todoHandler := handlers.NewTodoHandler(
		services.NewTaskService(
			dbrepo.NewPostgresTaskDBRepo(db),
		),
	)
	router.Post("/tasks", todoHandler.PostCreateTask)

	log.Print("[Todo] service up and running...\n")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router); err != nil {
		log.Fatalf("[Todo] error: %v", err)
	}
}
