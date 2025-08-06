// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package server

import (
	"embed"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/elgopher/piweb/internal/server/compiler"
	"github.com/elgopher/piweb/internal/server/indexhtml"
	"github.com/elgopher/piweb/internal/server/wasmexecjs"
)

//go:embed "html/*"
var embeddedHtmlDir embed.FS

func newHandler() *Handler {
	return &Handler{
		wasmExecJS: wasmexecjs.Get(),
	}
}

type Handler struct {
	wasmExecJS []byte
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	file := strings.TrimPrefix(request.RequestURI, "/")

	if file == "wasm_exec.js" {
		writeContent(file, h.wasmExecJS, writer)
		return
	}

	if file == "main.wasm" {
		wasm, err := compiler.BuildMainWasm()
		if err != nil {
			log.Print(err)
			writer.WriteHeader(500)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		writeContent(file, wasm, writer)
		return
	}

	if file == "" {
		file = "index.html"
	}

	workdir, _ := os.Getwd()
	workdirFile := filepath.Join(workdir, file)
	content, err := os.ReadFile(workdirFile)
	if err == nil {
		writeContent(file, content, writer)
		return
	}

	content, err = embeddedHtmlDir.ReadFile("html/" + file)
	if err == nil {
		writeContent(file, content, writer)
		return
	}

	writer.WriteHeader(404)
	_, _ = writer.Write([]byte(http.StatusText(http.StatusNotFound)))
}

func writeContent(file string, content []byte, writer http.ResponseWriter) {
	if file == "index.html" {
		content = indexhtml.PutScripts(content)
	}

	contentType := mime.TypeByExtension(filepath.Ext(file))
	writer.Header().Set("Content-Type", contentType)

	// Enable SharedArrayBuffer:
	writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

	_, _ = writer.Write(content)
}
