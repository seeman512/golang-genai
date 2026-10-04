package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerPort    int    `json:"server_port"`
	Environment   string `json:"environment"`
	DatabaseURL   string `json:"database_url,omitempty"`
	DebugMode     bool   `json:"debug_mode"`
	AdminPassword string `json:"-"`
}

func SaveConfig(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")

	if err != nil {
		return fmt.Errorf("parse json error: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("save file %s error: %w", path, err)
	}

	return nil
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("open file %s error: %w", path, err)
	}

	err = json.Unmarshal(data, &cfg)

	if err != nil {
		return cfg, fmt.Errorf("parse json %s error: %w", path, err)
	}

	return cfg, nil
}
