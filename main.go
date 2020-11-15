package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", httpHandler)

	port := 8080

	log.Printf("Starting HTTP server on port %d", port)
	err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		nil,
	)
	if err != nil {
		log.Fatalf("Unable to start HTTP server on port %d: %v\n", port, err)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}
