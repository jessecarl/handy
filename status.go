// © Copyright 2016 Jesse Allen. All rights reserved.
// Released under the MIT license found in the LICENSE file.

package handy

import (
	"net/http"
	"strconv"
	"strings"
)

// ServeStatus provides a handler that will respond with the http status
// indicated by the path. Only 2xx, 4xx, and 5xx status codes are supported
// at the moment.
func ServeStatus() http.Handler {
	var whiteListCodes = []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNonAuthoritativeInfo,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent,
		http.StatusBadRequest,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCodeString := strings.Trim(r.URL.Path, "/")
		statusCode, err := strconv.Atoi(statusCodeString)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		var ok bool
		for _, code := range whiteListCodes {
			if code == statusCode {
				ok = true
				break
			}
		}
		if !ok {
			http.NotFound(w, r)
			return
		}
		statusText := http.StatusText(statusCode)
		w.Header().Set("x-status-code", statusCodeString)
		w.Header().Set("x-status", statusText)
		w.WriteHeader(statusCode)
		w.Write([]byte(statusCodeString + " " + statusText + "\n"))
	})
}
