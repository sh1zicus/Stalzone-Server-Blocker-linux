package config

import (
	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

func ApplySelection(cfg *model.Config, pools []model.Pool) {

	selected := make(map[string]struct{}, len(cfg.Selected))

	for _, name := range cfg.Selected {
		selected[name] = struct{}{}
	}

	for pi := range pools {
		for ti := range pools[pi].Tunnels {

			_, ok := selected[pools[pi].Tunnels[ti].Name]

			pools[pi].Tunnels[ti].Selected = ok
		}
	}
}

func UpdateSelection(cfg *model.Config, pools []model.Pool) {

	cfg.Selected = cfg.Selected[:0]

	for _, pool := range pools {
		for _, t := range pool.Tunnels {

			if t.Selected {
				cfg.Selected = append(cfg.Selected, t.Name)
			}
		}
	}
}
