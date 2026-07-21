package tui

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/viewport"

	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

type Model struct {
	cfg *model.Config
	log *slog.Logger

	state *model.State

	viewport viewport.Model

	searchMode bool
	search     string

	status string

	stalzoneRunning bool

	width  int
	height int
}
