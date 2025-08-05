// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"syscall/js"
)

func NewByteBuffer(v js.Value) *ByteBuffer {
	return &ByteBuffer{
		jsVariable: v,
		data:       make([]byte, 0, 1),
	}
}

type ByteBuffer struct {
	jsVariable js.Value
	data       []byte
}

func (b *ByteBuffer) Read() {
	length := b.jsVariable.Call("length").Int()

	if length == 0 {
		b.data = b.data[:0]
		return
	}

	// ensure data is big enough
	if cap(b.data) < length {
		// make data at least 2 times bigger than before
		newSize := 2 * cap(b.data) * (length / cap(b.data))
		b.data = make([]byte, newSize)
	}

	b.data = b.data[:length] // change the length of the slice (an array is not resized)

	js.CopyBytesToGo(b.data, b.jsVariable.Get("buf"))

	b.jsVariable.Call("clear")
}

func (b *ByteBuffer) Data() []byte {
	return b.data
}
