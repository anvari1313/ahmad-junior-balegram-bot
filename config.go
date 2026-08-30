package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds all application settings loaded from the environment.
type Config struct {
	Telegram struct {
		Token  string `env:"TELEGRAM_BOT_TOKEN"`
		APIURL string `env:"TELEGRAM_API_URL" envDefault:"https://tapi.bale.ai"`
	}
	ZAI struct {
		APIKey  string `env:"ZAI_API_KEY"`
		BaseURL string `env:"ZAI_BASE_URL" envDefault:"https://api.z.ai/api/paas/v4/"`
		Model   string `env:"ZAI_MODEL"   envDefault:"glm-5"`
	}
}

// loadConfig reads configuration from environment variables.
func loadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// validate checks that required fields are present, returning a clear error
// naming the missing variable.
func (c *Config) validate() error {
	if c.Telegram.Token == "" {
		return fmt.Errorf("config: TELEGRAM_BOT_TOKEN is not set")
	}
	if c.ZAI.APIKey == "" {
		return fmt.Errorf("config: ZAI_API_KEY is not set")
	}
	return nil
}
