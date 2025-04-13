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

// Review represents a single review comment for a specific file and line
type Review struct {
	Comment      string
	FilePath     string
	LinePosition int
}

// ReviewList represents a list of reviews
type ReviewList []Review
