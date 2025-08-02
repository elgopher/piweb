// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build js && wasm

package piweb

import (
	"github.com/elgopher/piweb/internal"
)

func Run() {
	internal.Run()
	// do not close the wasm module (for now):
	select {}
}
