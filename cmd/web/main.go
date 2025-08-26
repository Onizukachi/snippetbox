package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Onizukachi/snippetbox/pkg/models"
	"github.com/Onizukachi/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

var contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets interface {
		Insert(title, content, expires string) (int, error)
		Get(id int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	users interface {
		Insert(name, email, password string) error
		Authenticate(email, password string) (int, error)
		Get(id int) (*models.User, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "hikaru:potolok149@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secter", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:       infoLog,
		errorLog:      errLog,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
