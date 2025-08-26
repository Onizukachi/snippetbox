package main

import (
	"bytes"
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

func TestShopwSnippet(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(app.routes())

	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("Test content")},
		{"Non-existing ID", "/snippet/2", http.StatusNotFound, nil},
		{"Not valid ID", "/snippet/-10", http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.Get(t, test.urlPath)

			if code != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, code)
			}

			if !bytes.Contains(body, test.wantBody) {
				t.Errorf("want body to contain %q", test.wantBody)
			}
		})
	}
}
