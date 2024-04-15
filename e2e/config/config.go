package config

import (
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	APIBaseURL    string
	MongoURL      string
	MongoUser     string
	MongoPassword string
}

var Config config

func init() {
	godotenv.Load(".env")
	Config = config{
		APIBaseURL:    os.Getenv("API_BASE_URL"),
		MongoURL:      os.Getenv("MONGO_URL"),
		MongoUser:     os.Getenv("MONGO_USER"),
		MongoPassword: os.Getenv("MONGO_PASSWORD"),
	}
}
