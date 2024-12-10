package main

import (
	"fmt"
	"log"
	"os"

	"ffmpeg-ai/ai"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ffmpeg-ai \"your request here\"")
		os.Exit(1)
	}

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get service type from args or default to OpenAI
	serviceType := "openai"
	if len(os.Args) > 2 {
		serviceType = os.Args[2]
	}

	var aiService ai.Service
	
	switch serviceType {
	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			log.Fatal("ANTHROPIC_API_KEY not found in environment")
		}
		aiService = ai.NewAnthropicService(apiKey)
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Fatal("OPENAI_API_KEY not found in environment")
		}
		aiService = ai.NewOpenAIService(apiKey)
	case "ollama":
		model := os.Getenv("OLLAMA_MODEL")
		if model == "" {
			model = "llama2" // default model
		}
		aiService = ai.NewOllamaService(model)
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
