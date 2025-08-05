// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func ReleaseZip() ([]byte, error) {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)

	files, err := listReleaseFiles()
	if err != nil {
		return nil, fmt.Errorf("error listing release files: %w", err)
	}

	for _, file := range files {
		content, err := GetFile(file, *ReleaseGoBuild)
		if err != nil {
			return nil, fmt.Errorf("getting file %s failed: %w", file, err)
		}
		if err := addFileToZip(writer, file, content); err != nil {
			return nil, fmt.Errorf("adding file to zip failed: %v", err)
		}
	}

	writer.Close()

	return buf.Bytes(), nil
}

func listReleaseFiles() ([]string, error) {
	files := []string{
		"main.wasm",
		"wasm_exec.js",
	}

	embeddedEntry, err := embeddedHtmlDir.ReadDir("html")
	if err != nil {
		return nil, fmt.Errorf("error reading embedded html directory: %w", err)
	}
	for _, entry := range embeddedEntry {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	workdirEntries, err := os.ReadDir(*HtmlDir)
	if err != nil {
		return nil, fmt.Errorf("error reading working directory: %w", err)
	}
	for _, entry := range workdirEntries {
		if filepath.Ext(entry.Name()) != ".go" && !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	files = removeDuplicates(files)

	return files, nil
}

func addFileToZip(zipWriter *zip.Writer, filename string, content []byte) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = writer.Write(content)
	if err != nil {
		return fmt.Errorf("writing zip file failed: %w", err)
	}

	return nil
}

func removeDuplicates(input []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))

	for _, s := range input {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}

	return result
}
