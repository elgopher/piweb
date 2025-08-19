// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build windows || linux || darwin || unix

package piweb

import (
	"context"
	"os"
	"os/signal"

	"github.com/elgopher/piweb/internal/server"
)

func run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	server.GoBuild = &GoBuild
	server.ReleaseGoBuild = &ReleaseGoBuild

	server.Start(ctx, 8080)
}
