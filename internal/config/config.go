package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

func Path() string {

	home, _ := os.UserHomeDir()

	return filepath.Join(
		home,
		".config",
		"stalzone-blocker",
		"config.json",
	)
}

func Load() (*model.Config, error) {

	cfg := model.DefaultConfig()

	data, err := os.ReadFile(Path())
	if err != nil {
		return &cfg, nil
	}

	err = json.Unmarshal(data, &cfg)

	return &cfg, err
}

func Save(cfg *model.Config) error {

	path := Path()

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
