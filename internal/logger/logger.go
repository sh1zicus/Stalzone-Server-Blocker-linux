package logger

import (
	"log/slog"
	"os"

	"github.com/sh1zicus/stalzone-server-blocker/internal/util"
)

func New() (*slog.Logger, error) {

	file, err := os.OpenFile(
		util.LogFile(),
				 os.O_CREATE|os.O_APPEND|os.O_WRONLY,
				 0644,
	)

	if err != nil {
		return nil, err
	}

	return slog.New(
		slog.NewTextHandler(file, nil),
	), nil
}
