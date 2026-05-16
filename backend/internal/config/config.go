package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DBUser      string `env:"POSTGRES_USER" env-default:"user"`
	DBPassword  string `env:"POSTGRES_PASSWORD" env-default:"1234"`
	DBHost      string `env:"POSTGRES_HOST" env-default:"localhost"`
	DBPort      string `env:"POSTGRES_PORT" env-default:"5432"`
	DBName      string `env:"POSTGRES_DB_NAME" env-default:"time_capsule"`
	PostgresURL string `env:"DATABASE_URL"`

	MinioAccessKey  string `env:"MINIO_ROOT_USER" env-default:"minioaccesskey"`
	MinioSecretKey  string `env:"MINIO_ROOT_PASSWORD" env-default:"miniosecretkey"`
	MinioEndpoint   string `env:"MINIO_ENDPOINT" env-default:"minio-api.localhost"`
	MinioUseSSL     bool   `env:"MINIO_USE_SSL" env-default:"false"`
	MinioBucketName string `env:"MINIO_BUCKET_NAME" env-default:"time-capsule"`

	SMTPHost     string `env:"SMTP_HOST"`
	SMTPPort     string `env:"SMTP_PORT"`
	SMTPPassword string `env:"SMTP_PASSWORD"`
	SMTPFrom     string `env:"SMTP_FROM"`
	SMTPTimeout  int    `env:"SMTP_TIMEOUT" env-default:"10"`

	CronCapsuleDispatch string `env:"CRON_CAPSULE_DISPATCH"`

	LogLevel string `env:"LOG_LEVEL" env-default:"info"`

	CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" env-separator:"," env-default:"http://frontend.localhost,http://localhost:8001"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		slog.Debug(".env file not loaded; falling back to environment", "error", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read env: %w", err)
	}

	return &cfg, nil
}
