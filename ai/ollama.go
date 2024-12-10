package ai

import (
	"context"
	"fmt"

	"github.com/jmorganca/ollama-go"
)

type OllamaService struct {
	client *ollama.Client
	model  string
}

func NewOllamaService(model string) *OllamaService {
	return &OllamaService{
		client: ollama.NewClient("http://localhost:11434"),
		model:  model,
	}
}

func (s *OllamaService) GenerateFFmpegCommand(prompt string) (string, error) {
	aiPrompt := fmt.Sprintf("Generate only the ffmpeg command for the following request: %s. "+
		"Respond only with the command, no explanations.", prompt)

	resp, err := s.client.Generate(context.Background(), &ollama.GenerateRequest{
		Model:   s.model,
		Prompt:  aiPrompt,
		Stream:  false,
	})

	if err != nil {
		return "", fmt.Errorf("error getting completion: %w", err)
	}

	return resp.Response, nil
}
