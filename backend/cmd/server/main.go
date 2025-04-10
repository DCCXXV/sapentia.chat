package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:5173"},
        AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))

    e.POST("/api/chat", handleChatMessage)

    e.Logger.Fatal(e.Start(":8080"))
}

type ChatRequest struct {
    Message string `json:"message"`
}

type ChatResponse struct {
    Reply string `json:"reply"`
}

func handleChatMessage(c echo.Context) error {
    req := new(ChatRequest)
    if err := c.Bind(req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }


    aiReply := "Esta es una respuesta simulada de la IA"

    res := &ChatResponse{
        Reply : aiReply,
    }

    return c.JSON(http.StatusOK, res)
}
