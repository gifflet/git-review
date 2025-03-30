package ai

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIProvider struct {
	llm llms.LLM
}

func NewOpenAIProvider(token, model string) (Provider, error) {
	llm, err := openai.New(
		openai.WithToken(token),
		openai.WithModel(model),
	)
	if err != nil {
		return nil, err
	}
	return &OpenAIProvider{llm: llm}, nil
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, p.llm, prompt)
}
