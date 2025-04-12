package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DCCXXV/sapentia.chat/backend/internal/gemini"
	"github.com/DCCXXV/sapentia.chat/backend/internal/utils"
)

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

type ChatHandler struct {
	geminiClient *gemini.Client
}

func NewChatHandler(gc *gemini.Client) *ChatHandler {
	if gc == nil {
		log.Fatal("Gemini client cannot be nil in ChatHandler")
	}
	return &ChatHandler{
		geminiClient: gc,
	}
}

func (h *ChatHandler) HandleChatMessage(c echo.Context) error {
	req := new(ChatRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}
	if req.Message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Message cannot be empty")
	}

	sanitizedUserMessage := utils.SanitizeInput(req.Message)

	ctx := c.Request().Context()
	aiReply, err := h.geminiClient.GenerateContent(ctx, sanitizedUserMessage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get response from AI")
	}

	res := &ChatResponse{
		Reply: aiReply,
	}

	return c.JSON(http.StatusOK, res)
}
