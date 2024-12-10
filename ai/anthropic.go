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
	cfg := anthropic.ClientConfig{
		APIKey: apiKey,
	}
	client := anthropic.NewClient(cfg)
	return &AnthropicService{
		client: client,
	}
}

func (s *AnthropicService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the command, no explanations.", prompt)

	resp, err := s.client.Messages.Create(context.Background(), &anthropic.MessageRequest{
		Model:    "claude-3-opus-20240229",
		MaxTokens: anthropic.Int(1000),
		System:   anthropic.String("You are an expert in ffmpeg commands. Respond only with the command, no explanations."),
		Messages: []anthropic.MessageParam{
			{
				Role:    "user",
				Content: []anthropic.MessageContent{
					{
						Type: "text",
						Text: aiPrompt,
					},
				},
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("error getting completion: %w", err)
	}

	return resp.Content[0].Text, nil
}
