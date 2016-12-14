// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/jessecarl/handy"
)

func TestPace_inOrder(t *testing.T) {
	testCases := []struct {
		name  string
		count int
		delay time.Duration
	}{
		{"none", 0, 0},
		{"one", 1, 10 * time.Millisecond},
		{"two", 2, 10 * time.Millisecond},
		{"many", 15, 15 * time.Millisecond},
		{"more than many", 150, time.Millisecond},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var (
				i  int
				mu sync.Mutex
			)

			s := httptest.NewServer(handy.Pace(
				tc.count,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer w.Write([]byte("a response"))
					func() {
						mu.Lock()
						defer mu.Unlock()
						j, err := strconv.Atoi(r.FormValue("i"))
						if err != nil {
							t.Fatalf("unexpected error parsing request count: %+v", err)
						}
						if j < i {
							t.Fatalf("processing request %d before request %d, out of order", j, i)
						}
						i = j
					}()
					<-time.After(tc.delay * 2)
				}),
			))

			tick := make(chan int)
			var wg sync.WaitGroup
			go func() {
				for i := range tick {
					go func(i int) {
						defer wg.Done()
						res, err := http.Get(fmt.Sprintf("%s/?i=%d", s.URL, i))
						if err != nil {
							t.Fatalf("unexpected error getting result: %+v", err)
						}
						defer res.Body.Close()
					}(i)
				}
			}()

			for i := 0; i < 3*tc.count; i++ {
				wg.Add(1)
				<-time.After(tc.delay)
				tick <- i
			}
			wg.Wait()
			close(tick)
		})
	}
}

func TestPace_waits(t *testing.T) {
	testCases := []struct {
		name  string
		count int
		delay time.Duration
	}{
		{"none", 0, 0},
		{"one", 1, time.Millisecond},
		{"two", 2, time.Millisecond},
		{"many", 15, 5 * time.Millisecond},
		{"more than many", 64, time.Millisecond},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var (
				currentCount int
				mu           sync.Mutex
			)

			s := httptest.NewServer(handy.Pace(
				tc.count,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					defer w.Write([]byte("a response"))
					mu.Lock()
					currentCount++
					if currentCount > tc.count {
						t.Errorf("executing concurrent request %d, expected no more than %d", currentCount, tc.count)
					}
					mu.Unlock()
					<-time.After(tc.delay * time.Duration(rand.Intn(10)))
					mu.Lock()
					currentCount--
					mu.Unlock()
				}),
			))

			tick := make(chan int)
			var wg sync.WaitGroup
			go func() {
				for i := range tick {
					go func(i int) {
						defer wg.Done()
						res, err := http.Get(fmt.Sprintf("%s/?i=%d", s.URL, i))
						if err != nil {
							t.Fatalf("unexpected error getting result: %+v", err)
						}
						defer res.Body.Close()
					}(i)
				}
			}()

			for i := 0; i < 3*tc.count; i++ {
				wg.Add(1)
				tick <- i
			}
			wg.Wait()
			close(tick)
		})
	}
}

func TestPace_concurrent(t *testing.T) {

	gen := make(chan time.Duration)
	go func() {
		gen <- time.Millisecond * 10
		gen <- time.Millisecond
		close(gen)
	}()

	s := httptest.NewServer(handy.Pace(
		3, // one more than the number of expected requests
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer w.Write([]byte("a response"))
			delay := <-gen
			<-time.After(delay)
		}),
	))

	doRequest := func(i int, response chan []byte) {
		defer close(response)
		res, err := http.Get(s.URL)
		if err != nil {
			t.Fatalf("unexpected error getting result: %+v", err)
		}
		defer res.Body.Close()
	}
	firstResponse := make(chan []byte)
	go doRequest(1, firstResponse)
	<-time.After(time.Millisecond)
	secondResponse := make(chan []byte)
	go doRequest(2, secondResponse)

	var complete bool
	for {
		// if request n can return before request n-1, it is concurrent
		select {
		case <-firstResponse:
			if !complete {
				t.Errorf("unexpected first request returned first")
			}
			return
		case <-secondResponse:
			complete = true
		}
	}
}
