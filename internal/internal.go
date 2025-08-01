// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import "github.com/elgopher/piweb/internal/window"

func Run() {
	window.Eval("alert('hello world')")
}
