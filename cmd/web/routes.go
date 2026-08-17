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
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// NOTE: we wrap Middlewares around "mux" here with "alice"
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
