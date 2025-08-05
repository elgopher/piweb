// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	_ "embed"
	"syscall/js"
	"unsafe"

	"github.com/elgopher/pi/piaudio"
)

var window = js.Global()

var audioWorkletNode = window.Get("audioWorkletNode")

type Backend struct{}

func (b *Backend) LoadSample(sample *piaudio.Sample) {
	msg := window.Get("Object").New()
	msg.Set("type", "loadSample")
	msg.Set("id", getPointerAddr(sample))
	uint8array := window.Get("Uint8Array").New(sample.Len())
	js.CopyBytesToJS(uint8array, int8SliceToByteSlice(sample.Data()))
	msg.Set("buffer", uint8array)
	msg.Set("rate", sample.SampleRate())

	audioWorkletNode.Get("port").Call("postMessage", msg)
}

func getPointerAddr(sample *piaudio.Sample) uintptr {
	return uintptr(unsafe.Pointer(sample))
}

func int8SliceToByteSlice(b []int8) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(b))), len(b))
}

func (b *Backend) UnloadSample(sample *piaudio.Sample) {
	msg := window.Get("Object").New()
	msg.Set("type", "unloadSample")
	msg.Set("id", getPointerAddr(sample))

	audioWorkletNode.Get("port").Call("postMessage", msg)

}

func (b *Backend) SetSample(ch piaudio.Chan, sample *piaudio.Sample, offset int, delay float64) {
	// TODO write message to ByteBuffer
	// Use double-buffer + dynamic switch
	// When buffer is full, create a new 2x bigger buffer.
	// Send message to worklet about new buffer. Write commands to this new buffer.
	// AUdio worklet reads all commands from old buffer until all commands
	// are read. Then it switches to the new buffer received.

	// TODO SHOULD SetSample directly update JS memory? Or should it modify Go memory
	// and update the buffer when Update finishes?
}

func (b *Backend) SetLoop(_ piaudio.Chan, start int, length int, loopType piaudio.LoopType, delay float64) {
}

func (b *Backend) SetPitch(_ piaudio.Chan, pitch float64, delay float64) {}

func (b *Backend) SetVolume(_ piaudio.Chan, vol float64, delay float64) {}

func (b *Backend) ClearChan(ch piaudio.Chan, delay float64) {}
