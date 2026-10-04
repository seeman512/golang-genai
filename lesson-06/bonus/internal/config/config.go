package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Config describes the server configuration stored in JSON.
type Config struct {
	ServerPort    int    `json:"server_port"`
	Environment   string `json:"environment"`
	DatabaseURL   string `json:"database_url,omitempty"`
	DebugMode     bool   `json:"debug_mode"`
	AdminPassword string `json:"-"`
}

// FileStorage abstracts the file operations needed by the configuration
// manager. Its use allows callers and tests to provide storage failures
// without requiring real filesystem failures.
type FileStorage interface {
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
}

type osFileStorage struct{}

func (osFileStorage) Create(name string) (io.WriteCloser, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func (osFileStorage) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// SaveConfig writes cfg to path as indented JSON.
func SaveConfig(path string, cfg Config) error {
	return SaveConfigWithStorage(osFileStorage{}, path, cfg)
}

// SaveConfigWithStorage writes cfg using the supplied storage implementation.
func SaveConfigWithStorage(storage FileStorage, path string, cfg Config) (err error) {
	if storage == nil {
		return fmt.Errorf("save config %s: file storage is nil", path)
	}

	file, err := storage.Create(path)
	if err != nil {
		return fmt.Errorf("create config file %s: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close config file %s: %w", path, closeErr))
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("write config %s: %w", path, err)
	}

	return nil
}

// LoadConfig reads a JSON configuration from path.
func LoadConfig(path string) (Config, error) {
	return LoadConfigWithStorage(osFileStorage{}, path)
}

// LoadConfigWithStorage reads a JSON configuration using the supplied storage
// implementation.
func LoadConfigWithStorage(storage FileStorage, path string) (cfg Config, err error) {
	if storage == nil {
		return cfg, fmt.Errorf("load config %s: file storage is nil", path)
	}

	file, err := storage.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("open file %s error: %w", path, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close config file %s: %w", path, closeErr))
		}
	}()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("parse json %s error: %w", path, err)
	}

	// A second decode ensures that trailing JSON values are not silently
	// accepted as part of an otherwise valid configuration.
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return cfg, fmt.Errorf("parse json %s error: multiple JSON values", path)
		}
		return cfg, fmt.Errorf("parse json %s error: trailing data: %w", path, err)
	}

	return cfg, nil
}

// ValidateConfig checks the port range and deployment environment.
func ValidateConfig(cfg Config) error {
	if cfg.ServerPort < 1024 || cfg.ServerPort > 65535 {
		return fmt.Errorf("server port %d must be between 1024 and 65535", cfg.ServerPort)
	}

	switch cfg.Environment {
	case "development", "staging", "production":
		return nil
	default:
		return fmt.Errorf("environment %q must be development, staging, or production", cfg.Environment)
	}
}
