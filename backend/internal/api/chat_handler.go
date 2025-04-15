package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DCCXXV/sapentia.chat/backend/internal/gemini"
	"github.com/DCCXXV/sapentia.chat/backend/internal/utils"
)

type ChatRequest struct {
	Message          string `json:"message"`
	SelectedModelID  string `json:"selectedModelId"`
	AssistedLearning bool   `json:"assistedLearning"`
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
	"2": "gemini-2.5-pro-exp-03-25",
}

const defaultModelName = "gemini-2.0-flash"

func NewChatHandler(gc *gemini.Client) (*ChatHandler, error) {
	if gc == nil {
		return nil, fmt.Errorf("gemini client cannot be nil in ChatHandler")
	}
	return &ChatHandler{
		geminiClient: gc,
	}, nil
}

func (h *ChatHandler) HandleChatMessage(c echo.Context) error {
	req := new(ChatRequest)
	if err := c.Bind(req); err != nil {
		c.Logger().Errorf("Error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	c.Logger().Debugf("Received chat request: Message='%s', ModelID='%s', AssistedLearning=%t", req.Message, req.SelectedModelID, req.AssistedLearning)

	if req.Message == "" {
		c.Logger().Warn("Received request with empty message")
		return echo.NewHTTPError(http.StatusBadRequest, "Message cannot be empty")
	}

	var targetModelName string
	if req.SelectedModelID == "" {
		c.Logger().Warn("Warning: SelectedModelID is empty in request, using default model.")
		targetModelName = defaultModelName
	} else {
		modelName, found := modelIDToNameMap[req.SelectedModelID]
		if !found {
			c.Logger().Warn("Warning: Received unknown model ID '%s', using default model.", req.SelectedModelID)
			targetModelName = defaultModelName
		} else {
			targetModelName = modelName
		}
	}

	c.Logger().Info("Processing chat message using model: %s", targetModelName)

	ctx := c.Request().Context()
	var finalMessage string

	if !req.AssistedLearning {
		finalMessage = req.Message
	} else {
		finalMessage = utils.AssistedLearningPrompt + req.Message
	}

	aiReply, err := h.geminiClient.GenerateContent(ctx, targetModelName, finalMessage)
	if err != nil {
		c.Logger().Errorf("Error calling Gemini API with model %s: %v", targetModelName, err)
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get response from AI using model %s", targetModelName))
	}

	res := &ChatResponse{
		Reply: aiReply,
	}

	return c.JSON(http.StatusOK, res)
}
