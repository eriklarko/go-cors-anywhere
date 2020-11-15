package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func httpHandler(res http.ResponseWriter, req *http.Request) {
	urlToForwardTo, err := url.Parse("https://icanhazip.com")
	if err != nil {
		fmt.Fprintf(res, "unable to parse URL: %v", err)
		return
	}

	httputil.NewSingleHostReverseProxy(urlToForwardTo).
		ServeHTTP(res, req)
}
