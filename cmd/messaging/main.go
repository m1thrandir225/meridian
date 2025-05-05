package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/m1thrandir225/meridian/internal/messaging/application"
	kafkainfra "github.com/m1thrandir225/meridian/internal/messaging/infrastructure/kafka"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
)

type Config struct {
	HTPPServerAddress string
	DatabaseURL       string
	KafkaBrokers      []string
	KafkaDefaultTopic string
}

// TODO: implement
func loadConfig() Config {
	return Config{}
}

func main() {
	cfg := loadConfig() // TODO: implement
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(os.Stdout, "[MessagingService] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting message service...")

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
	}
	defer dbPool.Close()

	// --- Kafka Producers ---
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, config)
	if err != nil {
		logger.Fatalf("Failed to create Kafka sync producer: %v", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			logger.Printf("ERROR closing Kafka producer: %v", err)
		}
	}()

	logger.Println("Kafka sync producer initialized.")

	eventPublisher := kafkainfra.NewSaramaEventPublisher(producer, cfg.KafkaDefaultTopic)
	logger.Println("Kafka event publisher initialized.")

	repository := persistence.NewPostgresChannelRepository(dbPool)

	_ = application.NewChannelService(repository, eventPublisher)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down the server")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
