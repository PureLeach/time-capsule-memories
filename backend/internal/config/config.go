package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var instance *Config
var once sync.Once

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
}

func LoadConfig() (*Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки .env файла")
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}

	// config.MinioEndpoint = "minio-api.localhost"
	// config.MinioEndpoint = fmt.Sprintf("%s:%s", config.MinioHost, config.MinioPort)

	return &config, nil
}

// Функция для получения конфигурации (Singleton)
func GetConfig() *Config {
	// Используем sync.Once для гарантии, что конфигурация будет загружена только один раз
	once.Do(func() {
		var err error
		instance, err = LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("Не удалось загрузить конфигурацию: %v", err))
		}
	})

	return instance
}
