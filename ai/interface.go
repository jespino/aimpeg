package ai

// Service defines the interface for AI services that can generate ffmpeg commands
type Service interface {
	GenerateFFmpegCommand(prompt string) (string, error)
}
