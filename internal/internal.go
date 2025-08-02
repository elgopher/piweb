// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	_ "embed"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/piweb/internal/audio"
	"github.com/elgopher/piweb/internal/window"
)

//go:embed "backend.js"
var backendJS []byte

func Run() {
	piaudio.Backend = &audio.Backend{}

	window.Set("api", api)
	window.Eval(string(backendJS))
}
