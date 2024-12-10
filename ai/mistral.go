package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MistralService struct {
	apiKey string
	model  string
	client *http.Client
}

type mistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type mistralRequest struct {
	Model       string           `json:"model"`
	Messages    []mistralMessage `json:"messages"`
	Temperature float64          `json:"temperature,omitempty"`
	TopP        float64          `json:"top_p,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type mistralResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewMistralService(apiKey string, model string) *MistralService {
	if model == "" {
		model = "mistral-tiny"
	}
	return &MistralService{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *MistralService) makeRequest(ctx context.Context, messages []mistralMessage) (string, error) {
	reqBody := mistralRequest{
		Model:       s.model,
		Messages:    messages,
		Temperature: 0.7,
		TopP:        0.9,
		MaxTokens:   1000,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.mistral.ai/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var mistralResp mistralResponse
	if err := json.Unmarshal(body, &mistralResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if len(mistralResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return mistralResp.Choices[0].Message.Content, nil
}

func (s *MistralService) ExplainFFmpegCommand(command string) (string, error) {
	ctx := context.Background()
	messages := []mistralMessage{
		{
			Role:    "system",
			Content: "You are an expert in ffmpeg commands. Explain commands clearly and technically.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Explain what this ffmpeg command does in detail: %s", command),
		},
	}

	return s.makeRequest(ctx, messages)
}

func (s *MistralService) GenerateFFmpegCommand(prompt string) (string, error) {
	ctx := context.Background()
	messages := []mistralMessage{
		{
			Role:    "system",
			Content: "You are an expert in ffmpeg. Generate only the raw command. No explanations, no quotes, no code blocks, no extra text.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return s.makeRequest(ctx, messages)
}
