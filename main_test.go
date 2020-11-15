package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLToForwardTo(t *testing.T) {
	testCases := map[string]struct {
		// the URL we want to forward to
		url string

		// the result we expect from the getURLToForwardTo function
		expected string
	}{
		"only host name": {
			url:      "foo.com",
			expected: "http://foo.com",
		},
		"with protocol and port": {
			url:      "http://foo.com:1337",
			expected: "http://foo.com:1337",
		},
		"with protocol": {
			url:      "http://foo.com",
			expected: "http://foo.com",
		},
		"with port": {
			url:      "foo.com:1337",
			expected: "http://foo.com:1337",
		},
		"adds https if port 443": {
			url:      "foo.com:443",
			expected: "https://foo.com:443",
		},
		"overwrites protocol if port 443": {
			url:      "http://foo.com:443",
			expected: "https://foo.com:443",
		},
	}

	serverURL := "https://example.com:8080"
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(
				"GET", fmt.Sprintf("%s/%s", serverURL, testCase.url),
				nil,
			)
			assert.NoError(t, err, "failed creating test request")

			actual, err := getURLToForwardTo(req)
			assert.NoError(t, err, "unable to read URL to forward to from test request")

			assert.Equal(t, testCase.expected, actual.String())
		})
	}
}
