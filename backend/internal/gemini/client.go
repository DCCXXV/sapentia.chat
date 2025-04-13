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
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating GenAI client: %w", err)
	}
	log.Println("Gemini client initialized successfully (ready to use any model).")

	return &Client{
		genaiClient: client,
	}, nil
}

func (c *Client) GenerateContent(ctx context.Context, modelName string, prompt string) (string, error) {
	if modelName == "" {
		log.Println("Error: GenerateContent called with empty modelName.")
		return "", fmt.Errorf("modelName cannot be empty when calling GenerateContent")
	}

	model := c.genaiClient.GenerativeModel(modelName)

	log.Printf("Sending to Gemini (model: %s): %s", modelName, prompt)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Printf("Error generating content from Gemini (model: %s): %v", modelName, err)
		return "", fmt.Errorf("failed to get response from AI (model: %s): %w", modelName, err)
	}

	aiReplyContent := extractTextFromResponse(resp)
	if aiReplyContent == "" {
		log.Printf("Gemini response (model: %s) was empty or had unexpected format: %+v", modelName, resp)
		blockReason := "unknown or no content"
		if resp != nil && resp.PromptFeedback != nil && resp.PromptFeedback.BlockReason != genai.BlockReasonUnspecified {
			blockReason = fmt.Sprintf("PROMPT_BLOCK_%s", resp.PromptFeedback.BlockReason.String())
		} else if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].FinishReason == genai.FinishReasonSafety {
			blockReason = fmt.Sprintf("FINISH_REASON_%s", resp.Candidates[0].FinishReason.String())
		}
		return "", fmt.Errorf("AI (model: %s) returned an empty or blocked response (reason: %s)", modelName, blockReason)
	}

	log.Printf("Received from Gemini (model: %s): %s", modelName, aiReplyContent)
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
	if candidate.FinishReason != genai.FinishReasonStop && candidate.FinishReason != genai.FinishReasonUnspecified {
		log.Printf("Content generation stopped unexpectedly. FinishReason: %s", candidate.FinishReason.String())
		if candidate.FinishReason == genai.FinishReasonSafety && len(candidate.SafetyRatings) > 0 {
			log.Printf("Safety Ratings: %+v", candidate.SafetyRatings)
		}
		return ""
	}

	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		log.Println("Candidate content or parts are nil/empty.")
		return ""
	}

	var fullText string
	for _, part := range candidate.Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			fullText += string(textPart)
		} else {
			log.Printf("Encountered non-text part in response: %T", part)
		}
	}

	if fullText == "" {
		log.Println("Extracted text is empty after processing parts.")
	}

	return fullText
}
