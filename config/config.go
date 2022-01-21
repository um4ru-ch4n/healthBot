package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		BotConfig `yaml:"bot"`
	}

	BotConfig struct {
		BotID       string `yaml:"bot_id" env:"bot_id"`
		BotToken    string `yaml:"bot_token" env:"bot_token"`
		BotUsername string `yaml:"bot_username" env:"bot_username"`
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
