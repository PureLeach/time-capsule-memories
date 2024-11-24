package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

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
	MinioHost       string `env:"MINIO_HOST" env-default:"localhost"`
	MinioPort       string `env:"MINIO_ENDPOINT" env-default:"9000"`
	MinioUseSSL     bool   `env:"MINIO_USE_SSL" env-default:"false"`
	MinioBucketName string `env:"MINIO_BUCKET_NAME" env-default:"time-capsule"`
	MinioEndpoint   string
}

func LoadConfig() (Config, error) {
	var config Config

	// Загружаем переменные из .env файла
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file")
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return Config{}, err
	}

	// Формирование MinioEndpoint
	config.MinioEndpoint = fmt.Sprintf("%s:%s", config.MinioHost, config.MinioPort)

	return config, nil
}
