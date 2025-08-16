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
	"github.com/redis/go-redis/v9"

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

	return &Config{
		HTTPPort:           httpPort,
		DatabaseURL:        dbURL,
		KafkaBrokers:       strings.Split(kafkaBrokerStr, ","),
		KafkaDefaultTopic:  kafkaDefaultTopic,
		GRPCPort:           grpcPort,
		IdentityGRPCURL:    identityGRPCURL,
		RedisURL:           redisURL,
		IntegrationGRPCURL: integrationGRPCURL,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(os.Stdout, "[MessagingService] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting message service...")

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
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
		logger.Fatalf("Failed to create Kafka sync producer: %v", err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			logger.Printf("ERROR closing Kafka producer: %v", err)
		}
	}()

	logger.Println("Kafka sync producer initialized.")

	eventPublisher := kafka.NewSaramaEventPublisher(producer, cfg.KafkaDefaultTopic)
	logger.Println("Kafka event publisher initialized.")

	repository := persistence.NewPostgresChannelRepository(dbPool)
	logger.Println("Database pool initialized.")

	identityClient, err := services.NewIdentityClient(cfg.IdentityGRPCURL)
	if err != nil {
		logger.Fatalf("Failed to create identity client: %v", err)
	}
	defer identityClient.Close()

	integrationClient, err := services.NewIntegrationClient(cfg.IntegrationGRPCURL)
	if err != nil {
		logger.Fatalf("Failed to create integration client: %v", err)
	}
	defer integrationClient.Close()

	channelService := services.NewChannelService(repository, eventPublisher, identityClient, integrationClient)
	logger.Println("Channel service initialized.")

	messageService := services.NewMessageService(repository, eventPublisher, identityClient, integrationClient)
	logger.Println("Message service initialized.")

	httpHandler := handlers.NewHttpHandler(channelService, messageService, redisCache)
	logger.Println("HTTP Handler initialized")

	wsHandler := handlers.NewWebSocketHandler(channelService, messageService, redisClient, identityClient)
	logger.Println("WebSocket Handler initialized")

	// -- GIN ROUTE SETUP --
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.SetupRoutes(router, httpHandler, wsHandler)
	logger.Println("HTTP Routes initialized")

	// -- HTTP SERVER  --
	httpServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Printf("Starting gRPC server on %s", cfg.GRPCPort)
		if err := handlers.StartGRPCServer(
			cfg.GRPCPort,
			channelService,
			messageService,
			wsHandler,
			redisCache,
		); err != nil {
			logger.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	go func() {
		logger.Printf("HTTP Server listening on %s", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down the server")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
