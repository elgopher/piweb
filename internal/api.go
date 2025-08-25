// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"syscall/js"
	"time"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/piweb/internal/audio"
)

// api provides functions available in JavaScript code
var api = window.Get("Object").New()

func init() {
	api.Set("tick", js.FuncOf(tick))
}

// snapshotPi makes a snapshot of Pi values to avoid memory allocation
// on each call from JS to GO.
func snapshotPi() {
	api.Set("tps", pi.TPS())

	screen := pi.Screen()
	api.Set("screenWidth", screen.W())
	api.Set("screenHeight", screen.H())
}

// When true, indicates that the last update+draw cycle
// exceeded the tick duration of 1/TPS (e.g., 33 ms for TPS=30).
var skipNextDraw bool

func tick(this js.Value, args []js.Value) any {
	started := time.Now()

	ticks := args[0].Int()

	for i := 0; i < ticks; i++ {
		audio.UpdateTime()

		if !paused {
			piloop.Target().Publish(piloop.EventFrameStart)
		}
		piloop.DebugTarget().Publish(piloop.EventFrameStart)

		// handling input only once
		if i == 0 {
			mouse.Update()
			keyboard.Update()
			gamepad.Update()
		}

		if !paused {
			pi.Update()
			piloop.Target().Publish(piloop.EventUpdate)
		}
		piloop.DebugTarget().Publish(piloop.EventUpdate)

		if !paused {
			piloop.Target().Publish(piloop.EventLateUpdate)
		}
		piloop.DebugTarget().Publish(piloop.EventLateUpdate)

		if i == ticks-1 { // draw only on the last tick
			if !skipNextDraw {
				if !paused {
					pi.Draw()
					piloop.Target().Publish(piloop.EventDraw)
				}
				piloop.DebugTarget().Publish(piloop.EventDraw)

				if !paused {
					piloop.Target().Publish(piloop.EventLateDraw)
				}
				piloop.DebugTarget().Publish(piloop.EventLateDraw)

				updateImageData()
			} else {
				skipNextDraw = false
			}
		}

		pi.Time += 1.0 / float64(pi.TPS())
		pi.Frame++

		audio.SendCommands()
	}

	elapsed := time.Since(started).Seconds()
	if elapsed > 1.0/float64(pi.TPS()) {
		skipNextDraw = true // game is too slow. Try to keep up by discarding next pi.Draw()
	}

	snapshotPi()

	return nil
}

func updateImageData() {
	data := window.Get("imageData").Get("data")
	CopyCanvasToUint8ClampedArray(data, pi.Screen())
	window.Call("updateCanvas")
}
