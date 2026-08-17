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

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// NOTE: we wrap Middlewares around "mux" here with "alice"
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
