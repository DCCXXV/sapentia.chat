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
	"github.com/ziflex/lecho/v3"

	"github.com/DCCXXV/sapentia.chat/backend/internal/api"
	"github.com/DCCXXV/sapentia.chat/backend/internal/config"
	"github.com/DCCXXV/sapentia.chat/backend/internal/gemini"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	e := echo.New()

	lechoLogger := lecho.New(os.Stdout)

	e.Logger = lechoLogger
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	ctx := context.Background()
	geminiClient, err := gemini.NewClient(ctx, cfg.GeminiAPIKey, lechoLogger)
	if err != nil {
		lechoLogger.Fatalf("Failed 2to initialize Gemini client: %v", err)
	}

	defer func() {
		if err := geminiClient.Close(); err != nil {
			lechoLogger.Printf("Error closing gemini client: %v", err)
		}
	}()

	chatHandler, err := api.NewChatHandler(geminiClient)

	if err != nil {
		lechoLogger.Error("Failed to create chat handler")
	}

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
		e.Logger.Info("Starting server on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server unexpectedly: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Server forced to shutdown: ", err)
	}

	e.Logger.Info("Server exiting gracefully")
}
