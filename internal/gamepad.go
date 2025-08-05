// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"syscall/js"

	"github.com/elgopher/pi/pipad"
)

func StartGamepad(events *ByteBuffer) *Gamepad {
	g := &Gamepad{}
	g.Start(events)
	return g
}

type Gamepad struct {
	connectionEvents []pipad.EventConnection

	events *ByteBuffer
}

func (k *Gamepad) Start(events *ByteBuffer) {
	k.events = events

	window.Call("addEventListener", "gamepadconnected", js.FuncOf(k.gamepadconnected))
	window.Call("addEventListener", "gamepaddisconnected ", js.FuncOf(k.gamepaddisconnected))
}

func (k *Gamepad) gamepadconnected(this js.Value, args []js.Value) any {
	gamepad := args[0].Get("gamepad")

	index := gamepad.Get("index").Int()
	event := pipad.EventConnection{
		Type:   pipad.EventConnect,
		Player: index,
	}
	k.connectionEvents = append(k.connectionEvents, event)

	return nil
}

func (k *Gamepad) gamepaddisconnected(this js.Value, args []js.Value) any {
	gamepad := args[0].Get("gamepad")

	index := gamepad.Get("index").Int()
	event := pipad.EventConnection{
		Type:   pipad.EventDisconnect,
		Player: index,
	}
	k.connectionEvents = append(k.connectionEvents, event)

	return nil
}

func (k *Gamepad) Update() {
	for _, event := range k.connectionEvents {
		pipad.ConnectionTarget().Publish(event)
	}
	k.connectionEvents = k.connectionEvents[:0]

	k.events.Read()
	data := k.events.Data()

	for i := 0; i < len(data)-2; i += 3 {
		event := data[i : i+3]
		var eventType = pipad.EventUp
		if event[0] == 2 {
			eventType = pipad.EventDown
		}
		padIndex := event[1]
		button := gamepadMapping[int(event[2])]

		newEvent := pipad.EventButton{
			Type:   eventType,
			Player: int(padIndex),
			Button: button,
		}

		pipad.ButtonTarget().Publish(newEvent)
	}
}

var gamepadMapping = map[int]pipad.Button{
	0:  pipad.A,
	1:  pipad.B,
	2:  pipad.X,
	3:  pipad.Y,
	12: pipad.Top,
	13: pipad.Bottom,
	14: pipad.Left,
	15: pipad.Right,
}
