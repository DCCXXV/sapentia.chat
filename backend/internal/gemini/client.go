package gemini

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Client struct {
	genaiClient *genai.Client
	modelName   string
}

func NewClient(ctx context.Context, apiKey string, modelName string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}
	if modelName == "" {
		modelName = "gemini-2.0-flash-lite"
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating GenAI client: %w", err)
	}
	log.Println("Gemini client initialized successfully.")

	return &Client{
		genaiClient: client,
		modelName:   modelName,
	}, nil
}

func (c *Client) GenerateContent(ctx context.Context, prompt string) (string, error) {
	model := c.genaiClient.GenerativeModel(c.modelName)
	log.Printf("Sending to Gemini: %s", prompt)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error generating content from Gemini: %v", err)
		return "", fmt.Errorf("failed to get response from AI: %w", err)
	}

	aiReplyContent := extractTextFromResponse(resp)
	if aiReplyContent == "" {
		log.Printf("Gemini response was empty or had unexpected format: %+v", resp)
		return "", fmt.Errorf("AI returned an empty or blocked response")
	}

	log.Printf("Received from Gemini: %s", aiReplyContent)
	return aiReplyContent, nil
}

func (c *Client) Close() error {
	log.Println("Closing Gemini client.")
	return c.genaiClient.Close()
}

func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 {
		if resp != nil && resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != genai.BlockReasonUnspecified {
			log.Printf("Prompt blocked due to safety settings: Reason %s", resp.PromptFeedback.BlockReason.String())
		}
		return ""
	}

	candidate := resp.Candidates[0]
	if candidate.FinishReason == genai.FinishReasonSafety {
		log.Printf("Content blocked due to safety settings: %+v", candidate.SafetyRatings)
		return ""
	}

	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return ""
	}

	if textPart, ok := candidate.Content.Parts[0].(genai.Text); ok {
		return string(textPart)
	}

	log.Printf("First part of response was not text: %+v", candidate.Content.Parts[0])
	return ""
}
