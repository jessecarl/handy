// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		{http.StatusOK, allMethods, false},
		{http.StatusCreated, allMethods, false},
		{http.StatusAccepted, allMethods, false},
		{http.StatusNonAuthoritativeInfo, allMethods, false},
		{http.StatusNoContent, allMethods, true},
		{http.StatusResetContent, allMethods, false},
		{http.StatusPartialContent, allMethods, false},
		{http.StatusBadRequest, allMethods, false},
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
					t.Fatalf("request to /%03d returned x-status-code header %03d, expected %03d", tc.code, res.Header.Get("x-status-code"), tc.code)
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
		{http.StatusContinue, allMethods},
		{http.StatusSwitchingProtocols, allMethods},
		{http.StatusProcessing, allMethods},
		{http.StatusMultiStatus, allMethods},
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
					t.Fatalf("request to /%03d returned x-status-code header %03d, expected %03d", tc.code, res.Header.Get("x-status-code"), "")
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
