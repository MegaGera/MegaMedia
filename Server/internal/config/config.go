package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	MegagoalServerApiUrl string
	MongoDBURI           string
	ValidateUri          string
}

var Cfg Config // Global config instance

func LoadConfig() {
	// Determine which .env file to load
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // Default to development
	}

	envFile := ".env.development"
	if env == "production" {
		envFile = ".env.production"
	}

	// Load environment variables
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}

	// Populate config struct
	Cfg = Config{
		AppEnv:               os.Getenv("APP_ENV"),
		MongoDBURI:           os.Getenv("MONGODB_URI"),
		MegagoalServerApiUrl: os.Getenv("MEGAGOAL_SERVER_API_URL"),
		ValidateUri:          os.Getenv("VALIDATE_URI"),
	}

	log.Printf("Loaded configuration for %s", env)
}
