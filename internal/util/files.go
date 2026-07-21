package util

import (
	"os"
	"path/filepath"
)

func FindServersJSON() (string, error) {

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
