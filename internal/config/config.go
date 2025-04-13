package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gifflet/git-review/internal/ai"
	"github.com/gifflet/git-review/internal/types"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	AI struct {
		Provider string `mapstructure:"provider"`
		OpenAI   struct {
			Token string `mapstructure:"token"`
			Model string `mapstructure:"model"`
		} `mapstructure:"openai"`
		SystemPrompt string `mapstructure:"system_prompt"`
	} `mapstructure:"ai"`
}

func (c *Config) GetOpenAIToken() string {
	return c.AI.OpenAI.Token
}

func (c *Config) GetOpenAIModel() string {
	return c.AI.OpenAI.Model
}

func (c *Config) GetSystemPrompt() string {
	if c.AI.SystemPrompt == "" {
		return `You are a code review assistant. Your task is to review code changes and provide constructive feedback.

Please analyze the provided git diff and suggest improvements focusing on:
1. Code quality and best practices
2. Potential bugs or issues
3. Performance considerations
4. Security implications
5. Maintainability and readability

Format your response as a list of reviews, where each review contains:
- The file path
- The line position in the file
- A clear and constructive comment explaining the suggested improvement

Only suggest meaningful, high-impact changes that would significantly improve the code.
Avoid nitpicking or suggesting minor stylistic changes unless they significantly impact readability.`
	}
	return c.AI.SystemPrompt
}

// GetAIProvider creates a new AI provider based on configuration
func (c *Config) GetAIProvider() (types.Provider, error) {
	return ai.NewProvider(c.AI.Provider, c)
}

// LoadConfig reads configuration from files and environment variables
func LoadConfig(projectPath string) (*Config, error) {
	v := viper.New()

	// Set config name and type
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Set default values
	v.SetDefault("ai.provider", "openai")
	v.SetDefault("ai.openai.model", "gpt-4o")
	v.SetDefault("ai.system_prompt", "You are a helpful assistant that reviews code changes. You are given a git diff with the changes made to the code. You need to review the changes and provide a list of potential improvements.")

	// Add global config paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Linux/Mac global config path
	v.AddConfigPath(filepath.Join(homeDir, ".config", "git-review"))

	// Windows global config path
	v.AddConfigPath(filepath.Join(homeDir, "AppData", "Roaming", "git-review"))

	// Local project config
	v.AddConfigPath(projectPath)
	v.SetConfigName(".gitreview")

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("configuration file not found. Please create one of the following files:\n"+
				"- %s\n"+
				"- %s\n"+
				"- %s\n"+
				"\nExample configuration:\n"+
				"ai:\n"+
				"  provider: \"openai\"\n"+
				"  openai:\n"+
				"    token: \"your-token-here\"\n"+
				"    model: \"gpt-4o\"",
				filepath.Join(homeDir, ".config", "git-review", "config.yaml"),
				filepath.Join(homeDir, "AppData", "Roaming", "git-review", "config.yaml"),
				".gitreview.yaml")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Set environment variable prefix
	v.SetEnvPrefix("GIT_REVIEW")
	v.AutomaticEnv()

	// Create config struct
	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
