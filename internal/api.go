// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/piweb/internal/window"
	"syscall/js"
)

// api provides functions available in backend.js
var api = map[string]any{
	"screenSize": js.FuncOf(screenSize),
	"tps":        js.FuncOf(tps),
	"init":       js.FuncOf(piInit),
	"update":     js.FuncOf(update),
	"draw":       js.FuncOf(draw),
}

func screenSize(this js.Value, args []js.Value) interface{} {
	screen := pi.Screen()

	o := window.NewObject()
	o.Set("w", screen.W())
	o.Set("h", screen.H())

	return o
}

func tps(this js.Value, args []js.Value) interface{} {
	return pi.TPS()
}

func piInit(this js.Value, args []js.Value) interface{} {
	if pi.Init != nil {
		pi.Init()
	}
	piloop.Target().Publish(piloop.EventInit)

	return nil
}

func update(this js.Value, args []js.Value) interface{} {
	pi.Update()
	piloop.Target().Publish(piloop.EventUpdate)
	return nil
}

func draw(this js.Value, args []js.Value) interface{} {
	pi.Draw()
	piloop.Target().Publish(piloop.EventDraw)
	return nil
}
