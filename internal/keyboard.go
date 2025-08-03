// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi/pikey"
	"syscall/js"
)

func StartKeyboard() *Keyboard {
	k := &Keyboard{
		pressedKeys: make(map[pikey.Key]struct{}),
	}
	k.Start()
	return k
}

type Keyboard struct {
	events      []pikey.Event
	pressedKeys map[pikey.Key]struct{}
}

func (k *Keyboard) Start() {
	window.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		code := args[0].Get("code").String()
		key, found := keymap[code]
		if !found {
			return nil
		}

		repeat := args[0].Get("repeat").Bool()
		if repeat == true {
			return nil
		}

		event := pikey.Event{Type: pikey.EventDown, Key: key}
		k.events = append(k.events, event)
		k.pressedKeys[key] = struct{}{}

		if virtualKey, virtualKeyFound := virtualKeys[key]; virtualKeyFound {
			additionalEvent := pikey.Event{Type: pikey.EventDown, Key: virtualKey}
			k.events = append(k.events, additionalEvent)
			k.pressedKeys[virtualKey] = struct{}{}
		}

		return nil
	}))

	window.Call("addEventListener", "keyup", js.FuncOf(func(this js.Value, args []js.Value) any {
		code := args[0].Get("code").String()
		key, found := keymap[code]
		if !found {
			return nil
		}

		event := pikey.Event{Type: pikey.EventUp, Key: key}
		k.events = append(k.events, event)
		delete(k.pressedKeys, key)

		if virtualKey, virtualKeyFound := virtualKeys[key]; virtualKeyFound {
			additionalEvent := pikey.Event{Type: pikey.EventUp, Key: virtualKey}
			k.events = append(k.events, additionalEvent)
			delete(k.pressedKeys, virtualKey)
		}

		return nil
	}))

	window.Call("addEventListener", "blur", js.FuncOf(func(this js.Value, args []js.Value) any {
		for key := range k.pressedKeys {
			event := pikey.Event{Type: pikey.EventUp, Key: key}
			k.events = append(k.events, event)
		}
		clear(k.pressedKeys)

		return nil
	}))
}

var keymap = map[string]pikey.Key{
	"KeyA":         pikey.A,
	"KeyB":         pikey.B,
	"KeyC":         pikey.C,
	"KeyD":         pikey.D,
	"KeyE":         pikey.E,
	"KeyF":         pikey.F,
	"KeyG":         pikey.G,
	"KeyH":         pikey.H,
	"KeyI":         pikey.I,
	"KeyJ":         pikey.J,
	"KeyK":         pikey.K,
	"KeyL":         pikey.L,
	"KeyM":         pikey.M,
	"KeyN":         pikey.N,
	"KeyO":         pikey.O,
	"KeyP":         pikey.P,
	"KeyQ":         pikey.Q,
	"KeyR":         pikey.R,
	"KeyS":         pikey.S,
	"KeyT":         pikey.T,
	"KeyU":         pikey.U,
	"KeyV":         pikey.V,
	"KeyW":         pikey.W,
	"KeyX":         pikey.X,
	"KeyY":         pikey.Y,
	"KeyZ":         pikey.Z,
	"AltLeft":      pikey.AltLeft,
	"AltRight":     pikey.AltRight,
	"ArrowDown":    pikey.Down,
	"ArrowLeft":    pikey.Left,
	"ArrowRight":   pikey.Right,
	"ArrowUp":      pikey.Up,
	"Backquote":    pikey.Backquote,
	"Backslash":    pikey.Backslash,
	"Backspace":    pikey.Backspace,
	"BracketLeft":  pikey.BracketLeft,
	"BracketRight": pikey.BracketRight,
	"CapsLock":     pikey.CapsLock,
	"Comma":        pikey.Comma,
	"ControlLeft":  pikey.CtrlLeft,
	"ControlRight": pikey.CtrlRight,
	"Digit0":       pikey.Digit0,
	"Digit1":       pikey.Digit1,
	"Digit2":       pikey.Digit2,
	"Digit3":       pikey.Digit3,
	"Digit4":       pikey.Digit4,
	"Digit5":       pikey.Digit5,
	"Digit6":       pikey.Digit6,
	"Digit7":       pikey.Digit7,
	"Digit8":       pikey.Digit8,
	"Digit9":       pikey.Digit9,
	"Enter":        pikey.Enter,
	"Equal":        pikey.Equal,
	"Escape":       pikey.Esc,
	"F1":           pikey.F1,
	"F2":           pikey.F2,
	"F3":           pikey.F3,
	"F4":           pikey.F4,
	"F5":           pikey.F5,
	"F6":           pikey.F6,
	"F7":           pikey.F7,
	"F8":           pikey.F8,
	"F9":           pikey.F9,
	"F10":          pikey.F10,
	"F11":          pikey.F11,
	"F12":          pikey.F12,
	"Minus":        pikey.Minus,
	"Period":       pikey.Period,
	"Quote":        pikey.Quote,
	"Semicolon":    pikey.Semicolon,
	"ShiftLeft":    pikey.ShiftLeft,
	"ShiftRight":   pikey.ShiftRight,
	"Slash":        pikey.Slash,
	"Space":        pikey.Space,
	"Tab":          pikey.Tab,
	"Alt":          pikey.Alt,     // virtual key
	"Control":      pikey.Control, // virtual key
	"Shift":        pikey.Shift,   // virtual key
}

var virtualKeys = map[pikey.Key]pikey.Key{
	pikey.AltLeft:    pikey.Alt,
	pikey.AltRight:   pikey.Alt,
	pikey.CtrlLeft:   pikey.Control,
	pikey.CtrlRight:  pikey.Control,
	pikey.ShiftLeft:  pikey.Shift,
	pikey.ShiftRight: pikey.Shift,
}

func (k *Keyboard) Update() {
	for _, event := range k.events {
		pikey.Target().Publish(event)
	}

	k.events = k.events[:0]
}
