// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package server

import (
	"errors"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	file := strings.TrimPrefix(request.RequestURI, "/")

	if file == "release.zip" {
		zip, err := ReleaseZip()
		if err != nil {
			writerInternalServerError(err, writer)
			return
		}

		writeContent(file, zip, writer)
		return
	}

	if file == "" {
		file = "index.html"
	}

	content, err := GetFile(file)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writer.WriteHeader(404)
			_, _ = writer.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}

		writerInternalServerError(err, writer)
		return
	}

	writeContent(file, content, writer)
}

func writeContent(file string, content []byte, writer http.ResponseWriter) {
	contentType := mime.TypeByExtension(filepath.Ext(file))
	writer.Header().Set("Content-Type", contentType)

	// Enable SharedArrayBuffer:
	writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")

	_, _ = writer.Write(content)
}

func writerInternalServerError(err error, writer http.ResponseWriter) {
	log.Print(err)
	writer.WriteHeader(http.StatusInternalServerError)
	_, _ = writer.Write([]byte(err.Error()))
}
