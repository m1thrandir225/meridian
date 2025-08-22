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

	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/m1thrandir225/meridian/internal/messaging/application/handlers"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
	"github.com/m1thrandir225/meridian/internal/messaging/infrastructure/persistence"
)

type Config struct {
	HTTPPort           string
	DatabaseURL        string
	KafkaBrokers       []string
	KafkaDefaultTopic  string
	GRPCPort           string
	IdentityGRPCURL    string
	IntegrationGRPCURL string
	RedisURL           string
	Environment        string
	LogLevel           string
}

func loadConfig() (*Config, error) {
	kafkaBrokerStr := os.Getenv("MESSAGING_KAFKA_BROKERS")
	if kafkaBrokerStr == "" {
		return nil, fmt.Errorf("missing MESSAGING_KAFKA_BROKERS")
	}

	httpPort := os.Getenv("MESSAGING_HTTP_PORT")
	if httpPort == "" {
		return nil, fmt.Errorf("missing MESSAGING_HTTP_PORT")
	}

	kafkaDefaultTopic := os.Getenv("MESSAGING_KAFKA_DEFAULT_TOPIC")
	if kafkaDefaultTopic == "" {
		return nil, fmt.Errorf("missing MESSAGING_KAFKA_DEFAULT_TOPIC")
	}

	dbURL := os.Getenv("MESSAGING_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("missing MESSAGING_DB_URL")
	}

	grpcPort := os.Getenv("MESSAGING_GRPC_PORT")
	if grpcPort == "" {
		return nil, fmt.Errorf("missing MESSAGING_GRPC_PORT")
	}

	identityGRPCURL := os.Getenv("IDENTITY_GRPC_URL")
	if identityGRPCURL == "" {
		return nil, fmt.Errorf("missing IDENTITY_GRPC_URL")
	}

	redisURL := os.Getenv("MESSAGING_REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("missing MESSAGING_REDIS_URL")
	}

	integrationGRPCURL := os.Getenv("INTEGRATION_GRPC_URL")
	if integrationGRPCURL == "" {
		return nil, fmt.Errorf("missing INTEGRATION_GRPC_URL")
	}

	environment := os.Getenv("MESSAGING_ENVIRONMENT")
	if environment == "" {
		environment = "development"
		fmt.Printf("WARN: MESSAGING_ENVIRONMENT is not set, using default %s\n", environment)
	}

	level := os.Getenv("MESSAGING_LOG_LEVEL")
	if level == "" {
		level = "info"
		fmt.Printf("WARN: MESSAGING_LOG_LEVEL is not set, using default %s\n", level)
	}

	return &Config{
		HTTPPort:           httpPort,
		DatabaseURL:        dbURL,
		KafkaBrokers:       strings.Split(kafkaBrokerStr, ","),
		KafkaDefaultTopic:  kafkaDefaultTopic,
		GRPCPort:           grpcPort,
		IdentityGRPCURL:    identityGRPCURL,
		RedisURL:           redisURL,
		IntegrationGRPCURL: integrationGRPCURL,
		Environment:        environment,
		LogLevel:           level,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	logger := logging.NewLogger(logging.Config{
		ServiceName: "[MessagingService]",
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})
	logger.Info("Starting service...")

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
		DB:   0,
	})
	defer redisClient.Close()
	redisCache := cache.NewRedisCache(redisClient)

	// --- Kafka Producers ---
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, config)
	if err != nil {
		logger.Fatal("Failed to create Kafka sync producer", zap.Error(err))
	}

	defer func() {
		if err := producer.Close(); err != nil {
			logger.Error("ERROR closing Kafka producer", zap.Error(err))
		}
	}()

	logger.Info("Kafka sync producer initialized.")

	eventPublisher := kafka.NewSaramaEventPublisher(producer, cfg.KafkaDefaultTopic)
	logger.Info("Kafka event publisher initialized.")

	repository := persistence.NewPostgresChannelRepository(dbPool)
	logger.Info("Database pool initialized.")

	identityClient, err := services.NewIdentityClient(cfg.IdentityGRPCURL)
	if err != nil {
		logger.Fatal("Failed to create identity client", zap.Error(err))
	}
	defer identityClient.Close()

	integrationClient, err := services.NewIntegrationClient(cfg.IntegrationGRPCURL)
	if err != nil {
		logger.Fatal("Failed to create integration client", zap.Error(err))
	}
	defer integrationClient.Close()

	channelService := services.NewChannelService(
		repository,
		eventPublisher,
		identityClient,
		integrationClient,
		logger,
	)
	logger.Info("Channel service initialized.")

	messageService := services.NewMessageService(
		repository,
		eventPublisher,
		identityClient,
		integrationClient,
		logger,
	)
	logger.Info("Message service initialized.")

	httpHandler := handlers.NewHttpHandler(
		channelService,
		messageService,
		redisCache,
		logger,
	)
	logger.Info("HTTP Handler initialized")

	wsHandler := handlers.NewWebSocketHandler(
		channelService,
		messageService,
		redisClient,
		identityClient,
		logger,
	)
	logger.Info("WebSocket Handler initialized")

	// -- GIN ROUTE SETUP --
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(logging.GinLoggingMiddleware(logger))
	router.Use(logging.GinRecoveryMiddleware(logger))

	handlers.SetupRoutes(router, httpHandler, wsHandler)
	logger.Info("HTTP Routes initialized")

	// -- HTTP SERVER  --
	httpServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Info("Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := handlers.StartGRPCServer(
			cfg.GRPCPort,
			channelService,
			messageService,
			wsHandler,
			redisCache,
			logger,
		); err != nil {
			logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	go func() {
		logger.Info("HTTP Server listening", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down the server")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
