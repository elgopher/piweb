// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi/pipad"
	"syscall/js"
)

const maxGamepads = 4

var navigator = window.Get("navigator")

func StartGamepad() *Gamepad {
	k := &Gamepad{}
	for i := 0; i < maxGamepads; i++ {
		k.state[i] = newGamepadState()
	}
	k.reusedState = newGamepadState()
	k.Start()
	return k
}

type Gamepad struct {
	connectionEvents []pipad.EventConnection
	state            [maxGamepads]gamepadState
	reusedState      gamepadState // to avoid allocation
}

func (k *Gamepad) Start() {
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

	gamepadSnapshot := navigator.Call("getGamepads") // This cannot be avoided. Allocates a lot!

	for i := 0; i < maxGamepads; i++ {
		k.getGamepadState(gamepadSnapshot.Index(i), k.reusedState)
		k.state[i].PublishEvents(k.reusedState, i)
	}
}

func (k *Gamepad) getGamepadState(gamepad js.Value, out gamepadState) {
	if gamepad.Truthy() {
		buttons := gamepad.Get("buttons")

		for webBtn, pipadBtn := range gamepadMapping {
			out[pipadBtn] = buttons.Index(webBtn).Get("pressed").Bool()
		}

		axes := gamepad.Get("axes")

		verticalAxis := axes.Index(1).Float()
		if verticalAxis > 0.5 {
			out[pipad.Bottom] = true
		}
		if verticalAxis < -0.5 {
			out[pipad.Top] = true
		}

		horizontalAxis := axes.Index(0).Float()
		if horizontalAxis > 0.5 {
			out[pipad.Right] = true
		}
		if horizontalAxis < -0.5 {
			out[pipad.Left] = true
		}
	} else {
		for _, button := range gamepadMapping {
			out[button] = false
		}
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

func newGamepadState() gamepadState {
	out := gamepadState{}
	for _, pipadBtn := range gamepadMapping {
		out[pipadBtn] = false
	}
	return out
}

type gamepadState map[pipad.Button]bool

func (s gamepadState) PublishEvents(newState gamepadState, player int) {
	for button, oldBtnState := range s {
		newBtnState := newState[button]
		if newBtnState != oldBtnState {
			eventType := pipad.EventDown
			if !newBtnState {
				eventType = pipad.EventUp
			}

			event := pipad.EventButton{
				Type:   eventType,
				Button: button,
				Player: player,
			}
			pipad.ButtonTarget().Publish(event)

			s[button] = newBtnState
		}
	}
}
