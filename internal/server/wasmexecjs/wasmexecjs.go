// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package wasmexecjs

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Get() []byte {
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		log.Fatalf("Failed to get GOROOT: %v", err)
	}

	goRoot := strings.Trim(string(out), "\n")
	goRoot = filepath.Clean(goRoot)

	wasmExecPath := filepath.Join(goRoot, "lib", "wasm", "wasm_exec.js")

	bytes, err := os.ReadFile(wasmExecPath)
	if err != nil {
		log.Fatalf("Failed to read wasm_exec.js: %v", err)
	}

	return bytes
}
