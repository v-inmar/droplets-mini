package main

import (
	"context"
	"droplets_mini_eventservice/internal/broker"
	dbrepo "droplets_mini_eventservice/internal/db_repo"
	eventmodels "droplets_mini_eventservice/internal/models/event_models"
	"droplets_mini_eventservice/internal/services"
	"droplets_mini_eventservice/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	_ = godotenv.Load("services/event/.env")

	dsn := os.Getenv("DSN")
	log.Print("[Event] connecting to database server...")
	dbServer := utils.NewPostgresDB(dsn)
	db, err := dbServer.Connect(ctx)
	if err != nil {
		log.Printf("[Event] unable to connect to database server: %s, error: %v", dsn, err)
		return
	}
	defer db.Close()
	log.Printf("[Event] connected to database server: %s", dsn)

	log.Printf("[Event] migrating database models...")
	migrate := utils.NewPostgresMigrate("file://migrations", dsn)
	if err := migrate.Run(); err != nil {
		log.Printf("[Event] error while migrating to database: %s, error: %v", dsn, err)
		return
	}
	log.Print("[Event] migrated successfully")

	log.Printf("[Event] connecting to history database server...")
	historyDSN := os.Getenv("HISTORYDB_DSN")
	historyDBServer := utils.NewPostgresDB(historyDSN)
	historyDB, err := historyDBServer.Connect(ctx)
	if err != nil {
		log.Printf("[Event] unable to connect to database server: %s, error: %v", historyDSN, err)
		return
	}
	defer historyDB.Close()
	log.Printf("[Event] connected to database server: %s", historyDSN)

	kafkaAddr := os.Getenv("KAFKA_BROKERS")

	// Adding explicit MinBytes/MaxBytes to the kafka.ReaderConfig
	// resulted in the consumer receiving partition 0 and processing messages.
	kafkaBrokerReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaAddr},
		Topic:    "tasks",
		GroupID:  "task-service",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	kafkaBrokerWriterDLQ := kafka.Writer{
		Addr:  kafka.TCP(kafkaAddr),
		Topic: "tasks.dlq",
	}

	consumer := broker.NewKafkaBrokerConsumer(kafkaBrokerReader)
	defer consumer.Close()

	dlqProducer := broker.NewKafkaBrokerProducer(&kafkaBrokerWriterDLQ)
	defer dlqProducer.Close()

	postgresEPRepo := dbrepo.NewPostgresEventProcessedDBRepo(db)
	procEventService := services.NewProcessEventService(postgresEPRepo)

	postgresHSRepo := dbrepo.NewPostgresHistoryDBRepo(historyDB)
	historyService := services.NewHistoryService(postgresHSRepo)

	for {
		msg, err := consumer.Consume(ctx)
		if err != nil {

			// should be done to all services accepting context and returns error
			// future work
			if ctx.Err() != nil {
				log.Print("[Event] shutdown signal received")
				break
			}

			log.Printf("[Event] unable to consume: %v", err)
			continue
		}

		var eventTask eventmodels.EventTask
		if err := json.Unmarshal(msg.Value, &eventTask); err != nil {
			log.Printf("[Event] failed to unmarshal message event: %v", err)
			continue
		}

		eventProcModel, err := procEventService.GetOrCreateEventProcess(ctx, eventTask.EventPID)
		if err != nil {
			log.Printf("[Event] unbale to getorcreate event process model for event pid: %d with error: %v", eventTask.EventPID, err)
			continue
		}

		// already processed
		if eventProcModel.EventProcessed {
			if err := consumer.Commit(ctx, msg); err != nil {
				log.Printf("[Event] kafka commit error on already processed event: %v", err)
			}
			continue
		}

		// retry already reached
		// uses a constant to match retry limit
		if eventProcModel.EventRetry >= utils.MaxRetry {
			// publish in dlq
			// if this publish fails, it will retry
			if err := dlqProducer.Produce(ctx, eventTask, msg.Key); err != nil {
				log.Printf("[Event] unable to publish event in dlq: %v", err)
				continue
			}

			// if this commit fails, duplication is going to happen
			// needs better solution, this is ok for now
			if err := consumer.Commit(ctx, msg); err != nil {
				log.Printf("[Event] kafka commit error after publishing to DLQ: %v", err)
			}
			continue

			// kafka commit offset
			// continue
		}

		// increment retry by 1
		eventProcModel, err = procEventService.UpdateEventProcess(ctx, eventProcModel.EventPID, eventProcModel.EventRetry+1, false)
		if err != nil {
			log.Printf("[Event] error while incrementing retry: %v", err)
			continue
		}

		if err := historyService.CreateHistoryItem(ctx, &eventTask); err != nil {
			if errors.Is(err, dbrepo.ErrEvenPIDExist) {
				// prevent duplicate depending on event pid whihc is unique per event
				// so commit here
				if err := consumer.Commit(ctx, msg); err != nil {
					log.Printf("[Event] kafka commit error on already existing history item: %v", err)
				}
				continue
			}

			log.Printf("[Event] unable to create history item: %v", err)
			continue
		}

		// make EventProcessed true
		_, err = procEventService.UpdateEventProcess(ctx, eventProcModel.EventPID, eventProcModel.EventRetry, true)
		if err != nil {
			log.Printf("[Event] error while turning event processed into true: %v", err)
			continue
		}

		// commit kafka
		if err := consumer.Commit(ctx, msg); err != nil {
			log.Printf("[Event] end of execution kafka commit error: %v", err)
		}

	}

	log.Print("[Event] shutting down...")

}
