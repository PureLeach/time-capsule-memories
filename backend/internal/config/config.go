package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Config содержит настройки подключения к базе данных
type Config struct {
	DBName   string `env:"POSTGRES_DB" env-default:"time_capsule"`
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"password"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
}

// LoadConfig загружает значения конфигурации из переменных окружения
func LoadConfig() (Config, error) {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
