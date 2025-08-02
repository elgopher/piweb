// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piloop"
	"syscall/js"
	"time"
)

// api provides functions available in JavaScript code
var api = window.Get("Object").New()

func init() {
	api.Set("init", js.FuncOf(piInit))
	api.Set("tick", js.FuncOf(tick))
}

// snapshotPi makes a snapshot of Pi values to avoid memory allocation
// on each call from JS to GO.
func snapshotPi() {
	api.Set("tps", js.ValueOf(pi.TPS()))

	screen := pi.Screen()
	api.Set("screenWidth", screen.W())
	api.Set("screenHeight", screen.H())
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

	ticks := args[0].Int()

	for i := 0; i < ticks; i++ {
		piloop.Target().Publish(piloop.EventFrameStart)

		pi.Update()
		piloop.Target().Publish(piloop.EventUpdate)

		if i == ticks-1 { // draw only on the last tick
			if !skipNextDraw {
				pi.Draw()
				piloop.Target().Publish(piloop.EventDraw)

				data := window.Get("imageData").Get("data")
				CopyCanvasToUint8ClampedArray(data, pi.Screen())
				window.Call("updateCanvas")
			} else {
				skipNextDraw = false
			}
		}

		pi.Time += 1.0 / float64(pi.TPS())
		pi.Frame++
	}

	elapsed := time.Since(started).Seconds()
	if elapsed > 1.0/float64(pi.TPS()) {
		skipNextDraw = true // game is too slow. Try to keep up by discarding next pi.Draw()
	}

	snapshotPi()

	return nil
}
