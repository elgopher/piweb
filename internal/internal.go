// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	_ "embed"
	"syscall/js"

	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/piweb/internal/audio"
)

var window = js.Global()

var (
	keyboard = StartKeyboard()
	gamepad  *Gamepad
	mouse    Mouse
)

var paused bool

func Run() {
	piaudio.Backend = &audio.Backend{}

	window.Set("api", api)
	snapshotPi()

	eventsByteBuffer := NewByteBuffer(window.Get("gamepad").Get("events"))
	gamepad = StartGamepad(eventsByteBuffer)

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
