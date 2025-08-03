// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	_ "embed"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/piweb/internal/audio"
	"syscall/js"
)

var (
	//go:embed "gameloop.js"
	gameLoopJS []byte

	//go:embed "canvas.js"
	canvasJS []byte
)

var window = js.Global()

var keyboard = StartKeyboard()

func Run() {
	piaudio.Backend = &audio.Backend{}

	window.Set("api", api)
	snapshotPi()

	window.Call("eval", string(gameLoopJS))
	window.Call("eval", string(canvasJS))

	window.Call("prepareCanvas")
}
