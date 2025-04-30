package ai

import (
	"fmt"

	"github.com/gifflet/git-review/internal/types"
)

// NewProvider creates a new AI provider based on configuration
func NewProvider(providerName string, c types.AIConfig) (types.Provider, error) {
	switch providerName {
	case "openai":
		return NewOpenAIProvider(c.GetOpenAIToken(), c.GetOpenAIModel(), c.GetSystemPrompt())
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", providerName)
	}
}
