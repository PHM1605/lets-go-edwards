package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// routes() method: returns "servemux" that includes all application's routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// File server for files in "./ui/static"
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// Register the Fileserver as a handler
	// NOTE: must strip off "/static" of REQUEST
	// e.g. "/static/css/style.css" => "css/style.css"
	// otherwise fullpath is "/static/static/css/style.css" which is wrong
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// New Middleware for Session handling
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// Public routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected routes => wrap "dynamic" in an extra middleware
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// NOTE: we wrap Middlewares around "mux" here with "alice"
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
