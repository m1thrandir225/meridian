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
	"github.com/m1thrandir225/meridian/internal/analytics/application/handlers"
	"github.com/m1thrandir225/meridian/internal/analytics/application/services"
	"github.com/m1thrandir225/meridian/internal/analytics/infrastructure/persistence"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"go.uber.org/zap"
)

type Config struct {
	HTTPPort      string
	KafkaBrokers  []string
	ConsumerGroup string
	DatabaseURL   string
	Environment   string
	LogLevel      string
}

func loadConfig() (*Config, error) {
	kafkaBrokerStr := os.Getenv("ANALYTICS_KAFKA_BROKERS")
	if kafkaBrokerStr == "" {
		return nil, fmt.Errorf("missing ANALYTICS_KAFKA_BROKERS")
	}

	httpPort := os.Getenv("ANALYTICS_HTTP_PORT")
	if httpPort == "" {
		return nil, fmt.Errorf("missing ANALYTICS_HTTP_PORT")
	}

	consumerGroup := os.Getenv("ANALYTICS_CONSUMER_GROUP")
	if consumerGroup == "" {
		consumerGroup = "analytics-service"
	}

	dbURL := os.Getenv("ANALYTICS_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("missing ANALYTICS_DB_URL")
	}

	environment := os.Getenv("ANALYTICS_ENVIRONMENT")
	if environment == "" {
		environment = "development"
		fmt.Printf("WARN: ANALYTICS_ENVIRONMENT is not set, using default %s\n", environment)
	}

	level := os.Getenv("ANALYTICS_LOG_LEVEL")
	if level == "" {
		level = "info"
		fmt.Printf("WARN: ANALYTICS_LOG_LEVEL is not set, using default %s\n", level)
	}

	return &Config{
		HTTPPort:      httpPort,
		KafkaBrokers:  strings.Split(kafkaBrokerStr, ","),
		ConsumerGroup: consumerGroup,
		DatabaseURL:   dbURL,
		Environment:   environment,
		LogLevel:      level,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.NewLogger(logging.Config{
		ServiceName: "[AnalyticsService]",
		Environment: cfg.Environment,
		LogLevel:    cfg.LogLevel,
	})
	logger.Info("Starting Analytics Service...")

	// Initialize database
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	// Initialize repository
	repository := persistence.NewPostgresAnalyticsRepository(dbPool)
	logger.Info("Database connection established")

	// Initialize analytics service
	analyticsService := services.NewAnalyticsService(repository, logger)
	logger.Info("Analytics service initialized")

	// Initialize event handler
	eventHandler := handlers.NewAnalyticsEventHandler(analyticsService, logger)
	logger.Info("Event handler initialized")

	// Initialize Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := kafka.NewEventConsumer(cfg.KafkaBrokers, cfg.ConsumerGroup, config)
	logger.Info("Kafka consumer initialized")

	topics := []string{
		"meridian.identity.events",
		"meridian.messaging.events",
		"meridian.integration.events",
	}

	// Start event consumer
	go func() {
		logger.Info("Starting Kafka consumer", zap.Strings("topics", topics), zap.Strings("kafka_brokers", cfg.KafkaBrokers), zap.String("consumer_group", cfg.ConsumerGroup))
		if err := consumer.ConsumeEvents(ctx, topics, eventHandler); err != nil {
			logger.Error("Error consuming events", zap.Error(err))
		}
	}()

	// Setup HTTP router
	router := handlers.SetupAnalyticsRouter(analyticsService, logger)
	logger.Info("HTTP router configured")

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		logger.Info("HTTP Server listening", zap.String("port", cfg.HTTPPort))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down the Analytics Service")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}

	logger.Info("Analytics Service stopped")
}
