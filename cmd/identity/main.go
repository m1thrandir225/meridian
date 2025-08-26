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
	"github.com/m1thrandir225/meridian/internal/identity/infrastructure/token_generator"
	"github.com/m1thrandir225/meridian/pkg/auth"
	"github.com/m1thrandir225/meridian/pkg/cache"
	"github.com/m1thrandir225/meridian/pkg/kafka"
	"github.com/m1thrandir225/meridian/pkg/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/identity/application/handlers"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/infrastructure/persistence"
)

type Config struct {
	HTTPPort             string
	KafkaBrokers         []string
	DatabaseURL          string
	RedisURL             string
	PasetoPublicKey      string
	PasetoPrivateKey     string
	AuthTokenValidity    time.Duration
	RefreshTokenValidity time.Duration
	KafkaDefaultTopic    string
	IntegrationGRPCURL   string
	GRPCPort             string
	Environment          string
	LogLevel             string
}

func loadConfig() (*Config, error) {
	dbURL := os.Getenv("IDENTITY_DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("missing IDENTITY_DB_URL")
	}

	kafkaBrokerStr := os.Getenv("IDENTITY_KAFKA_BROKERS")
	if kafkaBrokerStr == "" {
		return nil, fmt.Errorf("missing IDENTITY_KAFKA_BROKERS")
	}
	kafkaDefaultTopic := os.Getenv("IDENTITY_KAFKA_DEFAULT_TOPIC")
	if kafkaDefaultTopic == "" {
		return nil, fmt.Errorf("missing IDENTITY_KAFKA_DEFAULT_TOPIC")
	}

	httpPort := os.Getenv("IDENTITY_HTTP_PORT")
	if httpPort == "" {
		return nil, fmt.Errorf("missing IDENTITY_HTTP_PORT")
	}

	privKey := os.Getenv("IDENTITY_PASETO_PRIVATE_KEY")
	if privKey == "" {
		return nil, fmt.Errorf("missing IDENTITY_PASETO_PRIVATE_KEY")
	}
	pubKey := os.Getenv("IDENTITY_PASETO_PUBLIC_KEY")
	if pubKey == "" {
		return nil, fmt.Errorf("missing IDENTITY_PASETO_PUBLIC_KEY")
	}
	integrationGRPCURL := os.Getenv("INTEGRATION_GRPC_URL")
	if integrationGRPCURL == "" {
		return nil, fmt.Errorf("missing INTEGRATION_GRPC_URL")
	}
	tokenValidityStr := os.Getenv("AUTH_TOKEN_VALIDITY_MINUTES")
	tokenValidity := 15 * time.Minute
	if val, err := time.ParseDuration(tokenValidityStr + "m"); err == nil {
		tokenValidity = val
	} else if tokenValidityStr != "" {
		log.Printf("WARN: Invalid AUTH_TOKEN_VALIDITY_MINUTES '%s', using default %v", tokenValidityStr, tokenValidity)
	}

	refreshTokenValidityStr := os.Getenv("REFRESH_TOKEN_VALIDITY_MINUTES")
	refreshTokenValidity := (24 * time.Hour) * 7
	if val, err := time.ParseDuration(refreshTokenValidityStr + "h"); err == nil {
		refreshTokenValidity = val
	} else if refreshTokenValidityStr != "" {
		log.Printf("WARN: Invalid REFRESH_TOKEN_VALIDITY_MINUTES '%s', using default %v", refreshTokenValidityStr, refreshTokenValidity)
	}

	grpcPort := os.Getenv("IDENTITY_GRPC_PORT")
	if grpcPort == "" {
		return nil, fmt.Errorf("missing IDENTITY_GRPC_PORT")
	}

	redisURL := os.Getenv("IDENTITY_REDIS_URL")
	if redisURL == "" {
		return nil, fmt.Errorf("missing IDENTITY_REDIS_URL")
	}

	environment := os.Getenv("IDENTITY_ENVIRONMENT")
	if environment == "" {
		environment = "development"
		fmt.Printf("WARN: IDENTITY_ENVIRONMENT is not set, using default %s\n", environment)
	}

	logLevel := os.Getenv("IDENTITY_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
		fmt.Printf("WARN: IDENTITY_LOG_LEVEL is not set, using default %s\n", logLevel)
	}

	return &Config{
		DatabaseURL:          dbURL,
		KafkaBrokers:         strings.Split(kafkaBrokerStr, ","),
		HTTPPort:             httpPort,
		RedisURL:             redisURL,
		PasetoPrivateKey:     privKey,
		PasetoPublicKey:      pubKey,
		AuthTokenValidity:    tokenValidity,
		KafkaDefaultTopic:    kafkaDefaultTopic,
		RefreshTokenValidity: refreshTokenValidity,
		IntegrationGRPCURL:   integrationGRPCURL,
		GRPCPort:             grpcPort,
		Environment:          environment,
		LogLevel:             logLevel,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(
		logging.Config{
			ServiceName: "[IdentityService]",
			Environment: cfg.Environment,
			LogLevel:    cfg.LogLevel,
		},
	)
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())

	logger.Info("Starting service...")

	if err := os.Setenv("IDENTITY_PASETO_PUBLIC_KEY", cfg.PasetoPublicKey); err != nil {
		logger.Fatal("Failed to set INTEGRITY_PASETO_PUBLIC_KEY for shared auth", zap.Error(err))
	}

	if err := auth.LoadPublicKeyFromEnv(); err != nil {
		logger.Fatal("Failed to load PASETO public key for shared verifier", zap.Error(err))
	}

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
	if err != nil {
		logger.Fatal("Failed to create Redis cache", zap.Error(err))
	}

	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaCfg)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}
	defer func() { _ = syncProducer.Close() }()
	logger.Info("Kafka producer established.")

	repository := persistence.NewPostgresUserRepository(dbPool)
	pasetoGenerator, err := token_generator.NewPasetoTokenGenerator(cfg.PasetoPrivateKey)
	if err != nil {
		logger.Fatal("Failed to create token generator", zap.Error(err))
	}
	eventPublisher := kafka.NewSaramaEventPublisher(syncProducer, cfg.KafkaDefaultTopic)

	service := services.NewUserService(
		repository,
		pasetoGenerator,
		cfg.AuthTokenValidity,
		cfg.RefreshTokenValidity,
		eventPublisher,
		logger,
	)

	if err := service.EnsureDefaultAdmin(ctx); err != nil {
		logger.Fatal("Failed to ensure default admin user exists", zap.Error(err))
		return
	}

	tokenVerifier, err := auth.NewPasetoTokenVerifier()
	if err != nil {
		logger.Fatal("Failed to create token verifier", zap.Error(err))
	}

	router := handlers.SetupIdentityRouter(
		service,
		redisCache,
		tokenVerifier,
		cfg.IntegrationGRPCURL,
		logger,
	)

	httpServer := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errChan := make(chan error, 1)

	go func() {
		logger.Info("Starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := handlers.StartGRPCServer(
			cfg.GRPCPort,
			tokenVerifier,
			service,
			redisCache,
			logger,
		); err != nil {
			logger.Fatal("Failed to start gRPC server", zap.Error(err))
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
