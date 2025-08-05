// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	_ "embed"
	"syscall/js"
	"unsafe"

	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/piweb/internal/audio/encoder"
)

var (
	window          = js.Global()
	audioCtx        = window.Get("audioCtx")
	loadSample      = window.Get("loadSample")
	unloadSample    = window.Get("unloadSample")
	uint8Array      = window.Get("Uint8Array")
	audioCommands   = window.Get("audioCommands")
	temporaryBuffer = window.Get("temporaryBuffer")
	sendCommands    = window.Get("sendCommands")
	commands        = encoder.NewCommands(19)
)

func NewBackend() *Backend {
	return &Backend{
		sampleIDs: make(map[*piaudio.Sample]uint32),
	}
}

type Backend struct {
	sampleIDs    map[*piaudio.Sample]uint32
	nextSampleID uint32
}

func (b *Backend) LoadSample(sample *piaudio.Sample) {
	id, ok := b.sampleIDs[sample]
	if !ok {
		b.nextSampleID++
		id = b.nextSampleID
		b.sampleIDs[sample] = id
	}

	data := uint8Array.New(sample.Len())
	js.CopyBytesToJS(data, int8SliceToByteSlice(sample.Data()))

	loadSample.Invoke(id, data, sample.SampleRate())
}

// int8SliceToByteSlice converts an []int8 to []byte without allocating.
// int8 and byte (alias of uint8) share the same 1-byte representation on all
// Go architectures, so reinterpreting the slice using unsafe is portable and
// avoids the cost of copying audio data.
func int8SliceToByteSlice(b []int8) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(b))), len(b))
}

func (b *Backend) UnloadSample(sample *piaudio.Sample) {
	unloadSample.Invoke(b.sampleIDs[sample])
	delete(b.sampleIDs, sample)
}

const audioBufferSizeInSeconds = 0.02 // 20ms

func scheduledTime(delay float64) float64 {
	return delay + piaudio.Time + audioBufferSizeInSeconds
}

func fixChan(p piaudio.Chan) piaudio.Chan {
	return p & (piaudio.Chan1 | piaudio.Chan2 | piaudio.Chan3 | piaudio.Chan4)
}

func (b *Backend) SetSample(ch piaudio.Chan, sample *piaudio.Sample, offset int, delay float64) {
	ch = fixChan(ch)
	if ch == 0 {
		return
	}

	cmd := commands.AppendCommand(cmdSetSample)
	cmd.AppendU8(ch)
	cmd.AppendU32(b.sampleIDs[sample])
	cmd.AppendI32(int32(offset))
	cmd.AppendF64(scheduledTime(delay))
}

var loopTypeMapping = map[piaudio.LoopType]byte{
	piaudio.LoopNone:    0,
	piaudio.LoopForward: 1,
}

func (b *Backend) SetLoop(ch piaudio.Chan, start int, length int, loopType piaudio.LoopType, delay float64) {
	ch = fixChan(ch)
	if ch == 0 {
		return
	}

	cmd := commands.AppendCommand(cmdSetLoop)
	cmd.AppendU8(ch)
	cmd.AppendI32(int32(start))
	cmd.AppendI32(int32(start + length - 1))
	cmd.AppendU8(loopTypeMapping[loopType]) // LoopNone if loopType is invalid
	cmd.AppendF64(scheduledTime(delay))
}

func (b *Backend) ClearChan(ch piaudio.Chan, delay float64) {
	ch = fixChan(ch)
	if ch == 0 {
		return
	}

	cmd := commands.AppendCommand(cmdClearChan)
	cmd.AppendU8(ch)
	cmd.AppendF64(scheduledTime(delay))
}

func (b *Backend) SetPitch(ch piaudio.Chan, pitch float64, delay float64) {
	ch = fixChan(ch)
	if ch == 0 {
		return
	}

	cmd := commands.AppendCommand(cmdSetPitch)
	cmd.AppendU8(ch)
	cmd.AppendF64(pitch)
	cmd.AppendF64(scheduledTime(delay))
}

func (b *Backend) SetVolume(ch piaudio.Chan, vol float64, delay float64) {
	ch = fixChan(ch)
	if ch == 0 {
		return
	}

	cmd := commands.AppendCommand(cmdSetVolume)
	cmd.AppendU8(ch)
	cmd.AppendF64(vol)
	cmd.AppendF64(scheduledTime(delay))
}

func UpdateTime() {
	piaudio.Time = audioCtx.Get("currentTime").Float()
}

func SendCommands() {
	freeSpace := audioCommands.Call("freeSpace").Int()
	bytes := commands.Bytes(freeSpace)
	length := len(bytes)
	if length > 0 {
		js.CopyBytesToJS(temporaryBuffer, bytes)
		window.Set("temporaryBufferLength", length)
		sendCommands.Invoke() // consider calling sendCommands more often (once per 128 commands?)
		commands.Remove(len(bytes))
	}
}
