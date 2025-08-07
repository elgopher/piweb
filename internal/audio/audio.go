// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	_ "embed"

	"github.com/elgopher/pi/piaudio"
)

type Backend struct{}

func (b *Backend) LoadSample(sample *piaudio.Sample) {
	// TODO postMessage with simple buffer (not shared) - just copy data
}

func (b *Backend) UnloadSample(sample *piaudio.Sample) {
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
