package util

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "stalzone-blocker")
}

func StateDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "stalzone-blocker")
}

func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func LogFile() string {
	return filepath.Join(StateDir(), "log.txt")
}

func EnsureDirs() error {
	if err := os.MkdirAll(ConfigDir(), 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(StateDir(), 0755); err != nil {
		return err
	}

	return nil
}
