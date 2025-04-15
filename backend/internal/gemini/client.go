package gemini

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

var (
	ErrBlockedResponse = errors.New("gemini: response blocked by safety settings")
	ErrEmptyResponse   = errors.New("gemini: received empty or non-text response")
)

type Client struct {
	genaiClient *genai.Client
	logger      echo.Logger
}

func NewClient(ctx context.Context, apiKey string, logger echo.Logger) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logger.Errorf("Error creating GenAI client: %v", err)
		return nil, fmt.Errorf("error creating GenAI client: %w", err)
	}
	logger.Info("Gemini client initialized successfully (ready to use any model).")

	return &Client{
		genaiClient: client,
		logger:      logger,
	}, nil
}

func (c *Client) GenerateContent(ctx context.Context, modelName string, prompt string) (string, error) {
	if modelName == "" {
		c.logger.Error("Error: GenerateContent called with empty modelName.")
		return "", fmt.Errorf("modelName cannot be empty when calling GenerateContent")
	}

	model := c.genaiClient.GenerativeModel(modelName)

	c.logger.Infof("Sending to Gemini (model: %s): %s", modelName, prompt)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		c.logger.Errorf("Error generating content from Gemini (model: %s): %v", modelName, err)
		return "", fmt.Errorf("failed to get response from AI (model: %s): %w", modelName, err)
	}

	aiReplyContent, err := c.extractTextFromResponse(resp, modelName)
	if err != nil {
		return "", fmt.Errorf("failed to process response from AI (model: %s): %w", modelName, err)
	}

	c.logger.Infof("Received from Gemini (model: %s): %s", modelName, aiReplyContent)
	return aiReplyContent, nil
}

func (c *Client) Close() error {
	c.logger.Info("Closing Gemini client.")
	return c.genaiClient.Close()
}

func (c *Client) extractTextFromResponse(resp *genai.GenerateContentResponse, modelName string) (string, error) {
	if resp == nil {
		c.logger.Warnf("Gemini response object was nil (model: %s)", modelName)
		return "", ErrEmptyResponse
	}

	if resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != genai.BlockReasonUnspecified {
		reason := resp.PromptFeedback.BlockReason.String()
		c.logger.Warnf("Prompt blocked by Gemini (model: %s). Reason: %s", modelName, reason)
		return "", fmt.Errorf("%w: prompt block reason %s", ErrBlockedResponse, reason)
	}

	if len(resp.Candidates) == 0 {
		c.logger.Warnf("Gemini response (model: %s) contained no candidates.", modelName)
		return "", ErrEmptyResponse
	}

	candidate := resp.Candidates[0]
	finishReason := candidate.FinishReason

	if finishReason == genai.FinishReasonSafety {
		reasonStr := finishReason.String()
		c.logger.Warnf("Content generation stopped by safety settings (model: %s). FinishReason: %s", modelName, reasonStr)
		if len(candidate.SafetyRatings) > 0 {
			c.logger.Debugf("Safety Ratings (model: %s): %+v", modelName, candidate.SafetyRatings)
		}
		return "", fmt.Errorf("%w: finish reason %s", ErrBlockedResponse, reasonStr)
	}

	if finishReason != genai.FinishReasonStop && finishReason != genai.FinishReasonUnspecified {
		reasonStr := finishReason.String()
		c.logger.Warnf("Content generation stopped unexpectedly (model: %s). FinishReason: %s", modelName, reasonStr)
		return "", fmt.Errorf("%w: unexpected finish reason %s", ErrEmptyResponse, reasonStr)
	}

	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		c.logger.Warnf("Candidate content or parts are nil/empty (model: %s, finishReason: %s).", modelName, finishReason)
		return "", ErrEmptyResponse
	}

	var sb strings.Builder
	for _, part := range candidate.Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			sb.WriteString(string(textPart))
		} else {
			c.logger.Warnf("Encountered non-text part in response (model: %s): %T", modelName, part)
		}
	}

	fullText := sb.String()
	if fullText == "" {
		c.logger.Warnf("Extracted text is empty after processing parts (model: %s).", modelName)
		return "", ErrEmptyResponse
	}

	return fullText, nil
}
