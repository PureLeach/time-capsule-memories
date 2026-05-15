package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
)

// Config holds all environment configuration variables for the application
type Config struct {
	// PostgreSQL
	DBUser      string `env:"POSTGRES_USER" env-default:"user"`
	DBPassword  string `env:"POSTGRES_PASSWORD" env-default:"1234"`
	DBHost      string `env:"POSTGRES_HOST" env-default:"localhost"`
	DBPort      string `env:"POSTGRES_PORT" env-default:"5432"`
	DBName      string `env:"POSTGRES_DB_NAME" env-default:"time_capsule"`
	PostgresURL string `env:"DATABASE_URL"`

	// MinIO
	MinioAccessKey  string `env:"MINIO_ROOT_USER" env-default:"minioaccesskey"`
	MinioSecretKey  string `env:"MINIO_ROOT_PASSWORD" env-default:"miniosecretkey"`
	MinioEndpoint   string `env:"MINIO_ENDPOINT" env-default:"minio-api.localhost"`
	MinioUseSSL     bool   `env:"MINIO_USE_SSL" env-default:"false"`
	MinioBucketName string `env:"MINIO_BUCKET_NAME" env-default:"time-capsule"`

	// SMTP
	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     string `env:"SMTP_PORT"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPFrom     string `env:"SMTP_FROM"`
	SMTPTimeout  int    `env:"SMTP_TIMEOUT" env-default:"10"`

	// Scheduler
	CronCapsuleDispatch string `env:"CRON_CAPSULE_DISPATCH"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

func LoadConfig() (*Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		// Logging may not be initialised yet; this falls through the default
		// slog text handler and is normally filtered at info level.
		slog.Debug(".env file not loaded; falling back to environment", "error", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to read environment variables: %w", err)
	}

	return &config, nil
}

// GetConfig returns a singleton instance of the Config.
// It initializes the config once using sync.Once.
func GetConfig() *Config {
	once.Do(func() {
		var err error
		instance, err = LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load configuration: %v", err))
		}
	})
	return instance
}
