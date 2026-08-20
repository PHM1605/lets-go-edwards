package main

import (
	"lets-go-edwards/ui"
	"net/http"

	"github.com/justinas/alice"
)

// routes() method: returns "servemux" that includes all application's routes
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// File server for files in "./ui/static"
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Health check route
	mux.HandleFunc("GET /ping", ping)

	// New Middleware for Session handling
	// app.authenticate: middleware to 1/ check Session for UserID 2/ if UserID in DB => add {<isAuthenticatedContextKey>: true} to Context
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Public routes
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected routes => wrap "dynamic" in an extra middleware
	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	// View Profile: only available in "protected" route
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	// GET password update form
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	// POST password update request
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// NOTE: we wrap Middlewares around "mux" here with "alice"
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
