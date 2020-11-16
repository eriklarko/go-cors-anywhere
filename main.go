package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
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

	addCORSHeaders(res.Header(), req)
}

func getURLToForwardTo(req *http.Request) (*url.URL, error) {
	// We expect requests in the format
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

	// golang's url.Parse doesn't work great if the URL passed to it
	// doesn't include a protocol. To work around this we add "http" if the URL
	// doesn't appear to have a protocol.
	rawURLWithProtocol := addProtocolIfNotPresent("http", rawURL)

	url, err := url.Parse(rawURLWithProtocol)
	if err != nil {
		return nil, err
	}

	// We want to ensure https if the port is 443
	if url.Port() == "443" {
		url.Scheme = "https"
	}

	return url, nil
}

func addProtocolIfNotPresent(protocol, url string) string {
	// matches strings starting with at least one alphanumeric
	// character followed by "://"
	//
	// this variable could be moved out of this function but until I'm able to
	// measure the performance implications of compiling the regex on every
	// request it makes more sense to keep cohesion strong rather than guess
	// where the performance bottlenecks will be
	hasProtocol := regexp.MustCompile(`^[[:alpha:]]+://`)

	if hasProtocol.MatchString(url) {
		return url
	} else {
		return fmt.Sprintf("%s://%s", protocol, url)
	}
}

func addCORSHeaders(responseHeaders http.Header, req *http.Request) {
	// req is the request from the client

	// delete CORS headers set by the server we're proxying
	for headerName := range responseHeaders {
		if isAccessControlHeader(headerName) {
			responseHeaders.Del(headerName)
		}
	}

	// expose headers so that the browser allows the client to read them
	for headerName := range responseHeaders {
		if !isAccessControlHeader(headerName) {
			responseHeaders.Add("Access-Control-Expose-Headers", headerName)
		}
	}

	// allow requests from any origin. this could be better
	responseHeaders.Set("Access-Control-Allow-Origin", "*")

	// add allow-methods if methods aare specified by the client request
	if method := req.Header.Get("Access-Control-Request-Method"); method != "" {
		responseHeaders.Set("Access-Control-Allow-Methods", method)
	}

	// add allow-headers if headers are specified by the client request
	if headers := req.Header.Values("Access-Control-Request-Headers"); len(headers) > 0 {
		for _, header := range headers {
			responseHeaders.Add("Access-Control-Allow-Headers", header)
		}
	}
}

func isAccessControlHeader(headerName string) bool {
	return strings.HasPrefix(
		strings.ToLower(headerName),
		"access-control-",
	)
}
