// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package encoder_test

import (
	"math"
	"testing"

	"github.com/elgopher/piweb/internal/audio/encoder"
	"github.com/stretchr/testify/assert"
)

func TestEncoder_Bytes(t *testing.T) {
	commands := encoder.NewCommands(2)

	cmd := commands.AppendCommand(1)
	cmd.AppendU8(2)

	cmd = commands.AppendCommand(3)
	cmd.AppendU8(4)

	t.Run("maxLen bigger than off", func(t *testing.T) {
		bytes := commands.Bytes(16)
		assert.Equal(t, []byte{1, 2, 3, 4}, bytes)
	})

	t.Run("maxLen equal single command", func(t *testing.T) {
		bytes := commands.Bytes(2)
		assert.Equal(t, []byte{1, 2}, bytes)
	})

	t.Run("maxLen slightly more than single command, but less than two commands", func(t *testing.T) {
		bytes := commands.Bytes(3)
		assert.Equal(t, []byte{1, 2}, bytes)
	})
}

func TestEncoder_Remove(t *testing.T) {
	newEncoder := func() *encoder.Commands {
		commands := encoder.NewCommands(2)

		cmd := commands.AppendCommand(1)
		cmd.AppendU8(2)

		cmd = commands.AppendCommand(3)
		cmd.AppendU8(4)
		return commands
	}

	t.Run("should not remove anything", func(t *testing.T) {
		commands := newEncoder()
		commands.Remove(0)
		bytes := commands.Bytes(16)
		assert.Equal(t, []byte{1, 2, 3, 4}, bytes)
	})

	t.Run("should remove first command", func(t *testing.T) {
		encoder := newEncoder()
		encoder.Remove(2)
		bytes := encoder.Bytes(16)
		assert.Equal(t, []byte{3, 4}, bytes)
	})

	t.Run("number of bytes > off", func(t *testing.T) {
		commands := newEncoder()
		commands.Remove(7)
		bytes := commands.Bytes(16)
		assert.Empty(t, bytes)
	})
}

// 150us
func BenchmarkAppendCommand(b *testing.B) {
	commandSize := 16
	commands := encoder.NewCommands(commandSize)
	for i := 0; i < 10000; i++ {
		commands.AppendCommand(byte(i))
	}
	commands.Remove(math.MaxInt)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		for i := 0; i < 10000; i++ {
			commands.AppendCommand(byte(i))
		}
		commands.Remove(1)
		commands.Remove(math.MaxInt)
	}
}
