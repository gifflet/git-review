package types

import "context"

// Provider defines the interface for AI providers
type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// AIConfig holds the configuration for AI providers
type AIConfig interface {
	GetOpenAIToken() string
	GetOpenAIModel() string
	GetSystemPrompt() string
}
