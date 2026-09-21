package main

import (
	"context"
	"database/sql"
	"droplets_mini_eventservice/internal/broker"
	dbrepo "droplets_mini_eventservice/internal/db_repo"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"
	"droplets_mini_eventservice/internal/services"
	"droplets_mini_eventservice/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load("services/event/.env")

	dsn := os.Getenv("DSN")
	log.Print("[Event] connecting to database server...")
	dbServer := utils.NewPostgresDB(dsn)
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Fatalf("[Event] unable to connect to database server: %s, error: %v", dsn, err)
	}
	defer db.Close()
	log.Printf("[Event] connected to database server: %s", dsn)

	log.Printf("[Event] migrating database models...")
	migrate := utils.NewPostgresMigrate("file://migrations", dsn)
	if err := migrate.Run(); err != nil {
		log.Fatalf("[Event] error while migrating to database: %s, error: %v", dsn, err)
	}
	log.Print("[Event] migrated successfully")

	historyURL := os.Getenv("HISTORY_URL")

	kafkaAddr := os.Getenv("KAFKA_BROKERS")
	kafkaBrokerReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddr},
		Topic:   "tasks",
		GroupID: "task-service",
	})
	kafkaBrokerWriterDLQ := kafka.Writer{
		Addr:  kafka.TCP(kafkaAddr),
		Topic: "tasks.dlq",
	}

	consumer := broker.NewKafkaBrokerConsumer(kafkaBrokerReader)
	defer consumer.Close()

	dlqProducer := broker.NewKafkaBrokerProducer(&kafkaBrokerWriterDLQ)
	defer dlqProducer.Close()

	postgresRepo := dbrepo.NewPostgresEventProcessedDBRepo(db)
	procEventService := services.NewProcessEventService(postgresRepo)

	for {
		msg, err := consumer.Consume(ctx)
		if err != nil {
			log.Fatalf("[Event] unable to consume: %v", err)
		}

		var eventTask eventmodels.EventTask
		if err := json.Unmarshal(msg.Value, &eventTask); err != nil {
			log.Fatalf("[Event] failed to unmarshal message event: %v", err)
		}

		model, err := postgresRepo.ReadByEventPID(ctx, eventTask.EventPID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				log.Fatalf("[Event] unable to read event processed database model: %v", err)
			}
		}

		if model == nil {
			model, err = procEventService.CreateEventProcess(ctx, &eventTask)
			if err != nil {
				log.Fatalf("[Event] create event process service failed: %v", err)
			}
		}

		// already processed?
		if model.EventProcessed {
			if err := consumer.Commit(ctx, msg); err != nil {
				log.Fatalf("[Event] event already processed kafka commit failed: %v", err)
			}
			continue
		}

		// max retry already reached
		if model.EventRetry >= 5 {
			if err := dlqProducer.Produce(ctx, eventTask, msg.Key); err != nil {
				log.Fatalf("[Event] error while publishing to dlq: %v", err)
			}

			if err := consumer.Commit(ctx, msg); err != nil {
				log.Fatalf("[Event] event retry max reached kafka commit failed: %v", err)
			}
			continue
		}

		updatedModel, err := procEventService.UpdateEventProcess(ctx, eventTask.EventPID, model.EventRetry+1, false)
		if err != nil {
			log.Fatalf("[Event] error incrementing event retry: %v", err)
		}

		// synchronous call -- because i need the response to continue the kafka flow
		// this is ok since it doesnt take too long in my setup
		if err := procEventService.PostHistoryAPI(ctx, fmt.Sprintf("%s/items", historyURL), eventTask); err != nil {
			log.Printf("[Event] unable to post to history api: %v", err)
			continue
		}

		_, err = procEventService.UpdateEventProcess(ctx, updatedModel.EventPID, updatedModel.EventRetry, true)
		if err != nil {
			log.Printf("[Event] error completing event processed: %v", err)
			continue
		}

		if err := consumer.Commit(ctx, msg); err != nil {
			log.Printf("[Event] event retry max reached kafka commit failed: %v", err)
		}

	}

}
