package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		BotConfig `yaml:"bot"`
		Postgres  `yaml:"postgres"`
	}

	BotConfig struct {
		BotID       string `yaml:"bot_id" env:"BOT_ID"`
		BotToken    string `yaml:"bot_token" env:"BOT_TOKEN"`
		BotUsername string `yaml:"bot_username" env:"BOT_USERNAME"`
	}

	Postgres struct {
		Host     string `yaml:"host" env:"PSQL_HOST"`
		Port     string `yaml:"port" env:"PSQL_PORT"`
		Username string `yaml:"username" env:"PSQL_USERNAME"`
		DBName   string `yaml:"dbname" env:"PSQL_NAME"`
		SSLMode  string `yaml:"sslmode" env:"PSQL_SSLMODE"`
		Password string `yaml:"password" env:"PSQL_PASSWORD"`
	}
)

func NewConfig(filePath string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		if envErr := cleanenv.ReadEnv(&cfg); envErr != nil {
			return nil, fmt.Errorf("config: %v", err)
		}
	}
	return &cfg, nil
}
