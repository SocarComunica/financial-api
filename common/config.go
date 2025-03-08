package common

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Database    string
}

var config *Config

func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Unable to load .env file, will be using system environment variables")
	}

	config = &Config{
		Environment: GetEnv("ENV", "development"),
		Database:    GetEnv("DATABASE", "local.db"),
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}

func GetConfig() (*Config, error) {
	if config == nil {
		return nil, errors.New("config not instantiated")
	}

	return config, nil
}
