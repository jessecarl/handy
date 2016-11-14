// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package handy

import (
	"net/http"
	"sync"
)

type permanentHandler struct {
	redirects map[string]string
	sync.RWMutex
}

// ServePermanentRedirects provides an http.Handler that will permanently redirect
// any requests based on the redirects map
//     redirects[requestedURL] = redirectedURL
func ServePermanentRedirects(redirects map[string]string) http.Handler {
	h := new(permanentHandler)
	h.init(redirects)
	return http.Handler(h)
}

func (h *permanentHandler) init(redirects map[string]string) {
	h.Lock()
	defer h.Unlock()
	h.redirects = redirects
}

// Serve Permanent Redirects
func (h *permanentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.RLock()
	defer h.RUnlock()
	url, ok := h.redirects[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}
