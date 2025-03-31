package ai

import (
	"context"

	"github.com/gifflet/git-review/internal/types"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIProvider struct {
	llm          llms.LLM
	systemPrompt string
}

func NewOpenAIProvider(token, model, systemPrompt string) (types.Provider, error) {
	llm, err := openai.New(
		openai.WithToken(token),
		openai.WithModel(model),
	)
	if err != nil {
		return nil, err
	}
	return &OpenAIProvider{
		llm:          llm,
		systemPrompt: systemPrompt,
	}, nil
}

func (p *OpenAIProvider) Generate(ctx context.Context, prompt string) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, p.systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}
	completion, err := p.llm.GenerateContent(ctx, content)
	if err != nil {
		return "", err
	}
	return completion.Choices[0].Content, nil
}
