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
	client := anthropic.NewClient(apiKey)
	return &AnthropicService{
		client: client,
	}
}

func (s *AnthropicService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the command, no explanations.", prompt)

	msg := anthropic.NewMessage(aiPrompt)
	msg.Model = anthropic.Claude3Opus20240229
	msg.MaxTokens = 1000

	resp, err := s.client.Messages.Create(context.Background(), msg)
	if err != nil {
		return "", fmt.Errorf("error getting completion: %w", err)
	}

	return resp.Content[0].Text, nil
}
