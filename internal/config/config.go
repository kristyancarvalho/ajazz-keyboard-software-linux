package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/kristyancarvalho/ak820pro/internal/keyboard"
)

type Config struct {
	LastMode      keyboard.ModeConfig   `json:"last_mode"`
	Themes        map[string][][3]uint8 `json:"themes"`
	LastThemeName string                `json:"last_theme_name"`
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Themes: make(map[string][][3]uint8)}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return &Config{Themes: make(map[string][][3]uint8)}, nil
	}
	cfg := &Config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Themes == nil {
		cfg.Themes = make(map[string][][3]uint8)
	}
	return cfg, nil
}

func (c *Config) Save() error {
	if c == nil {
		return errors.New("config is nil")
	}
	path, err := configPath()
	if err != nil {
		return err
	}
	if c.Themes == nil {
		c.Themes = make(map[string][][3]uint8)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func configPath() (string, error) {
	if dir := os.Getenv("AK820PRO_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "config.json"), nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "ak820pro", "config.json"), nil
}
