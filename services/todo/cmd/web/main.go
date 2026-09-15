package main

import (
	"context"
	"droplets_mini/pkg/database"
	"droplets_mini/services/todo/internal/handler"
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
	handler := handler.NewHandlerRepo(db)

	router.Get("/health", handler.HealthCheckHandler)

	log.Print("[Todo] service up and running...\n")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router); err != nil {
		log.Fatalf("[Todo] error: %v", err)
	}
}
