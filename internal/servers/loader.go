package servers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

func Load(path string) ([]model.Pool, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var file model.ServersFile

	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse servers.json: %w", err)
	}

	return file.Pools, nil
}

func FindPath() (string, error) {

	candidates := []string{
		"./data/Servers.json",
		"./Servers.json",
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			abs, _ := filepath.Abs(p)
			return abs, nil
		}
	}

	return "", os.ErrNotExist
}
