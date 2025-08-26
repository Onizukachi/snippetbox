package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Onizukachi/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
)

func newTestApplication(t *testing.T) *application {
	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour

	templateCache, err := newTemplateCache("./../../ui/html")
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		infoLog:       log.New(io.Discard, "", 0),
		errorLog:      log.New(io.Discard, "", 0),
		session:       session,
		snippets:      &mock.SnippetModel{},
		users:         &mock.UserModel{},
		templateCache: templateCache,
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(handler http.Handler) *testServer {
	ts := httptest.NewTLSServer(handler)

	jar, _ := cookiejar.New(nil)
	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) Get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
