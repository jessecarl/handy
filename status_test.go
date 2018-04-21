// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jessecarl/handy"
)

func TestServeStatusValidCodes(t *testing.T) {
	s := httptest.NewServer(handy.ServeStatus())

	testCases := []struct {
		code          int
		methods       []string
		wantEmptyBody bool
	}{
		// 2xx
		{http.StatusOK, allMethods, false},
		{http.StatusCreated, allMethods, false},
		{http.StatusAccepted, allMethods, false},
		{http.StatusNonAuthoritativeInfo, allMethods, false},
		{http.StatusNoContent, allMethods, true},
		{http.StatusResetContent, allMethods, false},
		{http.StatusPartialContent, allMethods, false},

		// 4xx
		{http.StatusBadRequest, allMethods, false},
		{http.StatusUnauthorized, allMethods, false},
		{http.StatusPaymentRequired, allMethods, false},
		{http.StatusForbidden, allMethods, false},
		{http.StatusNotFound, allMethods, false},
		{http.StatusMethodNotAllowed, allMethods, false},
		{http.StatusNotAcceptable, allMethods, false},
		{http.StatusProxyAuthRequired, allMethods, false},
		{http.StatusRequestTimeout, allMethods, false},
		{http.StatusConflict, allMethods, false},
		{http.StatusGone, allMethods, false},
		{http.StatusLengthRequired, allMethods, false},
		{http.StatusPreconditionFailed, allMethods, false},
		{http.StatusRequestEntityTooLarge, allMethods, false},
		{http.StatusRequestURITooLong, allMethods, false},
		{http.StatusUnsupportedMediaType, allMethods, false},
		{http.StatusRequestedRangeNotSatisfiable, allMethods, false},
		{http.StatusExpectationFailed, allMethods, false},
		{http.StatusTeapot, allMethods, false},
		{http.StatusUnprocessableEntity, allMethods, false},
		{http.StatusLocked, allMethods, false},
		{http.StatusFailedDependency, allMethods, false},
		{http.StatusUpgradeRequired, allMethods, false},
		{http.StatusPreconditionRequired, allMethods, false},
		{http.StatusTooManyRequests, allMethods, false},
		{http.StatusRequestHeaderFieldsTooLarge, allMethods, false},
		{http.StatusUnavailableForLegalReasons, allMethods, false},

		// 5xx
		{http.StatusInternalServerError, allMethods, false},
		{http.StatusNotImplemented, allMethods, false},
		{http.StatusBadGateway, allMethods, false},
		{http.StatusServiceUnavailable, allMethods, false},
		{http.StatusGatewayTimeout, allMethods, false},
		{http.StatusHTTPVersionNotSupported, allMethods, false},
		{http.StatusVariantAlsoNegotiates, allMethods, false},
		{http.StatusInsufficientStorage, allMethods, false},
		{http.StatusLoopDetected, allMethods, false},
		{http.StatusNotExtended, allMethods, false},
		{http.StatusNetworkAuthenticationRequired, allMethods, false},
	}
	for _, tc := range testCases {
		tc := tc
		for _, method := range tc.methods {
			t.Run(method+": "+http.StatusText(tc.code), func(t *testing.T) {
				t.Parallel() // service should be safe for concurrent use
				req, err := http.NewRequest(method, s.URL+"/"+strconv.Itoa(tc.code), nil)
				if err != nil {
					t.Fatalf("constructing request: %+v", err)
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("doing request: %+v", err)
				}
				defer res.Body.Close()
				if res.StatusCode != tc.code {
					t.Fatalf("request to /%03d returned status code %03d, expected %03d", tc.code, res.StatusCode, tc.code)
				}
				if res.Header.Get("x-status-code") != strconv.Itoa(tc.code) {
					t.Fatalf("request to /%03d returned x-status-code header %s, expected %03d", tc.code, res.Header.Get("x-status-code"), tc.code)
				}
				if res.Header.Get("x-status") != http.StatusText(tc.code) {
					t.Fatalf("request to /%03d returned x-status header %q, expected %q", tc.code, res.Header.Get("x-status"), http.StatusText(tc.code))
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("reading response body: %+v", err)
				}
				var wantBody []byte
				if method == http.MethodHead || tc.wantEmptyBody {
					wantBody = []byte{}
				} else {
					wantBody = []byte(fmt.Sprintf("%03d %s\n", tc.code, http.StatusText(tc.code)))
				}
				if !bytes.Equal(body, wantBody) {
					t.Fatalf("request to /%03d returned body %q, expected %q", tc.code, body, wantBody)
				}
			})
		}
	}
}

