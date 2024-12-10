package main

import (
	"aimpeg/ai"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
		fmt.Println("Usage: aimpeg your request here")
		os.Exit(1)
	}

	// Concatenate all arguments after the program name into a single string
	prompt := strings.Join(os.Args[1:], " ")

	// Get config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting home directory:", err)
	}
	configDir := filepath.Join(homeDir, ".config")
	configFile := filepath.Join(configDir, "aimpeg")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Error creating config directory:", err)
	}

	// Load config file
	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal("Error loading config file:", err)
	}

	var aiService ai.Service

	// Try services in order until we find one that's configured
	if config.OpenAI.APIKey != "" {
		aiService = ai.NewOpenAIService(config.OpenAI.APIKey)
	} else if config.Anthropic.APIKey != "" {
		aiService = ai.NewAnthropicService(config.Anthropic.APIKey)
	} else if config.Ollama.Endpoint != "" {
		aiService = ai.NewOllamaService(config.Ollama.Model)
	} else {
		log.Fatal("No AI service is properly configured. Please check your config file.")
	}

	// Get ffmpeg command
	command, err := aiService.GenerateFFmpegCommand(prompt)
	if err != nil {
		log.Fatalf("Error generating command: %v", err)
	}

	// Print the generated command
	fmt.Println("\nGenerated command:", command)

	for {
		fmt.Print("\nDo you want to run this command? (yes/no/explain): ")
		var response string
		fmt.Scanln(&response)

		switch strings.ToLower(response) {
		case "yes":
			// Execute the command
			cmdParts := strings.Fields(command)
			if len(cmdParts) == 0 {
				log.Fatal("Empty command received from AI")
			}

			cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				log.Fatalf("Error executing command: %v", err)
			}
			return

		case "no":
			fmt.Println("Command execution cancelled")
			return

		case "explain":
			explanation, err := aiService.ExplainFFmpegCommand(command)
			if err != nil {
				log.Fatalf("Error getting explanation: %v", err)
			}
			fmt.Printf("\nExplanation:\n%s\n", explanation)
			fmt.Print("\nDo you want to run this command? (yes/no): ")
			fmt.Scanln(&response)
			if strings.ToLower(response) == "yes" {
				cmdParts := strings.Fields(command)
				if len(cmdParts) == 0 {
					log.Fatal("Empty command received from AI")
				}

				cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					log.Fatalf("Error executing command: %v", err)
				}
			} else {
				fmt.Println("Command execution cancelled")
			}
			return

		default:
			fmt.Println("Invalid response. Please answer 'yes', 'no', or 'explain'")
		}
	}
}
