package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	GeminiAPIKey string
	AllowOrigins []string
	ServerPort   string
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables directly")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowedOrigins := []string{"http://localhost:5173"}

	return &AppConfig{
		GeminiAPIKey: apiKey,
		AllowOrigins: allowedOrigins,
		ServerPort:   port,
	}, nil
}
