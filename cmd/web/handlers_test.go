package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

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
	ts := newTestServer(t, app.routes())

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

func TestSignUpUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	defer ts.Close()
	_, _, body := ts.Get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Logf("CSRF token from form: %s", csrfToken)

	// Посмотрим, что есть в cookie jar
	for _, c := range ts.Client().Jar.Cookies(&url.URL{Scheme: "https", Host: ts.URL[8:]}) {
		t.Logf("Cookie: %s = %s", c.Name, c.Value)
	}

	tests := []struct {
		name      string
		nameField string
		email     string
		password  string
		csrfToken string
		wantCode  int
		wantBody  []byte
	}{
		{
			name:      "Valid submission",
			nameField: "Bob",
			email:     "bob@example.com",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusSeeOther,
			wantBody:  nil,
		},
		{
			name:      "Empty name",
			nameField: "",
			email:     "bob@example.com",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field cannot be blank"),
		},
		{
			name:      "Empty email",
			nameField: "Bob",
			email:     "",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field cannot be blank"),
		},
		{
			name:      "Empty password",
			nameField: "Bob",
			email:     "bob@example.com",
			password:  "",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field cannot be blank"),
		},
		{
			name:      "Invalid email (incomplete domain)",
			nameField: "Bob",
			email:     "bob@example.",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field is invalid"),
		},
		{
			name:      "Invalid email (missing @)",
			nameField: "Bob",
			email:     "bobexample.com",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field is invalid"),
		},
		{
			name:      "Invalid email (missing local part)",
			nameField: "Bob",
			email:     "@example.com",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field is invalid"),
		},
		{
			name:      "Short password",
			nameField: "Bob",
			email:     "bob@example.com",
			password:  "pa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("This field is too short (minimum is 10 characters)"),
		},
		{
			name:      "Duplicate email",
			nameField: "Bob",
			email:     "dupe@example.com",
			password:  "validPa$$word",
			csrfToken: csrfToken,
			wantCode:  http.StatusOK,
			wantBody:  []byte("Address is already in use"),
		},
		{
			name:      "Invalid CSRF Token",
			nameField: "Bob",
			email:     "bob@example.com",
			password:  "validPa$$word",
			csrfToken: "wrongToken",
			wantCode:  http.StatusBadRequest,
			wantBody:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.nameField)
			form.Add("email", test.email)
			form.Add("password", test.password)
			form.Add("csrf_token", test.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			if code != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, code)
			}

			if !bytes.Contains(body, test.wantBody) {
				t.Errorf("want body to contain %q", test.wantBody)
			}
		})
	}
}
