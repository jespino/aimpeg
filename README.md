# AImpeg

AImpeg is a command-line tool that uses AI to generate and explain FFmpeg commands based on natural language descriptions.

## Features

- Generates FFmpeg commands from natural language descriptions
- Supports multiple AI providers:
  - OpenAI (GPT-3.5)
  - Anthropic (Claude)
  - Ollama (local models)
  - Mistral AI
- Interactive command confirmation workflow
- Detailed command explanations
- Configuration via TOML file

## Installation

### Prerequisites

- Go 1.21 or later
- FFmpeg installed on your system

### Install from source

```bash
go install github.com/jespino/aimpeg@latest
```

### Manual installation

1. Clone the repository
```bash
git clone https://github.com/jespino/aimpeg.git
cd aimpeg
```

2. Build and install
```bash
go install
```

## Configuration

Create a configuration file at `~/.config/aimpeg` with your AI provider credentials:

```toml
[openai]
api_key = "your-api-key"
model = "gpt-3.5-turbo"

[anthropic]
api_key = "your-api-key"
model = "claude-3-opus-20240229"

[ollama]
endpoint = "http://localhost:11434"
model = "llama2"

[mistral]
api_key = "your-api-key"
model = "mistral-tiny"
```

At least one provider must be configured.

## Usage

```bash
aimpeg "your request here"
```

Example:
```bash
aimpeg "convert video.mp4 to 720p and add subtitles from subs.srt"
```

The tool will:
1. Generate an FFmpeg command
2. Show you the command
3. Ask if you want to:
   - Run the command
   - Get an explanation
   - Cancel
4. If you choose "explain", it will:
   - Show a detailed explanation
   - Ask if you want to run the command

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
