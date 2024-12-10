package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OllamaService struct {
	client *openai.Client
	model  string
}

func NewOllamaService(model string) *OllamaService {
	return &OllamaService{
		client: openai.NewClientWithConfig(openai.DefaultConfig("")),
		model:  model,
	}
}

func (s *OllamaService) ExplainFFmpegCommand(command string) (string, error) {
	aiPrompt := fmt.Sprintf("Explain what this ffmpeg command does in detail: %s", command)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: aiPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error getting explanation: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *OllamaService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the command, no explanations.", prompt)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: aiPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("error getting completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
