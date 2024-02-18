package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ApiPort     int
	SecretKey   string
	DisableAuth bool
	MongoURL    string
	MongoPort   int
}

var Config = ServerConfig{}

func Load() error {
	godotenv.Load(".env")

	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil || apiPort == 0 {
		return fmt.Errorf("invalid API_PORT env variable, got %v", os.Getenv("API_PORT"))
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("invalid SECRET_KEY env variable")
	}

	disableAuth := os.Getenv("DISABLE_AUTH")
	if disableAuth != "true" && disableAuth != "false" {
		return fmt.Errorf("invalid DISABLE_AUTH env variable, expected 'true' or 'false', got %v", os.Getenv("DISABLE_AUTH"))
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return fmt.Errorf("missing MONGO_URL env variable")
	}

	mongoPort, err := strconv.Atoi(os.Getenv("MONGO_PORT"))
	if err != nil || mongoPort == 0 {
		return fmt.Errorf("invalid MONGO_PORT env variable, got %v", os.Getenv("MONGO_PORT"))
	}

	Config = ServerConfig{
		ApiPort:     apiPort,
		SecretKey:   secretKey,
		DisableAuth: disableAuth == "true",
		MongoURL:    mongoURL,
		MongoPort:   mongoPort,
	}
	return nil
}
