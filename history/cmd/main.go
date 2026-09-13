package main

import (
	"context"
	"droplets_mini/history/internal/repo"
	"droplets_mini/pkg/database"
	"droplets_mini/pkg/events"
	"droplets_mini/pkg/stream"
	"droplets_mini/pkg/subjects"
	"droplets_mini/pkg/utilities"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	// signal for graceful shutdown from different executors i.e. K8s, ctrl+C, etc
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// load env variables
	_ = godotenv.Load("history/.env")

	// -------- postgres -------- //
	dbServer := database.NewPostgresDB(os.Getenv("DB_DSN"))
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Fatalf("[History] unable to connect to database: %v\n", err)
	}
	defer db.Close()
	log.Printf("[History] connected to database\n")

	// ----- migrate ---------- //
	log.Printf("[History] migrating database models")
	migrateSrc := "file://history/migrations"
	isContainer := strings.ToLower(os.Getenv("IS_CONTAINER")) == "true"
	if isContainer {
		migrateSrc = "file://migrations"
	}
	pm := database.NewPostgresMigrate(migrateSrc, os.Getenv("DB_DSN"))
	if err := _migrate(pm); err != nil {
		log.Fatalf("[History] error while migrating: %v", err)
	}
	log.Printf("[History] successfully migrated")

	// ------------- nats ---------- //
	nc, err := utilities.NATSConnect(ctx, os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatalf("[History] NATS connection failed: %v", err)
	}
	defer nc.Close()
	log.Printf("[History] connected to NATS server with the url: %s\n", nc.ConnectedUrl())

	// ------------ jetstream ----------- //
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("[History] unable to create jetstream instance")
	}
	log.Println("[History] initialized Jetstream")

	// ---------- consumer ---------------- //
	consumer, err := js.CreateOrUpdateConsumer(ctx, stream.StreamTodo, jetstream.ConsumerConfig{
		Durable: "history-service-durable",
		FilterSubjects: []string{
			subjects.SubjectTodoCreated,
			subjects.SubjectTodoCompleted,
			subjects.SubjectTodoDeleted,
			subjects.SubjectTodoUpdated,
		},
		AckPolicy:  jetstream.AckExplicitPolicy,
		MaxDeliver: 5,
	})
	if err != nil {
		log.Fatalf("[History] Consumer creation failed: %v", err)
	}

	// --------- consumer callback ------------ //
	cc, err := consumer.Consume(func(msg jetstream.Msg) {

		switch msg.Subject() {
		case subjects.SubjectTodoCreated:
			// unmarshal message data into event
			var event events.TodoCreatedEvent
			if err := json.Unmarshal(msg.Data(), &event); err != nil {
				log.Printf("[History] unable to unmarshal data: %v", err)
				return
			}

			// start transaction
			tx, err := db.BeginTxx(ctx, nil)
			if err != nil {
				log.Printf("[History] unable to start database transaction: %v", err)
				return
			}
			defer tx.Rollback()
			procTx := repo.NewProcessedRepoTx(tx, ctx)
			historyTx := repo.NewHistoryRepoTx(tx, ctx)

			_, err = procTx.Create(event.EventID)
			if err != nil {
				if errors.Is(err, repo.ErrEventIDExist) {
					log.Printf(
						"[History] event id %d already processed; acknowledging duplicate",
						event.EventID,
					)

					if err := msg.Ack(); err != nil {
						log.Printf("[History] acknowledge failed: %v", err)
					}

					return
				}
				log.Printf("[History] unable to insert event id into processed event table: %v", err)
				return
			}

			_, err = historyTx.Create("CREATED", event.TodoValue, event.TodoID)
			if err != nil {
				log.Printf("[History] unable to insert item into history table: %v", err)
				return
			}

			// commit happens here
			if err := tx.Commit(); err != nil {
				log.Printf("[History] error while commiting transaction for event id: %d. %v", event.EventID, err)
				return
			}

			if err := msg.Ack(); err != nil {
				log.Printf("[History] acknowledge failed: %v", err)
			}

			log.Printf("[History] '%s' for event %d - processed and inserted in history", subjects.SubjectTodoCreated, event.EventID)

		case subjects.SubjectTodoCompleted,
			subjects.SubjectTodoDeleted,
			subjects.SubjectTodoUpdated:

			if err := msg.Ack(); err != nil {
				log.Printf("[History] acknowledge failed: %v", err)
				return
			}

			log.Printf(
				"[History] '%s' acknowledged (dummy handler)",
				msg.Subject(),
			)

		}

	})
	if err != nil {
		log.Fatalf("[History] Consumer failed: %v", err)
	}
	defer cc.Stop()

	log.Println("[History] Durable JetStream listener running...")

	<-ctx.Done()

}

func _migrate(m database.Migrate) error {
	return m.Run()
}
