package messaging

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	HTPPServerAddress string
	DatabaseURL       string
}

// TODO: implement
func loadConfig() Config {
	return Config{}
}

func main() {
	// cfg := loadConfig() //TODO: implement
	_, cancel := context.WithCancel(context.Background())

	logger := log.New(os.Stdout, "[MessagingService] ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting message service...")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down the server")
	cancel()

	_, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
}
