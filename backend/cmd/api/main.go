package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	fmt.Println("API Server initial running on :8080...")
	http.ListenAndServe(":8080", nil)
}
