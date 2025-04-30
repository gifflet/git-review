package ai

import (
	"context"
	"testing"

	"github.com/gifflet/git-review/internal/types"
	"github.com/stretchr/testify/assert"
)

// MockProvider implements types.Provider for testing
type MockProvider struct {
	response string
	err      error
}

func (m *MockProvider) Generate(ctx context.Context, prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func TestReview(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		fileDiff       string
		mockResponse   string
		expectedError  error
		expectedReview types.ReviewList
	}{
		{
			name:     "successful JSON response",
			filePath: "test.go",
			fileDiff: "test diff content",
			mockResponse: `[
				{
					"Comment": "Consider using a more descriptive variable name",
					"FilePath": "test.go",
					"LinePosition": 10
				}
			]`,
			expectedError: nil,
			expectedReview: types.ReviewList{
				{
					Comment:      "Consider using a more descriptive variable name",
					FilePath:     "test.go",
					LinePosition: 10,
				},
			},
		},
		{
			name:     "successful text response",
			filePath: "test.go",
			fileDiff: "test diff content",
			mockResponse: `Line 10:
Consider using a more descriptive variable name`,
			expectedError: nil,
			expectedReview: types.ReviewList{
				{
					Comment:      "Consider using a more descriptive variable name",
					FilePath:     "test.go",
					LinePosition: 10,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &MockProvider{
				response: tt.mockResponse,
				err:      tt.expectedError,
			}

			reviews, err := Review(context.Background(), mockProvider, tt.filePath, tt.fileDiff)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedReview, reviews)
			}
		})
	}
}

func TestParseTextResponse(t *testing.T) {
	tests := []struct {
		name           string
		response       string
		filePath       string
		expectedReview types.ReviewList
	}{
		{
			name: "single review",
			response: `Line 10:
Consider using a more descriptive variable name`,
			filePath: "test.go",
			expectedReview: types.ReviewList{
				{
					Comment:      "Consider using a more descriptive variable name",
					FilePath:     "test.go",
					LinePosition: 10,
				},
			},
		},
		{
			name: "multiple reviews",
			response: `Line 10:
Consider using a more descriptive variable name

Line 20:
This function could be simplified`,
			filePath: "test.go",
			expectedReview: types.ReviewList{
				{
					Comment:      "Consider using a more descriptive variable name",
					FilePath:     "test.go",
					LinePosition: 10,
				},
				{
					Comment:      "This function could be simplified",
					FilePath:     "test.go",
					LinePosition: 20,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviews := parseTextResponse(tt.response, tt.filePath)
			assert.Equal(t, tt.expectedReview, reviews)
		})
	}
}
