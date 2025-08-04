// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pimouse"
	"syscall/js"
)

type Mouse struct {
	prev   pi.Position
	events []event
}

type event struct {
	isMove bool
	pimouse.EventMove
	pimouse.EventButton
}

func (m *Mouse) Start(canvas js.Value) {
	canvas.Call("addEventListener", "mousedown", js.FuncOf(m.mouseDown))
	canvas.Call("addEventListener", "mouseup", js.FuncOf(m.mouseUp))
	window.Call("addMouseMoveListener", canvas, js.FuncOf(m.mouseMove))
	window.Call("addEventListener", "blur", js.FuncOf(m.blur))
}

var mouseMap = map[int]pimouse.Button{
	0: pimouse.Left,
	2: pimouse.Right,
}

func (m *Mouse) mouseDown(this js.Value, args []js.Value) any {
	btnWeb := args[0].Get("button").Int()
	btn, found := mouseMap[btnWeb]
	if !found {
		return nil
	}

	m.events = append(m.events, event{
		isMove: false,
		EventButton: pimouse.EventButton{
			Type:   pimouse.EventButtonDown,
			Button: btn,
		},
	})

	return nil
}

func (m *Mouse) mouseUp(this js.Value, args []js.Value) any {
	btnWeb := args[0].Get("button").Int()
	btn, found := mouseMap[btnWeb]
	if !found {
		return nil
	}

	m.events = append(m.events, event{
		isMove: false,
		EventButton: pimouse.EventButton{
			Type:   pimouse.EventButtonUp,
			Button: btn,
		},
	})

	return nil
}

func (m *Mouse) mouseMove(this js.Value, args []js.Value) any {
	prev := m.prev
	m.prev.X = args[0].Int()
	m.prev.Y = args[1].Int()

	if prev == m.prev {
		return nil
	}

	m.events = append(m.events, event{
		isMove: true,
		EventMove: pimouse.EventMove{
			Position: m.prev,
			Previous: prev,
		},
	})

	return nil
}

func (m *Mouse) blur(this js.Value, args []js.Value) any {
	return nil
}

func (m *Mouse) Update() {
	for _, e := range m.events {
		if e.isMove {
			if !paused {
				pimouse.MoveTarget().Publish(e.EventMove)
			}
			pimouse.MoveDebugTarget().Publish(e.EventMove)
		} else {
			if !paused {
				pimouse.ButtonTarget().Publish(e.EventButton)
			}
			pimouse.ButtonDebugTarget().Publish(e.EventButton)
		}
	}
	m.events = m.events[:0]
}
