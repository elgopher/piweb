// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/piweb/internal/window"
	"syscall/js"
)

// api provides functions available in JavaScript code
var api = window.NewObject()

func init() {
	api.Set("tps", js.ValueOf(pi.TPS()))
	api.Set("init", js.FuncOf(piInit))
	api.Set("tick", js.FuncOf(tick))
}

func piInit(this js.Value, args []js.Value) any {
	if pi.Init != nil {
		pi.Init()
	}
	piloop.Target().Publish(piloop.EventInit)

	return nil
}

func tick(this js.Value, args []js.Value) any {
	// make a snapshot of pi values to avoid memory allocation
	// on each call from JS to GO.
	api.Set("tps", js.ValueOf(pi.TPS()))

	piloop.Target().Publish(piloop.EventFrameStart)

	pi.Update()
	piloop.Target().Publish(piloop.EventUpdate)
	pi.Draw()
	piloop.Target().Publish(piloop.EventDraw)

	return nil
}
