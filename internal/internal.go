// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	_ "embed"
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pievent"
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

var (
	keyboard = StartKeyboard()
	gamepad  = StartGamepad()
	mouse    Mouse
)

var paused bool

func Run() {
	piaudio.Backend = &audio.Backend{}

	window.Set("api", api)
	snapshotPi()

	window.Call("eval", string(gameLoopJS))
	window.Call("eval", string(canvasJS))
	window.Call("prepareCanvas")

	mouse.Start(window.Get("canvas"))

	pidebug.Target().SubscribeAll(func(event pidebug.Event, handler pievent.Handler) {
		switch event {
		case pidebug.EventPause:
			paused = true
		case pidebug.EventResume:
			paused = false
		}
	})
}
