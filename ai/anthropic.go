package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicService struct {
	client *anthropic.Client
}

func NewAnthropicService(apiKey string) *AnthropicService {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &AnthropicService{
		client: client,
	}
}

func (s *AnthropicService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the command, no explanations.", prompt)

	message, err := s.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		MaxTokens: anthropic.F(int64(1000)),
		System:    anthropic.F("You are an expert in ffmpeg commands. Respond only with the command, no explanations."),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(aiPrompt)),
		}),
	})

	if err != nil {
		return "", fmt.Errorf("error getting completion: %w", err)
	}

	return message.Content[0].Text, nil
}
