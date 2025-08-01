// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package window

import "syscall/js"

var window = js.Global()

func Eval(code string) {
	window.Call("eval", code)
}
