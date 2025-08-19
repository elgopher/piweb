// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package server

import (
	"embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/elgopher/piweb/internal/server/compiler"
	"github.com/elgopher/piweb/internal/server/indexhtml"
	"github.com/elgopher/piweb/internal/server/wasmexecjs"
)

//go:embed "html/*"
var embeddedHtmlDir embed.FS

var wasmExecJS = wasmexecjs.Get()

var ErrNotFound = errors.New("not found")

var (
	GoBuild        *string
	ReleaseGoBuild *string
)

func GetFile(file string, goBuild string) ([]byte, error) {
	if file == "wasm_exec.js" {
		return wasmExecJS, nil
	}

	if file == "main.wasm" {
		return compiler.BuildMainWasm(goBuild)
	}

	workdir, _ := os.Getwd()
	workdirFile := filepath.Join(workdir, file)
	content, err := os.ReadFile(workdirFile)
	if err == nil {
		return process(file, content), nil
	}

	content, err = embeddedHtmlDir.ReadFile("html/" + file)
	if err == nil {
		return process(file, content), nil
	}

	return nil, ErrNotFound
}

func process(file string, content []byte) []byte {
	if file == "index.html" {
		content = indexhtml.PutScripts(content)
	}

	return content
}
