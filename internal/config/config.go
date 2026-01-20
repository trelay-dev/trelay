package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	loadEnvFile(".env")
}

// loadEnvFile reads a .env file and sets environment variables.
func loadEnvFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return // .env file is optional
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already set (env vars take precedence)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	App      AppConfig
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Driver   string
	Path     string
	URL      string
	MaxConns int
}

// AuthConfig holds authentication settings.
type AuthConfig struct {
	APIKey       string
	JWTSecret    string
	TokenExpiry  time.Duration
}

// AppConfig holds application-specific settings.
type AppConfig struct {
	BaseURL           string
	DefaultDomain     string
	CustomDomains     []string
	AnalyticsEnabled  bool
	IPAnonymization   bool
	SlugLength        int
	MaxURLLength      int
	RateLimitPerMin   int
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:     getEnvDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ShutdownTimeout: getEnvDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Driver:   getEnv("DB_DRIVER", "sqlite3"),
			Path:     getEnv("DB_PATH", "trelay.db"),
			URL:      getEnv("DB_URL", ""),
			MaxConns: getEnvInt("DB_MAX_CONNS", 10),
		},
		Auth: AuthConfig{
			APIKey:      getEnv("API_KEY", ""),
			JWTSecret:   getEnv("JWT_SECRET", ""),
			TokenExpiry: getEnvDuration("TOKEN_EXPIRY", 24*time.Hour),
		},
		App: AppConfig{
			BaseURL:          getEnv("BASE_URL", "http://localhost:8080"),
			DefaultDomain:    getEnv("DEFAULT_DOMAIN", ""),
			CustomDomains:    getEnvList("CUSTOM_DOMAINS", nil),
			AnalyticsEnabled: getEnvBool("ANALYTICS_ENABLED", true),
			IPAnonymization:  getEnvBool("IP_ANONYMIZATION", true),
			SlugLength:       getEnvInt("SLUG_LENGTH", 6),
			MaxURLLength:     getEnvInt("MAX_URL_LENGTH", 2048),
			RateLimitPerMin:  getEnvInt("RATE_LIMIT_PER_MIN", 100),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required configuration is present.
func (c *Config) Validate() error {
	if c.Auth.APIKey == "" {
		return fmt.Errorf("API_KEY is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}
	if c.App.SlugLength < 4 || c.App.SlugLength > 32 {
		return fmt.Errorf("SLUG_LENGTH must be between 4 and 32")
	}
	return nil
}

// DSN returns the database connection string.
func (c *Config) DSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	if c.Database.Driver == "sqlite3" {
		return c.Database.Path
	}
	return c.Database.URL
}

// Address returns the server listen address.
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvList(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		items := strings.Split(value, ",")
		result := make([]string, 0, len(items))
		for _, item := range items {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}
