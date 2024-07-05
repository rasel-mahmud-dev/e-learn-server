package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	API_PORT     string
	DATABASE_URI string
}

var Cfg *Config

// LoadConfig loads configuration and stores it in the package-level variable
func LoadConfig() {
	err := godotenv.Load() // Load environment variables from .env file
	if err != nil {
		fmt.Errorf("error loading .env file: %w", err)
	}

	Cfg = &Config{
		API_PORT:     os.Getenv("API_PORT"),
		DATABASE_URI: os.Getenv("DATABASE_URI"),
	}

}
