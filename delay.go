// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy

import (
	"net/http"
	"path"
	"strings"
	"time"
)

// ServeWithDelay parses the last element in the path as a `time.Duration`
// and passes the request on to the next Handler after the requested delay.
// That next Handler gets a request that is identical to what a request without
// the delay.
func ServeWithDelay(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r)
		pathElements := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		if len(pathElements) == 0 {
			return
		}
		d, err := time.ParseDuration(pathElements[len(pathElements)-1])
		if err != nil {
			return
		}
		r.URL.Path = "/" + path.Join(pathElements[:len(pathElements)-1]...)
		<-time.After(d)
	})
}
