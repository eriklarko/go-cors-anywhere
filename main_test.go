package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			require.NoError(t, err, "failed creating test request")

			actual, err := getURLToForwardTo(req)
			require.NoError(t, err, "unable to read URL to forward to from test request")

			assert.Equal(t, testCase.expected, actual.String())
		})
	}
}

func TestCORSHeadersAdded(t *testing.T) {
	testCases := map[string]struct {
		// headers in the request from the client
		clientRequestHeaders http.Header

		// headers in the response from the server we're proxying
		serverResponseHeaders http.Header

		// headers returned to the client
		expectedHeaders http.Header
	}{
		"Adds expected CORS headers": {
			clientRequestHeaders: http.Header{
				"Access-Control-Request-Method":  []string{"GET"},
				"Access-Control-Request-Headers": []string{"X-Foo", "X-Bar"},
			},
			expectedHeaders: http.Header{
				"Access-Control-Allow-Origin":  []string{"*"},
				"Access-Control-Allow-Methods": []string{"GET"},
				"Access-Control-Allow-Headers": []string{"X-Foo", "X-Bar"},
			},
		},
		"Doesn't add allow-methods if not specified": {
			clientRequestHeaders: http.Header{},
			expectedHeaders: http.Header{
				"Access-Control-Allow-Origin":  []string{"*"},
				"Access-Control-Allow-Methods": nil,
			},
		},
		"Doesn't add allow-headers if not specified": {
			clientRequestHeaders: http.Header{},
			expectedHeaders: http.Header{
				"Access-Control-Allow-Origin":  []string{"*"},
				"Access-Control-Allow-Headers": nil,
			},
		},

		"overwrites Allow-Origin from other server": {
			serverResponseHeaders: http.Header{
				"Access-Control-Allow-Origin": []string{"example.com"},
			},
			expectedHeaders: http.Header{
				"Access-Control-Allow-Origin": []string{"*"},
			},
		},

		"exposes non-CORS headers to the client": {
			serverResponseHeaders: http.Header{
				"X-Foo":                            []string{"foo-val"},
				"X-Bar":                            []string{"bar-val"},
				"Access-Control-Allow-Credentials": []string{"true"},
			},
			expectedHeaders: http.Header{
				"X-Foo":                            []string{"foo-val"},
				"X-Bar":                            []string{"bar-val"},
				"Access-Control-Expose-Headers":    []string{"X-Foo", "X-Bar"},
				"Access-Control-Allow-Credentials": nil,
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {

			req, err := http.NewRequest("GET", "foo.bar", nil)
			require.NoError(t, err, "failed creating test request")

			// add headers to incoming request
			req.Header = testCase.clientRequestHeaders

			actualHeaders := testCase.serverResponseHeaders.Clone()
			if actualHeaders == nil {
				actualHeaders = make(http.Header)
			}

			addCORSHeaders(actualHeaders, req)

			// log headers to make it easier to debug tests
			t.Logf("Expected headers: %v\n", testCase.expectedHeaders)
			t.Logf("Actual headers: %v\n", actualHeaders)

			for headerName, expectedValues := range testCase.expectedHeaders {
				actualValues := actualHeaders.Values(headerName)

				assert.Equal(t,
					len(expectedValues), len(actualValues),
					fmt.Sprintf(
						"Header %s didn't have the correct number of values",
						headerName,
					),
				)

				for _, expectedValue := range expectedValues {
					assert.Contains(t,
						actualValues,
						expectedValue,
						fmt.Sprintf(
							"Header %s didn't contain expected value %s",
							headerName,
							expectedValue,
						),
					)
				}
			}
		})
	}
}
