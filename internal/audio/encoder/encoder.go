// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package encoder

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Commands struct {
	buf         []byte
	off         int
	commandSize int
}

// NewCommands creates an Commands, which encodes equally-sized commands.
// commandSize is a length of each command.
func NewCommands(commandSize int) *Commands {
	return &Commands{
		buf:         make([]byte, commandSize),
		commandSize: commandSize,
	}
}

// Bytes zwraca komendy mieszczące się w maxLen. Zwrócony bufor zawsze
// zawiera pełne komendy.
func (e *Commands) Bytes(maxLen int) []byte {
	if maxLen > e.off {
		maxLen = e.off
	}
	i := (maxLen / e.commandSize) * e.commandSize
	return e.buf[:i]
}

func (e *Commands) Remove(bytes int) {
	if bytes == 0 {
		return
	}

	if bytes >= e.off {
		e.off = 0
		return
	}

	copy(e.buf, e.buf[bytes:e.off])
	e.off -= bytes
}

func (e *Commands) AppendCommand(op byte) Command {
	e.ensure(e.commandSize)

	off := e.off
	e.off += e.commandSize

	e.buf[off] = op

	return Command{
		buf:  e.buf[off : off+e.commandSize],
		size: e.commandSize,
		off:  1,
	}
}

func (e *Commands) ensure(n int) {
	if len(e.buf)-e.off < n {
		newBuf := make([]byte, e.off+n+e.off/2)
		copy(newBuf, e.buf[:e.off])
		e.buf = newBuf
	}
}

type Command struct {
	off  int
	buf  []byte
	size int
}

func (o *Command) ensure(bytes int) {
	newSize := o.off + bytes
	if newSize > len(o.buf) {
		msg := fmt.Sprintf("Cannot encode more bytes in the command. "+
			"The command size is %d, but program was trying to encode "+
			"%d bytes",
			o.size, newSize)
		panic(msg)
	}
}

func (o *Command) AppendU8(v uint8) {
	o.ensure(1)
	o.buf[o.off] = v
	o.off += 1
}

func (o *Command) AppendI32(v int32) {
	o.ensure(4)
	binary.LittleEndian.PutUint32(o.buf[o.off:], uint32(v))
	o.off += 4
}

func (o *Command) AppendU32(v uint32) {
	o.ensure(4)
	binary.LittleEndian.PutUint32(o.buf[o.off:], v)
	o.off += 4
}

func (o *Command) AppendF32(v float32) {
	o.ensure(4)
	binary.LittleEndian.PutUint32(o.buf[o.off:], math.Float32bits(v))
	o.off += 4
}

func (o *Command) AppendF64(v float64) {
	o.ensure(8)
	binary.LittleEndian.PutUint64(o.buf[o.off:], math.Float64bits(v))
	o.off += 8
}

func (o *Command) Bytes() []byte {
	return o.buf
}
