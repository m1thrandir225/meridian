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
	"github.com/m1thrandir225/meridian/pkg/kafka"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m1thrandir225/meridian/internal/identity/application/handlers"
	"github.com/m1thrandir225/meridian/internal/identity/application/services"
	"github.com/m1thrandir225/meridian/internal/identity/infrastructure/persistence"
)

type Config struct {
	HTTPPort             string
	KafkaBrokers         []string
	DatabaseURL          string
	PasetoPublicKey      string
	PasetoPrivateKey     string
	AuthTokenValidity    time.Duration
	RefreshTokenValidity time.Duration
	KafkaDefaultTopic    string
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

	return &Config{
		DatabaseURL:          dbURL,
		KafkaBrokers:         strings.Split(kafkaBrokerStr, ","),
		HTTPPort:             httpPort,
		PasetoPrivateKey:     privKey,
		PasetoPublicKey:      pubKey,
		AuthTokenValidity:    tokenValidity,
		KafkaDefaultTopic:    kafkaDefaultTopic,
		RefreshTokenValidity: refreshTokenValidity,
	}, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(os.Stdout, "[IdentityService] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting identity service...")

	if err := os.Setenv("IDENTITY_PASETO_PUBLIC_KEY", cfg.PasetoPublicKey); err != nil {
		log.Fatalf("Failed to set INTEGRITY_PASETO_PUBLIC_KEY for shared auth: %v", err)
	}

	if err := auth.LoadPublicKeyFromEnv(); err != nil {
		log.Fatalf("Failed to load PASETO public key for shared verifier: %v", err)
	}

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

	repository := persistence.NewPostgresUserRepository(dbPool)
	pasetoGenerator, err := token_generator.NewPasetoTokenGenerator(cfg.PasetoPrivateKey)
	if err != nil {
		log.Fatalf("Failed to create token generator: %v", err)
	}
	eventPublisher := kafka.NewSaramaEventPublisher(syncProducer, cfg.KafkaDefaultTopic)

	service := services.NewUserService(repository, pasetoGenerator, cfg.AuthTokenValidity, cfg.RefreshTokenValidity, eventPublisher)

	tokenVerifier, err := auth.NewPasetoTokenVerifier()

	router := handlers.SetupIdentityRouter(service, tokenVerifier)
	httpServer := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	errChan := make(chan error, 1)

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
