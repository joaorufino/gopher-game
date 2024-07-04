//go:build !js || !wasm
// +build !js !wasm

package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DataCallback defines the type of the callback function.
type DataCallback func(data []byte) error

// LoadData loads data from a file in non-WASM builds.
func LoadData(path string, callback DataCallback) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	return callback(data)
}

// loadDataFromURL loads data from a URL for local testing.
func loadDataFromURL(url string, callback DataCallback) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	return callback(data)
}
