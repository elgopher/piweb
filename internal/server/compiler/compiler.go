// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package compiler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var wasmFilePrefix = []byte("\x00asm")

func BuildMainWasm(goBuildCommand string) ([]byte, error) {
	outputPath := filepath.Join(tmpDir, "main.wasm")

	// Run `go build .`
	args, err := parseGoBuildCommand(goBuildCommand, outputPath)
	if err != nil {
		return nil, err
	}
	log.Print("Running ", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)

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
		return nil, fmt.Errorf("Problem reading output file %s: %w", outputPath, err)
	}

	// ensure real WASM was created
	if !bytes.HasPrefix(content, wasmFilePrefix) {
		workdir, _ := os.Getwd()
		return nil, fmt.Errorf("Package main not found in the current working directory: %s", workdir)
	}

	return content, nil
}

func parseGoBuildCommand(command string, outputPath string) ([]string, error) {
	tmpl, err := template.New("").Parse(command)
	if err != nil {
		return nil, fmt.Errorf("problem parsing go build command %s: %w", command, err)
	}

	var evaluatedCommand bytes.Buffer
	err = tmpl.Execute(&evaluatedCommand, struct {
		Output string
	}{
		Output: outputPath,
	})
	if err != nil {
		return nil, fmt.Errorf("problem evaluating go build command %s: %w", command, err)
	}

	args := strings.Split(evaluatedCommand.String(), " ")

	return args, nil
}

var tmpDir = func() string {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		log.Fatal("Error creating temporary directory")
	}
	return tmp
}()
