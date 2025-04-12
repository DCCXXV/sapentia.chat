package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/DCCXXV/sapentia.chat/backend/internal/api"
	"github.com/DCCXXV/sapentia.chat/backend/internal/config"
	"github.com/DCCXXV/sapentia.chat/backend/internal/gemini"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	geminiClient, err := gemini.NewClient(ctx, cfg.GeminiAPIKey, "gemini-1.5-flash-latest")
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}

	defer func() {
		if err := geminiClient.Close(); err != nil {
			log.Printf("Error closing gemini client: %v", err)
		}
	}()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	chatHandler := api.NewChatHandler(geminiClient)
	apiGroup := e.Group("/api")
	{
		apiGroup.POST("/chat", chatHandler.HandleChatMessage)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	startServer(e, cfg.ServerPort)
}

func startServer(e *echo.Echo, port string) {
	go func() {
		serverAddr := fmt.Sprintf(":%s", port)
		log.Printf("Starting server on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server unexpectedly: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting gracefully")
}
