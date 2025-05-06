package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/m1thrandir225/meridian/internal/messaging/application/handlers"
	"github.com/m1thrandir225/meridian/internal/messaging/application/services"
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
	brokers := []string{"kafka:9092"} // Default kafka brokers

	if brokerList := os.Getenv("KAFKA_BROKERS"); brokerList != "" {
		brokers = strings.Split(brokerList, ",")
	}
	return Config{
		HTPPServerAddress: os.Getenv("HTTP_SERVER_ADDRESS"),
		KafkaBrokers:      brokers,
		KafkaDefaultTopic: os.Getenv("KAFKA_DEFAULT_TOPIC"),
		DatabaseURL:       os.Getenv("MESSAGING_DB_URL"),
	}
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

	service := services.NewChannelService(repository, eventPublisher)

	httpHandler := handlers.NewHttpHandler(service)

	// -- GIN ROUTE SETUP --
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiV1 := router.Group("/api/v1")
	{
		channelsGroup := apiV1.Group("/channels")
		{
			channelsGroup.POST("/", httpHandler.CreateChannel)
			channelsGroup.GET("/:channelId", httpHandler.GetChannel)
			channelsGroup.POST("/:channelId/join", httpHandler.JoinChannel)
			messagesGroup := channelsGroup.Group("/:channelId/messages")
			{
				messagesGroup.POST("/", httpHandler.SendMessage)
			}
		}
	}
	// -- HTTP SERVER  --
	httpServer := &http.Server{
		Addr:    cfg.HTPPServerAddress,
		Handler: router,
	}

	go func() {
		logger.Printf("HTTP Server listening on %s", cfg.HTPPServerAddress)
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
