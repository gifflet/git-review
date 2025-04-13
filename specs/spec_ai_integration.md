# Git Review AI Integration
> Ingest the information from this file, implement the Low-Level Tasks, and generate the code that will satisfy the High and Mid-Level Objectives.

## High-Level Objective

- Implement the code review using AI for each file present in a range of commits 

## Mid-Level Objective

- The model, AI provider and API Key should be configured in the yaml file
- The review should be created for each changed file (in git diff format) in a range of commits 
- The created reviews should have only middle-high proposal changes

## Implementation Notes
- Use the package `github.com/tmc/langchaingo` as this have all LLMs clients supported
- Load the AI config from the yaml file
- Each provider should have its own interface implementation
- Implement tests for the complate new funcionality coverage

## Context

### Beginning context
- internal/ai/ai.go
- internal/ai/ai_openai.go
- internal/ai/ai_review.go
- internal/ai/ai_review_test.go
- internal/config/config.go
- cmd/root.go
- cmd/review.go
- internal/model/model.go

### Ending context  
- internal/ai/ai.go
- internal/ai/ai_openai.go
- internal/ai/ai_review.go
- internal/ai/ai_review_test.go
- internal/config/config.go
- cmd/root.go
- cmd/review.go
- internal/model/model.go

## Low-Level Tasks
> Ordered from start to finish

1. Add the new module `github.com/tmc/langchaingo`
```aider
UPDATE go.mod:
    ADD github.com/tmc/langchaingo
```

2. Create the ai service
```aider
CREATE internal/ai/ai.go:
    CREATE AiCall(modelName: string, apiKey: string, systemMessage: string, userMessage: string) -> string
        This function should be a abstract function in a way that each AI provider should follow its signature
```

3. Create the OpenAI service AI implemnentation from interface
```aider
CREATE internal/ai/ai_openai.go:
    CREATE AiCall(modelName: string, apiKey: string, systemMessage: string, userMessage: string) -> string

    Interface implementation. Use as reference the snippet bellow:

        import (
            "context"
            "fmt"
            "log"

            "github.com/tmc/langchaingo/llms"
            "github.com/tmc/langchaingo/schema"
            "github.com/tmc/langchaingo/llms/openai"
        )

        func main() {
            llm, err := openai.New(openai.WithModel("gpt-4o"))
            if err != nil {
                log.Fatal(err)
            }
            ctx := context.Background()

            content := []llms.MessageContent{
                llms.TextParts(schema.ChatMessageTypeSystem, "You are a go programming expert"),
                llms.TextParts(schema.ChatMessageTypeHuman, "Why is go a good language for production LLM applications?"),
            }

            completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
                fmt.Print(string(chunk))
                return nil
            }))
            if err != nil {
                log.Fatal(err)
            }
            _ = completion
        }
```

4. Create the model package
```aider
CREATE internal/model/model.go:
    CREATE Review struct {
        comment string,
        filePath string,
        linePosition int
    }
    CREATE ReviewList []Review
```

5. Create the review function
```aider
CREATE internal/ai/ai_review.go:
    CREATE Review(aiProvider: string, modelName: string, apiKey: string, filePath: string,  fileDiff: string) -> ReviewList
        This function should use the AiCall implementation depending on the AI provider (aiProvider - openai, anthropic, etc). 
        The userMessage should be the fileDiff as it is the diff of the file in the git diff format.
        The review should be a list of comments in the format of the Review struct.
        The system message should instruct the model to review the file and provide feedback in the format of the ReviewList struct.
```

6. Create the review test function
```aider
CREATE internal/ai/ai_review_test.go:
    CREATE ReviewTest(aiProvider: string, modelName: string, apiKey: string, filePath: string,  fileDiff: string) -> ReviewList
        This function should be a test function that should test the Review function.
```

7. Use the review function in the review command
```aider
UPDATE cmd/review.go:
    UPDATE executeReview()
        This function should use the Review function to create the review for the file.
        The review should be created for each changed file (in git diff format) in a range of commits
```
