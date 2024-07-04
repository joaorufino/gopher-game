package main

import (
	"log"
	"mime"
	"net/http"
)

func main() {
	if err := mime.AddExtensionType(".wasm", "application/wasm"); err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/", noCacheHandler(fs))

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// noCacheHandler wraps an http.Handler to add no-cache headers to the response.
func noCacheHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
