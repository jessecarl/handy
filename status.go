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
	var blackListCodes = []int{
		// 1xx
		http.StatusContinue,
		http.StatusSwitchingProtocols,
		http.StatusProcessing,

		// 2xx
		http.StatusMultiStatus,
		http.StatusAlreadyReported,
		http.StatusIMUsed,

		// 3xx
		http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusUseProxy,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCodeString := strings.Trim(r.URL.Path, "/")
		statusCode, err := strconv.Atoi(statusCodeString)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		statusText := http.StatusText(statusCode)
		if statusText == "" {
			http.NotFound(w, r)
			return
		}
		for _, code := range blackListCodes {
			if code == statusCode {
				http.NotFound(w, r)
				return
			}
		}
		w.Header().Set("x-status-code", statusCodeString)
		w.Header().Set("x-status", statusText)
		w.WriteHeader(statusCode)
		w.Write([]byte(statusCodeString + " " + statusText + "\n"))
	})
}
