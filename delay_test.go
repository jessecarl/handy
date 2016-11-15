// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jessecarl/handy"
)

func TestServeWithDelay(t *testing.T) {
	testCases := []struct {
		name      string
		method    string
		path      string
		wantPath  string
		wantDelay time.Duration
	}{
		{"no delay, GET",
			http.MethodGet,
			"/foo",
			"/foo",
			0,
		},
		{"no delay, no path, GET",
			http.MethodGet,
			"/",
			"/",
			0,
		},
		{"no delay, POST",
			http.MethodPost,
			"/bar/baz/",
			"/bar/baz/",
			0,
		},
		{"1s delay, PUT",
			http.MethodPut,
			"/foo/1s",
			"/foo",
			time.Second,
		},
		{"10ms delay, HEAD",
			http.MethodHead,
			"/10ms/",
			"/", // always have at least one slash
			time.Millisecond * 10,
		},
		{"300ms delay, GET",
			http.MethodGet,
			"/baz/bar/300ms/",
			"/baz/bar",
			time.Millisecond * 300,
		},
		{"1ms delay, GET",
			http.MethodGet,
			"/foo/bar/baz/1ms",
			"/foo/bar/baz",
			time.Millisecond,
		},
		{"100µs delay, GET",
			http.MethodGet,
			"/foo/bar/baz/100µs",
			"/foo/bar/baz",
			time.Microsecond * 100,
		},
		{"100µs delay 'us', GET",
			http.MethodGet,
			"/foo/bar/baz/100us",
			"/foo/bar/baz",
			time.Microsecond * 100,
		},
		{"10µs delay, GET",
			http.MethodGet,
			"/foo/bar/baz/10µs",
			"/foo/bar/baz",
			time.Microsecond * 10,
		},
		{"1µs delay, GET",
			http.MethodGet,
			"/foo/bar/baz/1µs",
			"/foo/bar/baz",
			time.Microsecond,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var start time.Time
			s := httptest.NewServer(handy.ServeWithDelay(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotDelay := time.Since(start)
				if r.URL.Path != tc.wantPath {
					t.Fatalf("handler call to %q, expected %q", r.URL.Path, tc.wantPath)
				}
				delayError := (float64(gotDelay) - float64(tc.wantDelay)) / float64(tc.wantDelay)
				if delayError > 0.05 && delayError < -0.05 {
					t.Fatalf("handler delay took %v, expected %v, more than %v delay", gotDelay, tc.wantDelay, delayError)
				}
				w.Write([]byte("delayed handler"))
			})))

			req, err := http.NewRequest(tc.method, s.URL+tc.path, nil)
			if err != nil {
				t.Fatalf("constructing request: %+v", err)
			}
			start = time.Now()
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("doing request: %+v", err)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("reading response body: %+v", err)
			}
			if !bytes.Equal(body, []byte("delayed handler")) && tc.method != http.MethodHead {
				t.Fatalf("request to %q returned body %q, expected %q", tc.path, body, []byte("delayed handler"))
			}
		})
	}
}
