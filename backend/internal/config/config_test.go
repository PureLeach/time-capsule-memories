package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validConfig() Config {
	return Config{
		PostgresURL:         "postgres://user:pass@localhost:5432/db",
		MinioEndpoint:       "localhost:9000",
		MinioAccessKey:      "key",
		MinioSecretKey:      "secret",
		MinioBucketName:     "time-capsule",
		SMTPHost:            "localhost",
		SMTPPort:            "1025",
		SMTPFrom:            "capsule@example.com",
		SMTPTimeout:         10,
		CronCapsuleDispatch: "*/5 * * * *",
		CORSAllowedOrigins:  []string{"http://localhost:8001"},
	}
}

func TestValidate_AcceptsCompleteConfig(t *testing.T) {
	cfg := validConfig()
	require.NoError(t, cfg.Validate())
}

func TestValidate_RejectsMissingRequiredFields(t *testing.T) {
	cases := map[string]func(*Config){
		"DATABASE_URL":          func(c *Config) { c.PostgresURL = "" },
		"MINIO_ENDPOINT":        func(c *Config) { c.MinioEndpoint = "" },
		"MINIO_ROOT_USER":       func(c *Config) { c.MinioAccessKey = "" },
		"MINIO_ROOT_PASSWORD":   func(c *Config) { c.MinioSecretKey = "" },
		"SMTP_HOST":             func(c *Config) { c.SMTPHost = "" },
		"SMTP_PORT":             func(c *Config) { c.SMTPPort = "" },
		"SMTP_FROM":             func(c *Config) { c.SMTPFrom = "" },
		"CRON_CAPSULE_DISPATCH": func(c *Config) { c.CronCapsuleDispatch = "" },
		"SMTP_TIMEOUT":          func(c *Config) { c.SMTPTimeout = 0 },
		"CORS_ALLOWED_ORIGINS":  func(c *Config) { c.CORSAllowedOrigins = nil },
	}

	for name, clear := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := validConfig()
			clear(&cfg)
			err := cfg.Validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), name)
		})
	}
}

func TestValidate_ReportsEveryProblemAtOnce(t *testing.T) {
	err := (&Config{SMTPTimeout: 10, CORSAllowedOrigins: []string{"x"}}).Validate()
	require.Error(t, err)
	for _, name := range []string{"DATABASE_URL", "MINIO_ENDPOINT", "SMTP_HOST", "CRON_CAPSULE_DISPATCH"} {
		require.Contains(t, err.Error(), name)
	}
}

// A binary that starts with a baked-in password is worse than one that refuses.
func TestValidate_NoDefaultCredentials(t *testing.T) {
	var cfg Config
	require.Error(t, cfg.Validate())
	require.Empty(t, cfg.MinioSecretKey)
	require.Empty(t, cfg.SMTPPassword)
}
