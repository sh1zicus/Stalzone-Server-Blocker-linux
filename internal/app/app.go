package app

import (
	"fmt"
	"os"

	"github.com/sh1zicus/stalzone-server-blocker/internal/config"
	"github.com/sh1zicus/stalzone-server-blocker/internal/logger"
	"github.com/sh1zicus/stalzone-server-blocker/internal/nft"
	"github.com/sh1zicus/stalzone-server-blocker/internal/servers"
	"github.com/sh1zicus/stalzone-server-blocker/internal/tui"
	"github.com/sh1zicus/stalzone-server-blocker/internal/util"
)

func Run() error {

	if err := util.EnsureDirs(); err != nil {
		return err
	}

	log, err := logger.New()
	if err != nil {
		return err
	}

	// Проверяем права на nftables до запуска TUI.
	if err := nft.CheckPermissions(); err != nil {
		exe, _ := os.Executable()
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		fmt.Fprintf(os.Stderr, "  sudo setcap cap_net_admin+ep %s\n", exe)
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	path, err := util.FindServersJSON()
	if err != nil {
		return err
	}

	pools, err := servers.Load(path)
	if err != nil {
		return err
	}

	config.ApplySelection(cfg, pools)

	log.Info("pools loaded", "count", len(pools))

	return tui.Run(cfg, log, pools)
}
