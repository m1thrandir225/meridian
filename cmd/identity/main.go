package identity

import (
	"context"
	"log"
	"os"
)

type Config struct{}

func loadConfig() Config {
	return Config{}
}

func main() {
	_ = loadConfig()

	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(os.Stdout, "[IdentityService] ", log.LstdFlags|log.Lshortfile)

	logger.Println("Starting identity service...")
}
