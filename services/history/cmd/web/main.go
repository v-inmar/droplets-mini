package main

import (
	"context"
	"droplets_mini/pkg/database"
	dbrepo "droplets_mini/services/history/internal/db_repo"
	"droplets_mini/services/history/internal/handlers"
	"droplets_mini/services/history/internal/services"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load("services/history/.env")

	// -------- postgres -------- //
	dbServer := database.NewPostgresDB(os.Getenv("DSN"))
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Fatalf("[History] unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Printf("[History] connected to database\n")

	// ----- migrate ---------- //
	log.Printf("[History] migrating database models")
	migrateSrc := fmt.Sprintf("file://%s", os.Getenv("MIGRATION_SRC"))
	pm := database.NewPostgresMigrate(migrateSrc, os.Getenv("DSN"))
	if err := pm.Run(); err != nil {
		log.Fatalf("[History] error while migrating: %v", err)
	}

	log.Printf("[History] successfully migrated")

	router := chi.NewRouter()

	healthHandler := handlers.NewHealthHandler(services.NewHealthService())
	router.Get("/health", healthHandler.GetHealthHandler)

	historyHandler := handlers.NewHistoryHandler(
		services.NewHistoryService(
			dbrepo.NewPostgresHistoryDBRepo(db),
		),
	)
	router.Get("/items", historyHandler.GetAllHistoryHandler)

	log.Print("[History] service up and running...\n")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router); err != nil {
		log.Fatalf("[History] error: %v", err)
	}
}
