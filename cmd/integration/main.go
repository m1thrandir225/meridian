package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/integration/application/handlers"
	"github.com/m1thrandir225/meridian/internal/integration/application/services"
	"github.com/m1thrandir225/meridian/internal/integration/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/kafka"
)

type Config struct {
	HTTPPort          string
	KafkaBrokers      []string
	KafkaDefaultTopic string
	DatabaseURL       string
	GRPCPort          string
	MessagingGRPCURL  string
}

func loadConfig() (*Config, error) {
	dbURL := os.Getenv("INTEGRATION_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("missing INTEGRATION_DB_URL")
	}

	kafkaBrokerStr := os.Getenv("INTEGRATION_KAFKA_BROKERS")
	if kafkaBrokerStr == "" {
		return nil, fmt.Errorf("missing INTEGRATION_KAFKA_BROKERS")
	}
	kafkaDefaultTopic := os.Getenv("INTEGRATION_KAFKA_DEFAULT_TOPIC")
	if kafkaDefaultTopic == "" {
		return nil, fmt.Errorf("missing INTEGRATION_KAFKA_DEFAULT_TOPIC")
	}
	httpPort := os.Getenv("INTEGRATION_HTTP_PORT")
	if httpPort == "" {
		return nil, fmt.Errorf("missing INTEGRATION_HTTP_PORT")
	}
	grpcPort := os.Getenv("INTEGRATION_GRPC_PORT")
	if grpcPort == "" {
		return nil, fmt.Errorf("missing INTEGRATION_GRPC_PORT")
	}

	messagingGRPCURL := os.Getenv("MESSAGING_GRPC_URL")
	if messagingGRPCURL == "" {
		return nil, fmt.Errorf("missing MESSAGING_GRPC_URL")
	}

	return &Config{
		HTTPPort:          httpPort,
		KafkaBrokers:      strings.Split(kafkaBrokerStr, ","),
		KafkaDefaultTopic: kafkaDefaultTopic,
		DatabaseURL:       dbURL,
		GRPCPort:          grpcPort,
		MessagingGRPCURL:  messagingGRPCURL,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(os.Stdout, "[Integration Service] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting integration service")

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaCfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func() { _ = syncProducer.Close() }()
	log.Println("Kafka producer established.")

	repository := persistence.NewPostgresIntegrationRepository(dbPool)

	bcryptGenerator := services.NewBcryptTokenGenerator()
	eventPublisher := kafka.NewSaramaEventPublisher(syncProducer, cfg.KafkaDefaultTopic)

	service := services.NewIntegrationService(repository, bcryptGenerator, eventPublisher)

	router := handlers.SetupIntegrationRouter(service)
	httpServer := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errChan := make(chan error, 1)

	go func() {
		log.Printf("Starting gRPC Server on %s", cfg.GRPCPort)
		if err := handlers.StartGRPCServer(service, cfg.GRPCPort); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	go func() {
		log.Printf("Starting Identity HTTP server on %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("identity HTTP server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down http server.")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
