// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/piweb/internal/window"
	"syscall/js"
	"time"
)

// api provides functions available in JavaScript code
var api = window.NewObject()

func init() {
	api.Set("init", js.FuncOf(piInit))
	api.Set("tick", js.FuncOf(tick))
	snapshotPi()
}

// snapshotPi makes a snapshot of Pi values to avoid memory allocation
// on each call from JS to GO.
func snapshotPi() {
	api.Set("tps", js.ValueOf(pi.TPS()))
}

func piInit(this js.Value, args []js.Value) any {
	if pi.Init != nil {
		pi.Init()
	}
	piloop.Target().Publish(piloop.EventInit)

	return nil
}

// When true, indicates that the last update+draw cycle
// exceeded the tick duration of 1/TPS (e.g., 33 ms for TPS=30).
var skipNextDraw bool

func tick(this js.Value, args []js.Value) any {
	started := time.Now()

	snapshotPi()

	piloop.Target().Publish(piloop.EventFrameStart)

	pi.Update()
	piloop.Target().Publish(piloop.EventUpdate)

	if !skipNextDraw {
		pi.Draw()
		piloop.Target().Publish(piloop.EventDraw)
	} else {
		skipNextDraw = false
	}

	tickDuration := 1.0 / float64(pi.TPS())
	elapsed := time.Since(started).Seconds()
	if elapsed > tickDuration {
		skipNextDraw = true // game is too slow. Try to keep up by discarding next pi.Draw()
	}

	pi.Time += tickDuration
	pi.Frame++

	return nil
}
