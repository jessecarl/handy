// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy

import (
	"net/http"
)

// Pace creates a static sized pool of workers to handle requests
// with the work handler.
func Pace(count int, work http.Handler) http.Handler {
	sem := make(chan struct{}, count)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sem <- struct{}
		defer func(){
			<-sem
		}()
		work.ServeHTTP(w, r)
	})
}
