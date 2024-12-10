package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

type AnthropicService struct {
	client *anthropic.Client
}

func NewAnthropicService(apiKey string) *AnthropicService {
	client, _ := anthropic.NewClient(apiKey)
	return &AnthropicService{
		client: client,
	}
}

func (s *AnthropicService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the command, no explanations.", prompt)

	resp, err := s.client.CreateMessage(context.Background(), &anthropic.CreateMessageRequest{
		Model:     anthropic.Claude3Opus,
		MaxTokens: 1000,
		Messages: []anthropic.Message{
			{
				Role:    anthropic.UserRole,
				Content: aiPrompt,
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("error getting completion: %w", err)
	}

	return resp.Content[0].Text, nil
}
