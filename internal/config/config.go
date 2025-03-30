package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gifflet/git-review/internal/ai"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	AI struct {
		Provider string `mapstructure:"provider"`
		OpenAI   struct {
			Token string `mapstructure:"token"`
			Model string `mapstructure:"model" default:"gpt-4o"`
		} `mapstructure:"openai"`
	} `mapstructure:"ai"`
}

// GetAIProvider creates a new AI provider based on configuration
func (c *Config) GetAIProvider() (ai.Provider, error) {
	cfg := c.AI
	return ai.NewProvider(cfg.Provider, cfg.OpenAI.Token, cfg.OpenAI.Model)
}

// LoadConfig reads configuration from files and environment variables
func LoadConfig(projectPath string) (*Config, error) {
	v := viper.New()

	// Set config name and type
	v.SetConfigName("config")
	v.SetConfigType("yaml")

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
