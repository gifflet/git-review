package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gifflet/git-review/internal/types"
)

// Review generates a list of code review comments for a given file diff
func Review(ctx context.Context, provider types.Provider, filePath string, fileDiff string) (types.ReviewList, error) {
	// Create the user message with the file diff
	userMessage := fmt.Sprintf("Please review the following changes in file %s:\n\n%s", filePath, fileDiff)

	// Generate the review using the AI provider
	response, err := provider.Generate(ctx, userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to generate review: %w", err)
	}

	// Parse the response into a list of reviews
	var reviews types.ReviewList

	// Try to parse the response as a JSON array first
	if err := json.Unmarshal([]byte(response), &reviews); err != nil {
		// If JSON parsing fails, try to parse the response as a text format
		reviews = parseTextResponse(response, filePath)
	}

	return reviews, nil
}

// parseTextResponse parses a text response into a list of reviews
func parseTextResponse(response string, filePath string) types.ReviewList {
	var reviews types.ReviewList
	lines := strings.Split(response, "\n")

	var currentReview *types.Review
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if line contains a line number
		if strings.Contains(line, "Line") || strings.Contains(line, "line") {
			// If we have a previous review, add it to the list
			if currentReview != nil {
				reviews = append(reviews, *currentReview)
			}

			// Create a new review
			currentReview = &types.Review{
				FilePath: filePath,
			}

			// Try to extract line number using a more robust approach
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				numStr := strings.TrimSpace(strings.ToLower(parts[0]))
				numStr = strings.TrimPrefix(numStr, "line")
				numStr = strings.TrimSpace(numStr)
				if num, err := strconv.Atoi(numStr); err == nil {
					currentReview.LinePosition = num
				}
			}
		} else if currentReview != nil {
			// Append the line to the current review's comment
			if currentReview.Comment == "" {
				currentReview.Comment = line
			} else {
				currentReview.Comment += "\n" + line
			}
		}
	}

	// Add the last review if exists
	if currentReview != nil {
		reviews = append(reviews, *currentReview)
	}

	return reviews
}
