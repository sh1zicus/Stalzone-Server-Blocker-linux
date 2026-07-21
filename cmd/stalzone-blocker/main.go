package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sh1zicus/stalzone-server-blocker/internal/app"
	"github.com/sh1zicus/stalzone-server-blocker/internal/daemon"
	"github.com/sh1zicus/stalzone-server-blocker/internal/logger"
	"github.com/sh1zicus/stalzone-server-blocker/internal/nft"
	"github.com/sh1zicus/stalzone-server-blocker/internal/util"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "--daemon" {
		runDaemon()
		return
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func runDaemon() {

	if err := util.EnsureDirs(); err != nil {
		log.Fatal(err)
	}

	slog, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := nft.CheckPermissions(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}

	d := daemon.New(slog)

	if err := d.Reload(); err != nil {
		slog.Error("initial load failed", "err", err)
	}

	if err := d.Run(); err != nil {
		log.Fatal(err)
	}
}
