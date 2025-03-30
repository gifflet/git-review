package ai

import (
	"context"
	"fmt"
)

// Provider defines the interface for AI providers
type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// NewProvider creates a new AI provider based on configuration
func NewProvider(providerName, token, model string) (Provider, error) {
	switch providerName {
	case "openai":
		return NewOpenAIProvider(token, model)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerName)
	}
}
