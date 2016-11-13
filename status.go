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
		http.StatusUnauthorized,
		http.StatusPaymentRequired,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusMethodNotAllowed,
		http.StatusNotAcceptable,
		http.StatusProxyAuthRequired,
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusGone,
		http.StatusLengthRequired,
		http.StatusPreconditionFailed,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestURITooLong,
		http.StatusUnsupportedMediaType,
		http.StatusRequestedRangeNotSatisfiable,
		http.StatusExpectationFailed,
		http.StatusTeapot,
		http.StatusUnprocessableEntity,
		http.StatusLocked,
		http.StatusFailedDependency,
		http.StatusUpgradeRequired,
		http.StatusPreconditionRequired,
		http.StatusTooManyRequests,
		http.StatusRequestHeaderFieldsTooLarge,
		http.StatusUnavailableForLegalReasons,
		http.StatusInternalServerError,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusHTTPVersionNotSupported,
		http.StatusVariantAlsoNegotiates,
		http.StatusInsufficientStorage,
		http.StatusLoopDetected,
		http.StatusNotExtended,
		http.StatusNetworkAuthenticationRequired,
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
