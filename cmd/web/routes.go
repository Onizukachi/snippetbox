package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standartMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc((http.HandlerFunc(app.createSnippetForm))))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.showSnippet)))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(http.HandlerFunc(app.logoutUser)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standartMiddleware.Then(mux)
}
