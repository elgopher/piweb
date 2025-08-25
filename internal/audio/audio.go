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
	window                         = js.Global()
	audioCtx                       = window.Get("audioCtx")
	audioWorkletNode               = window.Get("audioWorkletNode")
	loadSample                     = window.Get("loadSample")
	unloadSample                   = window.Get("unloadSample")
	uint8Array                     = window.Get("Uint8Array")
	audioCommandsSharedArrayBuffer = window.Get("audioCommandsSharedArrayBuffer")
	audioCommands                  = window.Get("audioCommands")
	temporaryBuffer                = window.Get("temporaryBuffer")
	sendCommands                   = window.Get("sendCommands")
	commands                       = encoder.NewCommands(19)
)

func NewBackend() *Backend {
	return &Backend{}
}

type Backend struct{}

func (b *Backend) LoadSample(sample *piaudio.Sample) {
	id := getPointerAddr(sample)

	data := uint8Array.New(sample.Len())
	js.CopyBytesToJS(data, int8SliceToByteSlice(sample.Data()))

	loadSample.Invoke(id, data, sample.SampleRate())
}

func getSampleID(sample *piaudio.Sample) uint32 {
	return uint32(getPointerAddr(sample)) // TODO
}

func getPointerAddr(sample *piaudio.Sample) uintptr {
	return uintptr(unsafe.Pointer(sample))
}

func int8SliceToByteSlice(b []int8) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(b))), len(b))
}

func (b *Backend) UnloadSample(sample *piaudio.Sample) {
	unloadSample.Invoke(getPointerAddr(sample))
}

const audioBufferSizeInSeconds = 0.02 // 20ms

func scheduledTime(delay float64) float64 {
	return delay + piaudio.Time + audioBufferSizeInSeconds
}

func (b *Backend) SetSample(ch piaudio.Chan, sample *piaudio.Sample, offset int, delay float64) {
	cmd := commands.AppendCommand(cmdSetSample)
	cmd.AppendU8(ch)
	cmd.AppendU32(getSampleID(sample))
	cmd.AppendI32(int32(offset))
	cmd.AppendF64(scheduledTime(delay))
}

var loopTypeMapping = map[piaudio.LoopType]byte{
	piaudio.LoopNone:    0,
	piaudio.LoopForward: 1,
}

func (b *Backend) SetLoop(ch piaudio.Chan, start int, length int, loopType piaudio.LoopType, delay float64) {
	cmd := commands.AppendCommand(cmdSetLoop)
	cmd.AppendU8(ch)
	cmd.AppendI32(int32(start))
	cmd.AppendI32(int32(start + length - 1))
	cmd.AppendU8(loopTypeMapping[loopType]) // LoopNone if loopType is invalid
	cmd.AppendF64(scheduledTime(delay))
}

func (b *Backend) ClearChan(ch piaudio.Chan, delay float64) {
	cmd := commands.AppendCommand(cmdClearChan)
	cmd.AppendU8(ch)
	cmd.AppendF64(scheduledTime(delay))
}

func (b *Backend) SetPitch(ch piaudio.Chan, pitch float64, delay float64) {
	cmd := commands.AppendCommand(cmdSetPitch)
	cmd.AppendU8(ch)
	cmd.AppendF64(pitch)
	cmd.AppendF64(scheduledTime(delay))
}

func (b *Backend) SetVolume(ch piaudio.Chan, vol float64, delay float64) {
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
		sendCommands.Invoke()
		commands.Remove(len(bytes))
	}
}
