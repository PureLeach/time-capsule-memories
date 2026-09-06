// Package config loads and validates the application's environment settings.
package config

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Secrets deliberately carry no defaults: a missing credential must fail
// startup rather than fall back to a value baked into the binary.
type Config struct {
	PostgresURL string `env:"DATABASE_URL"`

	MinioAccessKey  string `env:"MINIO_ROOT_USER"`
	MinioSecretKey  string `env:"MINIO_ROOT_PASSWORD"`
	MinioEndpoint   string `env:"MINIO_ENDPOINT"`
	MinioUseSSL     bool   `env:"MINIO_USE_SSL" env-default:"false"`
	MinioBucketName string `env:"MINIO_BUCKET_NAME" env-default:"time-capsule"`

	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     string `env:"SMTP_PORT"`
	SMTPFrom     string `env:"SMTP_FROM"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPTimeout  int    `env:"SMTP_TIMEOUT" env-default:"10"`

	CronCapsuleDispatch string `env:"CRON_CAPSULE_DISPATCH"`

	LogLevel string `env:"LOG_LEVEL" env-default:"info"`

	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" env-separator:"," env-default:"http://frontend.localhost,http://localhost:8001"`

	// Exposes POST /send-test-email, which is unauthenticated and can send mail
	// to any address through the configured relay.
	EnableTestEmailEndpoint bool `env:"ENABLE_TEST_EMAIL_ENDPOINT" env-default:"false"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		slog.Debug(".env file not loaded; falling back to environment", "error", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate reports every problem at once, so a misconfigured deployment is not
// fixed one restart at a time.
func (c *Config) Validate() error {
	required := []struct {
		name  string
		value string
	}{
		{"DATABASE_URL", c.PostgresURL},
		{"MINIO_ENDPOINT", c.MinioEndpoint},
		{"MINIO_ROOT_USER", c.MinioAccessKey},
		{"MINIO_ROOT_PASSWORD", c.MinioSecretKey},
		{"MINIO_BUCKET_NAME", c.MinioBucketName},
		{"SMTP_HOST", c.SMTPHost},
		{"SMTP_PORT", c.SMTPPort},
		{"SMTP_FROM", c.SMTPFrom},
		{"CRON_CAPSULE_DISPATCH", c.CronCapsuleDispatch},
	}

	var errs []error
	for _, field := range required {
		if field.value == "" {
			errs = append(errs, fmt.Errorf("%s is required", field.name))
		}
	}

	if c.SMTPTimeout <= 0 {
		errs = append(errs, errors.New("SMTP_TIMEOUT must be a positive number of seconds"))
	}
	if len(c.CORSAllowedOrigins) == 0 {
		errs = append(errs, errors.New("CORS_ALLOWED_ORIGINS must list at least one origin"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid configuration: %w", errors.Join(errs...))
	}

	return nil
}
