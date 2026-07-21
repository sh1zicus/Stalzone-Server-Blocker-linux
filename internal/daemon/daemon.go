package daemon

import (
	"log/slog"
	"os/exec"
	"time"

	"github.com/sh1zicus/stalzone-server-blocker/internal/config"
	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
	"github.com/sh1zicus/stalzone-server-blocker/internal/nft"
	"github.com/sh1zicus/stalzone-server-blocker/internal/servers"
)

type Daemon struct {
	log    *slog.Logger
	cfg    *model.Config
	pools  []model.Pool
	running bool
}

func New(log *slog.Logger) *Daemon {
	return &Daemon{log: log}
}

func (d *Daemon) Run() error {
	d.log.Info("daemon started")

	for {
		wasRunning := d.running

		d.running = isStalzoneRunning()

		if d.running && !wasRunning {
			d.log.Info("stalzone detected, applying rules in 2s...")
			time.Sleep(2 * time.Second)
			d.applyRules()
		}

		if !d.running && wasRunning {
			d.log.Info("stalzone stopped, resetting rules")
			d.resetRules()
		}

		time.Sleep(1 * time.Second)
	}
}

func (d *Daemon) Reload() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	d.cfg = cfg

	path, err := servers.FindPath()
	if err != nil {
		return err
	}

	pools, err := servers.Load(path)
	if err != nil {
		return err
	}

	config.ApplySelection(d.cfg, pools)
	d.pools = pools

	d.log.Info("config reloaded", "pools", len(pools))

	if d.running {
		d.log.Info("reapplying rules...")
		d.applyRules()
	}

	return nil
}

func (d *Daemon) applyRules() {
	if err := nft.Apply(d.pools); err != nil {
		d.log.Error("apply nft failed", "err", err)
		return
	}
	d.log.Info("rules applied")
}

func (d *Daemon) resetRules() {
	if err := nft.Reset(); err != nil {
		d.log.Error("reset nft failed", "err", err)
		return
	}
	d.log.Info("rules reset")
}

func isStalzoneRunning() bool {
	cmd := exec.Command("pgrep", "-f", "stalzone.exe")
	return cmd.Run() == nil
}
