package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// set Headers in a Middleware - this will go BEFORE "servemux" so "servemux" is "next"
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set Headers
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

// Middleware to print Request's information on Screen when it arrives first
// make this a Method of "app" to use "app.logger()"
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

// Middleware to return 500 when a request has error
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Middleware to protect a route from unauthenticated access
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if User not authenticated => redirect to Login page
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Authenticated => add "Cache-Control: no-store" header
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

// Middleware to insert CSRF token to HTML
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
