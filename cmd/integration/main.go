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
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Config struct {
	HTTPPort          string
	KafkaBrokers      []string
	KafkaDefaultTopic string
	DatabaseURL       string
	RedisURL          string
	GRPCPort          string
	MessagingGRPCURL  string
	Environment       string
	LogLevel          string
}

func loadConfig() (*Config, error) {
	dbURL := os.Getenv("INTEGRATION_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("missing INTEGRATION_DB_URL")
	}
	redisURL := os.Getenv("INTEGRATION_REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("missing INTEGRATION_REDIS_URL")
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

	environment := os.Getenv("INTEGRATION_ENVIRONMENT")
	if environment == "" {
		environment = "development"
		fmt.Printf("WARN: INTEGRATION_ENVIRONMENT is not set, using default %s\n", environment)
	}

	level := os.Getenv("INTEGRATION_LOG_LEVEL")
	if level == "" {
		level = "info"
		fmt.Printf("WARN: INTEGRATION_LOG_LEVEL is not set, using default %s\n", level)
	}

	return &Config{
		HTTPPort:          httpPort,
		KafkaBrokers:      strings.Split(kafkaBrokerStr, ","),
		KafkaDefaultTopic: kafkaDefaultTopic,
		DatabaseURL:       dbURL,
		RedisURL:          redisURL,
		GRPCPort:          grpcPort,
		MessagingGRPCURL:  messagingGRPCURL,
		Environment:       environment,
		LogLevel:          level,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	logger := logging.NewLogger(logging.Config{
		ServiceName: "[IntegrationService]",
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})
	logger.Info("Starting integration service")

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	defer redisClient.Close()
	redisCache := cache.NewRedisCache(redisClient)

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaCfg)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer func() { _ = syncProducer.Close() }()
	logger.Info("Kafka producer established.")

	repository := persistence.NewPostgresIntegrationRepository(dbPool)

	bcryptGenerator := services.NewBcryptTokenGenerator()
	eventPublisher := kafka.NewSaramaEventPublisher(syncProducer, cfg.KafkaDefaultTopic)

	service := services.NewIntegrationService(
		repository,
		bcryptGenerator,
		eventPublisher,
		logger,
	)

	messageClient, err := services.NewMessagingClient(cfg.MessagingGRPCURL)
	if err != nil {
		logger.Fatal("Failed to create messaging client", zap.Error(err))
	}
	defer messageClient.Close()

	router := handlers.SetupIntegrationRouter(service, redisCache, messageClient, logger)
	httpServer := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errChan := make(chan error, 1)

	go func() {
		logger.Info("Starting gRPC Server", zap.String("port", cfg.GRPCPort))
		if err := handlers.StartGRPCServer(service, redisCache, logger, cfg.GRPCPort); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	go func() {
		logger.Info("Starting Identity HTTP server", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("identity HTTP server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down http server.")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
