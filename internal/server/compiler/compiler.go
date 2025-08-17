// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package compiler

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var wasmFilePrefix = []byte("\x00asm")

func BuildMainWasm() ([]byte, error) {
	outputPath := filepath.Join(tmpDir, "main.wasm")

	// Run `go build .`
	args := []string{"build", "-o", outputPath, "."}
	log.Print("Running go ", strings.Join(args, " "))

	cmd := exec.Command("go", args...)

	cmd.Env = os.Environ()
	cmd.Env = slices.DeleteFunc(cmd.Env, func(v string) bool {
		return strings.HasPrefix(v, "GOOS=") || strings.HasPrefix(v, "GOARCH=")
	})
	cmd.Env = append(cmd.Env, "GOOS=js", "GOARCH=wasm")

	cmd.Dir = "."

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Problem running go build: \n%s%w", stderr.String(), err)
	}

	if stdout.Len() > 0 {
		log.Print(stdout.String())
	}
	if stderr.Len() > 0 {
		log.Print(stderr.String())
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("Problem reading file %s: %w", outputPath, err)
	}

	// ensure real WASM was created
	if !bytes.HasPrefix(content, wasmFilePrefix) {
		workdir, _ := os.Getwd()
		return nil, fmt.Errorf("Package main not found in the current working directory: %s", workdir)
	}

	return content, nil
}

var tmpDir = func() string {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal("Error creating temporary directory")
	}
	return tmp
}()
