package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIService struct {
	client *openai.Client
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
	}
}

func (s *OpenAIService) ExplainFFmpegCommand(command string) (string, error) {
	aiPrompt := fmt.Sprintf("Explain what this ffmpeg command does in detail: %s", command)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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

func (s *OpenAIService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the raw command, no explanations or code block markers.", prompt)

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
