package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ApiPort       int
	SecretKey     string
	MongoURL      string
	MongoUser     string
	MongoPassword string
	CorsOrigin    string
}

var Config = ServerConfig{}

func Load() error {
	if os.Getenv("E2E") == "true" {
		godotenv.Load(".env.e2e")
	} else {
		godotenv.Load(".env")
	}

	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil || apiPort == 0 {
		return fmt.Errorf("invalid API_PORT env variable, got %v", os.Getenv("API_PORT"))
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("invalid SECRET_KEY env variable")
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return fmt.Errorf("missing MONGO_URL env variable")
	}

	mongoUser := os.Getenv("MONGO_USER")
	if mongoUser == "" {
		return fmt.Errorf("missing MONGO_USER env variable")
	}

	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoPassword == "" {
		return fmt.Errorf("missing MONGO_PASSWORD env variable")
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		return fmt.Errorf("missing CORS_ORIGIN env variable")
	}

	Config = ServerConfig{
		ApiPort:       apiPort,
		SecretKey:     secretKey,
		MongoURL:      mongoURL,
		MongoUser:     mongoUser,
		MongoPassword: mongoPassword,
		CorsOrigin:    corsOrigin,
	}
	return nil
}
