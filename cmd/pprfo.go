package main

import (
	"log"
	"net/http"
)

func startPprof() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
