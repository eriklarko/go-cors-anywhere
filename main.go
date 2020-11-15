package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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
	urlToForwardTo, err := getURLToForwardTo(req)
	if err != nil {
		fmt.Fprintf(res, "unable to parse URL: %v", err)
		return
	}

	httputil.NewSingleHostReverseProxy(urlToForwardTo).
		ServeHTTP(res, req)
}

func getURLToForwardTo(req *http.Request) (*url.URL, error) {
	// We expect request in the format
	//   protocol://host:port/url-to-foward to
	// where url-to-foward to can be (non-exhaustive)
	//   * google.com
	//   * https://google.com
	//   * google.com:443
	//
	// If the parameter `req` is a request to http://example.com/hello
	// then `req.URL.Path` is the string "/hello", so we simply
	// take the `req.URL.Path` without the first slash to be the URL
	// to forward to.

	rawURL := strings.TrimPrefix(req.URL.Path, "/")
	return url.Parse(rawURL)
}
