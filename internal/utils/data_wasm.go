//go:build js && wasm
// +build js,wasm

package utils

import (
	"errors"
	"log"
	"syscall/js"
)

// DataCallback defines the type of the callback function.
type DataCallback func(data []byte) error

// LoadData loads data from the specified path using the fetchData JavaScript function and calls the provided callback with the data.
func LoadData(path string, callback DataCallback) error {
	// Create a channel to signal when the fetch is complete
	done := make(chan struct{})

	// Define a Go function to be called by the JavaScript fetchData function.
	jsCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(done)
		if len(args) < 1 {
			log.Println("Data not provided")
			return nil
		}
		data := args[0].String()
		log.Printf("Fetched data: %s", data) // Print the fetched data as a string
		err := callback([]byte(data))
		if err != nil {
			log.Printf("Failed to handle data: %v", err)
		}
		return nil
	})
	defer jsCallback.Release()

	// Call the JavaScript fetchData function with the path and the Go callback function.
	result := js.Global().Call("fetchData", path, jsCallback)
	if result.IsNull() || result.IsUndefined() {
		return errors.New("fetchData call failed")
	}

	log.Println("fetchData call initiated")

	// Wait for the fetch operation to complete
	<-done
	return nil
}
