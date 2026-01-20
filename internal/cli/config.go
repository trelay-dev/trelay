package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds CLI configuration.
type Config struct {
	APIURL string `json:"api_url"`
	APIKey string `json:"api_key"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		APIURL: "http://localhost:8080",
	}
}

// ConfigPath returns the path to the config file.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "trelay", "config.json"), nil
}

// LoadConfig loads configuration from file.
func LoadConfig() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to file.
func SaveConfig(cfg *Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetClient creates an API client from config.
func GetClient() (*Client, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	// Override with environment variables
	if url := os.Getenv("TRELAY_API_URL"); url != "" {
		cfg.APIURL = url
	}
	if key := os.Getenv("TRELAY_API_KEY"); key != "" {
		cfg.APIKey = key
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API key not configured. Run 'trelay config set api-key <key>' or set TRELAY_API_KEY")
	}

	return NewClient(cfg.APIURL, cfg.APIKey), nil
}
