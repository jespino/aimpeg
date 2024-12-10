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

type OllamaService struct {
	model    string
	endpoint string
	client   *http.Client
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool           `json:"stream"`
	Options  *ollamaOptions `json:"options,omitempty"`
}

type ollamaOptions struct {
	Temperature    float64 `json:"temperature,omitempty"`
	NumPredict    int     `json:"num_predict,omitempty"`
	TopK          int     `json:"top_k,omitempty"`
	TopP          float64 `json:"top_p,omitempty"`
	Seed          int     `json:"seed,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
}

type ollamaChatResponse struct {
	Model    string `json:"model"`
	Message  struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

func NewOllamaService(model string) *OllamaService {
	return &OllamaService{
		model:    model,
		endpoint: "http://localhost:11434",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *OllamaService) makeRequest(ctx context.Context, messages []ollamaMessage) (string, error) {
	reqBody := ollamaChatRequest{
		Model:    s.model,
		Messages: messages,
		Stream:   false,
		Options: &ollamaOptions{
			Temperature: 0.7,
			TopK: 40,
			TopP: 0.9,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint+"/api/chat", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	var ollamaResp ollamaChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error parsing response: %w", err)
	}

	if ollamaResp.Error != "" {
		return "", fmt.Errorf("ollama error: %s", ollamaResp.Error)
	}

	return ollamaResp.Message.Content, nil
}

func (s *OllamaService) ExplainFFmpegCommand(command string) (string, error) {
	ctx := context.Background()
	messages := []ollamaMessage{
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

func (s *OllamaService) GenerateFFmpegCommand(prompt string) (string, error) {
	ctx := context.Background()
	messages := []ollamaMessage{
		{
			Role:    "system",
			Content: "You are an expert in ffmpeg. Generate only the raw command without code block markers or explanations.",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	return s.makeRequest(ctx, messages)
}
