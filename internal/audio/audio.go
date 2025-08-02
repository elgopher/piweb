// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"github.com/elgopher/pi/piaudio"
)

type Backend struct{}

func (b *Backend) LoadSample(sample *piaudio.Sample) {}

func (b *Backend) UnloadSample(sample *piaudio.Sample) {}

func (b *Backend) SetSample(ch piaudio.Chan, sample *piaudio.Sample, offset int, delay float64) {}

func (b *Backend) SetLoop(_ piaudio.Chan, start int, length int, loopType piaudio.LoopType, delay float64) {
}

func (b *Backend) SetPitch(_ piaudio.Chan, pitch float64, delay float64) {}

func (b *Backend) SetVolume(_ piaudio.Chan, vol float64, delay float64) {}

func (b *Backend) ClearChan(ch piaudio.Chan, delay float64) {}
