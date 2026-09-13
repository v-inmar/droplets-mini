package main

import (
	"context"
	"droplets_mini/app/internal/service"

	"droplets_mini/pkg/database"
	"droplets_mini/pkg/stream"
	"droplets_mini/pkg/utilities"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx := context.Background()

	// load env variables
	_ = godotenv.Load("app/.env")

	// -------- postgres -------- //
	dbServer := database.NewPostgresDB(os.Getenv("DB_DSN"))
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Fatalf("[App] unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Printf("[App] connected to database\n")

	// ----- migrate ---------- //
	log.Printf("[App] migrating database models")
	migrateSrc := "file://app/migrations"
	isContainer := strings.ToLower(os.Getenv("IS_CONTAINER")) == "true"
	if isContainer {
		migrateSrc = "file://migrations"
	}
	pm := database.NewPostgresMigrate(migrateSrc, os.Getenv("DB_DSN"))
	if err := _migrate(pm); err != nil {
		log.Fatalf("[App] error while migrating: %v", err)
	}
	log.Printf("[App] successfully migrated")

	// ------------- nats ---------- //
	nc, err := utilities.NATSConnect(ctx, os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("[App] NATS connection failed: %v", err)
	}
	defer nc.Close()
	log.Printf("[App] connected to NATS server with the url: %s\n", nc.ConnectedUrl())

	// ------------ jetstream ----------- //
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("[App] unable to create jetstream instance")
	}
	log.Println("[App] initialized Jetstream")

	// ----------- jetstream event service ------------ //
	jsService := service.NewJetstreamService(js)
	if err := jsService.CreateOrUpdateStream(ctx, stream.StreamTodo); err != nil {
		log.Fatalf("[App] failed to create/update stream: %v", err)
	}
	log.Printf("[App] connected to stream '%s'", stream.StreamTodo)

	// ---------- new db service ----- //
	dbService := service.NewDBService(db)

	// -------- some test values ------- //
	userTodos := []string{
		"Go to sleep",
		"Do some golang coding",
		"Walk for 1 hour",
		"Eat",
		"Eat some more",
	}

	for _, t := range userTodos {
		log.Printf("[App] inserting '%s' into the database...", t)
		todo, err := dbService.CreateNewTodo(ctx, t)
		if err != nil {
			log.Printf("[App] unable to save new todo in database: %v", err)
			continue
		}
		log.Printf("[App] '%v' inserted in database", todo.Value)

		// publish event todo.created with fat event
		log.Printf("[App] publishing into the queue...")
		if err := jsService.PublishTodoCreated(ctx, todo); err != nil {
			log.Printf("[App] unable to publish todo created event: %v", err)
			continue
		}
		log.Printf("[App] event published for todo: '%v'", todo.Value)
	}

}

func _migrate(m database.Migrate) error {
	return m.Run()
}
