package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ApiPort    int
	SecretKey  string
	MongoURI   string
	CorsOrigin string
}

var Config = ServerConfig{}

func Load() error {
	envFile := ".env"
	if os.Getenv("E2E") == "true" {
		envFile = ".env.e2e"
	} else if os.Getenv("DOCKER") == "true" {
		envFile = ".env.docker"
	}
	godotenv.Load(envFile)

	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil || apiPort == 0 {
		return fmt.Errorf("invalid API_PORT env variable, got %v", os.Getenv("API_PORT"))
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("invalid SECRET_KEY env variable")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return fmt.Errorf("missing MONGO_URI env variable")
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		return fmt.Errorf("missing CORS_ORIGIN env variable")
	}

	Config = ServerConfig{
		ApiPort:    apiPort,
		SecretKey:  secretKey,
		MongoURI:   mongoURI,
		CorsOrigin: corsOrigin,
	}
	return nil
}