// TestServeStatusUnsupportedCodes ensures that the set of
// valid http status codes that work is bounded. Some codes
// may be removed from this set as support is implemented.
func TestServeStatusUnsupportedCodes(t *testing.T) {
	s := httptest.NewServer(handy.ServeStatus())

	testCases := []struct {
		code    int
		methods []string
	}{
		// 1xx
		{http.StatusContinue, allMethods},
		{http.StatusSwitchingProtocols, allMethods},
		{http.StatusProcessing, allMethods},

		// 2xx
		{http.StatusMultiStatus, allMethods},
		{http.StatusAlreadyReported, allMethods},
		{http.StatusIMUsed, allMethods},

		// 3xx
		{http.StatusMultipleChoices, allMethods},
		{http.StatusMovedPermanently, allMethods},
		{http.StatusFound, allMethods},
		{http.StatusSeeOther, allMethods},
		{http.StatusNotModified, allMethods},
		{http.StatusUseProxy, allMethods},
		{http.StatusTemporaryRedirect, allMethods},
		{http.StatusPermanentRedirect, allMethods},

		// not covered
		{0, allMethods},
		{99, allMethods},
		{600, allMethods},
		{1234, allMethods},
	}

	for _, tc := range testCases {
		tc := tc
		for _, method := range tc.methods {
			t.Run(method+": "+http.StatusText(tc.code), func(t *testing.T) {
				t.Parallel() // service should be safe for concurrent use
				req, err := http.NewRequest(method, s.URL+"/"+strconv.Itoa(tc.code), nil)
				if err != nil {
					t.Fatalf("constructing request: %+v", err)
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("doing request: %+v", err)
				}
				defer res.Body.Close()
				if res.StatusCode != http.StatusNotFound {
					t.Fatalf("request to /%03d returned status code %03d, expected %03d", tc.code, res.StatusCode, http.StatusNotFound)
				}
				if res.Header.Get("x-status-code") != "" {
					t.Fatalf("request to /%03d returned x-status-code header %s, expected %s", tc.code, res.Header.Get("x-status-code"), "")
				}
				if res.Header.Get("x-status") != "" {
					t.Fatalf("request to /%03d returned x-status header %q, expected %q", tc.code, res.Header.Get("x-status"), "")
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("reading response body: %+v", err)
				}
				wantBody := []byte("404 page not found\n")
				if !bytes.Equal(body, wantBody) {
					t.Fatalf("request to /%03d returned body %q, expected %q", tc.code, body, wantBody)
				}
			})
		}
	}
}

func TestServeStatusNotCodes(t *testing.T) {
	s := httptest.NewServer(handy.ServeStatus())

	testCases := []struct {
		path    string
		methods []string
	}{
		{"", allMethods},
		{"/", allMethods},
		{"/foo", allMethods},
		{"/200/foo", allMethods},
		{"/bar/200", allMethods},
		{"/200/400", allMethods},
		{"/1/", allMethods},
		{"/602", allMethods},
	}

	for _, tc := range testCases {
		tc := tc
		for _, method := range tc.methods {
			t.Run(method+": "+tc.path, func(t *testing.T) {
				t.Parallel() // service should be safe for concurrent use
				req, err := http.NewRequest(method, s.URL+tc.path, nil)
				if err != nil {
					t.Fatalf("constructing request: %+v", err)
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("doing request: %+v", err)
				}
				defer res.Body.Close()
				if res.StatusCode != http.StatusNotFound {
					t.Fatalf("request to %q returned status code %03d, expected %03d", tc.path, res.StatusCode, http.StatusNotFound)
				}
				if res.Header.Get("x-status-code") != "" {
					t.Fatalf("request to /%q returned x-status-code header %s, expected %s", tc.path, res.Header.Get("x-status-code"), "")
				}
				if res.Header.Get("x-status") != "" {
					t.Fatalf("request to /%q returned x-status header %q, expected %q", tc.path, res.Header.Get("x-status"), "")
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("reading response body: %+v", err)
				}
				wantBody := []byte("404 page not found\n")
				if !bytes.Equal(body, wantBody) {
					t.Fatalf("request to /%q returned body %q, expected %q", tc.path, body, wantBody)
				}
			})
		}
	}
}

var (
	safeMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}
	idempotentMethods = append(safeMethods, http.MethodPut, http.MethodDelete)
	allMethods        = append(idempotentMethods, http.MethodPost, http.MethodPatch)
)

func BenchmarkServeStatus(b *testing.B) {
	s := httptest.NewServer(handy.ServeStatus())
	urls := []string{"/100", "/200", "/300", "/400", "/500", "/foo", "/200/foo"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := http.Get(s.URL + urls[rand.Int()%len(urls)])
		if err != nil {
			b.Fatalf("unexpected error making request: %+v", err)
		}
		res.Body.Close()
	}
}
