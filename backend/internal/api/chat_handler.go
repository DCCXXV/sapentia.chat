package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DCCXXV/sapentia.chat/backend/internal/gemini"
)

type ChatRequest struct {
	Message         string `json:"message"`
	SelectedModelID string `json:"selectedModelId"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

type ChatHandler struct {
	geminiClient *gemini.Client
}

var modelIDToNameMap = map[string]string{
	"0": "gemini-2.0-flash-lite",
	"1": "gemini-2.0-flash",
	"2": "gemini-2.5-pro-latest",
}

const defaultModelName = "gemini-2.0-flash"

func NewChatHandler(gc *gemini.Client) *ChatHandler {
	if gc == nil {
		log.Fatalf("Gemini client cannot be nil in ChatHandler")
	}
	return &ChatHandler{
		geminiClient: gc,
	}
}

func (h *ChatHandler) HandleChatMessage(c echo.Context) error {
	req := new(ChatRequest)
	if err := c.Bind(req); err != nil {
		log.Printf("Error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	log.Println(req.Message)
	log.Println(req.SelectedModelID)

	if req.Message == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Message cannot be empty")
	}

	var targetModelName string
	if req.SelectedModelID == "" {
		log.Println("Warning: SelectedModelID is empty in request, using default model.")
		targetModelName = defaultModelName
	} else {
		modelName, found := modelIDToNameMap[req.SelectedModelID]
		if !found {
			log.Printf("Warning: Received unknown model ID '%s', using default model.", req.SelectedModelID)
			targetModelName = defaultModelName
		} else {
			targetModelName = modelName
		}
	}

	log.Printf("Processing chat message using model: %s", targetModelName)

	ctx := c.Request().Context()

	aiReply, err := h.geminiClient.GenerateContent(ctx, targetModelName, req.Message)
	if err != nil {
		log.Printf("Error calling Gemini API with model %s: %v", targetModelName, err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get response from AI using model %s", targetModelName))
	}

	res := &ChatResponse{
		Reply: aiReply,
	}

	return c.JSON(http.StatusOK, res)
}
