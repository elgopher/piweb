// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package window

import "syscall/js"

var window = js.Global()

func Eval(code string) {
	window.Call("eval", code)
}

func Set(p string, x any) {
	window.Set(p, x)
}

func NewObject(args ...any) js.Value {
	return window.Get("Object").New(args...)
}
