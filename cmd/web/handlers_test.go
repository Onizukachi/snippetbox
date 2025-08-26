package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(app.routes())

	defer ts.Close()

	code, _, body := ts.Get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want to body equal %q", "OK")
	}
}
