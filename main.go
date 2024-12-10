package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ffmpeg-ai/ai"
	"github.com/BurntSushi/toml"
)

type Config struct {
	OpenAI struct {
		APIKey string `toml:"api_key"`
		Model  string `toml:"model"`
	} `toml:"openai"`
	Anthropic struct {
		APIKey string `toml:"api_key"`
		Model  string `toml:"model"`
	} `toml:"anthropic"`
	Ollama struct {
		Endpoint string `toml:"endpoint"`
		Model    string `toml:"model"`
	} `toml:"ollama"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ffmpeg-ai \"your request here\"")
		os.Exit(1)
	}

	// Get config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}
	configDir := filepath.Join(homeDir, ".config", "aimpeg")
	configFile := filepath.Join(configDir, "config.toml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Error creating config directory:", err)
	}

	// Load config file
	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal("Error loading config file:", err)
	}

	// Get service type from args or default to OpenAI
	serviceType := "openai"
	if len(os.Args) > 2 {
		serviceType = os.Args[2]
	}

	var aiService ai.Service
	
	switch serviceType {
	case "anthropic":
		if config.Anthropic.APIKey == "" {
			log.Fatal("Anthropic API key not found in config")
		}
		aiService = ai.NewAnthropicService(config.Anthropic.APIKey)
	case "openai":
		if config.OpenAI.APIKey == "" {
			log.Fatal("OpenAI API key not found in config")
		}
		aiService = ai.NewOpenAIService(config.OpenAI.APIKey)
	case "ollama":
		aiService = ai.NewOllamaService(config.Ollama.Model)
	default:
		log.Fatalf("Unknown service type: %s", serviceType)
	}

	// Get ffmpeg command
	command, err := aiService.GenerateFFmpegCommand(os.Args[1])
	if err != nil {
		log.Fatalf("Error generating command: %v", err)
	}

	// Print the ffmpeg command
	fmt.Println(command)
}
