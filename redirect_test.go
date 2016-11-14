// Copyright 2013 Jesse Allen. All rights reserved
// Released under the MIT license found in the LICENSE file.

package handy_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/jessecarl/handy"
)

func TestServePermanentRedirects(t *testing.T) {
	s := httptest.NewServer(handy.ServePermanentRedirects(map[string]string{
		"/google": "http://www.google.com",
		"/blank":  "",
	}))
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirected")
		},
	}
	// test unlisted redirect
	if r, err := c.Get(s.URL + "/nil"); err != nil {
		t.Error(err)
	} else {
		t.Log("Expected: StatusNotFound\tActual: ", r.StatusCode)
		if r.StatusCode != http.StatusNotFound {
			t.Fail()
		}
	}
	// test blank redirect
	if _, err := c.Get(s.URL + "/blank"); err != nil {
		e := err.(*url.Error)
		t.Log("Expected: /\tActual: ", e.URL)
		if e.URL != "/" { // assumes the "" to "/" automatically
			t.Fail()
		}
	} else {
		t.Error("Expected Redirection")
	}
	// test ok redirect
	if _, err := c.Get(s.URL + "/google"); err != nil {
		e := err.(*url.Error)
		t.Log("Expected: http://www.google.com\tActual: ", e.URL)
		if e.URL != "http://www.google.com" {
			t.Fail()
		}
	} else {
		t.Error("Expected Redirection")
	}
}
