package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFirstHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/first", nil)
	w := httptest.NewRecorder()
	firstHandler(w, req)
	defer w.Result().Body.Close()

	if w.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", w.Code)
	}

	if !strings.Contains(w.Body.String(), "first") {
		t.Errorf(
			`response body "%s" does not contain "first"`,
			w.Body.String(),
		)
	}
}

func TestSecondHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/second", nil)
	w := httptest.NewRecorder()
	secondHandler(w, req)
	defer w.Result().Body.Close()

	if w.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", w.Code)
	}

	if !strings.Contains(w.Body.String(), "second") {
		t.Errorf(
			`response body "%s" does not contain "second"`,
			w.Body.String(),
		)
	}
}
