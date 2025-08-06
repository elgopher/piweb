// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package indexhtml

import (
	"bytes"
	_ "embed"
)

//go:embed "scripts.html"
var scripts []byte

func PutScripts(content []byte) []byte {
	return bytes.Replace(content, []byte("$$SCRIPTS$$"), scripts, 1)
}
