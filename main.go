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

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not found in environment")
	}

	// Initialize AI service
	aiService := ai.NewOpenAIService(apiKey)

	// Get ffmpeg command
	command, err := aiService.GenerateFFmpegCommand(os.Args[1])
	if err != nil {
		log.Fatalf("Error generating command: %v", err)
	}

	// Print the ffmpeg command
	fmt.Println(command)
}
