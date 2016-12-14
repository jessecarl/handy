// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy

import (
	"net/http"
)

// Pace creates a static sized pool of workers to handle requests
// with the work handler.
func Pace(count int, work http.Handler) http.Handler {
	type args struct {
		w    http.ResponseWriter
		r    *http.Request
		done chan struct{}
	}
	workCh := make(chan args)
	for i := 0; i < count; i++ {
		go func() {
			for a := range workCh {
				work.ServeHTTP(a.w, a.r)
				close(a.done)
			}
		}()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		workCh <- args{w, r, done}
		<-done
	})
}
