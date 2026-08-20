package main

import (
	"bytes"
	"io"
	"lets-go-edwards/internal/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 1/ Test headers in Response
// 2/ Test the next handler, whether it's called by commonHeaders() correctly
func TestCommonHeaders(t *testing.T) {
	// Init a fresh "httptest.ResponseRecorder"
	rr := httptest.NewRecorder()
	// Dummy request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Dummy handler; which will be called after Middleware commonHeaders()
	// this "next-handler" will have status of 200 & body of "OK"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// Execute full middleware flow: request -> middleware -> next-handler -> response-record
	commonHeaders(next).ServeHTTP(rr, r)
	// Read out recorded result
	rs := rr.Result() // http.Response

	// Check Content-Security-Policy header
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)
	// Check Referrer-Policy header
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)
	// Check X-Content-Type-Options header
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)
	// Check X-Frame-Options header
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)
	// Check X-XSS-Protection header
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)
	// Check Server type header
	expectedValue = "Go"
	assert.Equal(t, rs.Header.Get("Server"), expectedValue)

	// Check status code of Response
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// Check body of Response ("OK" or not)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
