package domain

// AppConfig holds application-wide configuration stored in database.
type AppConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ConfigKey constants for application settings.
const (
	ConfigKeyAPIKeyHash     = "api_key_hash"
	ConfigKeyDefaultDomain  = "default_domain"
	ConfigKeyAnalyticsEnabled = "analytics_enabled"
	ConfigKeyIPAnonymization  = "ip_anonymization"
)
