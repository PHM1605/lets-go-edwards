package main

import (
	"bytes"
	"io"
	"lets-go-edwards/internal/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	// NOTE: Get a ResponseRecorder (similar to ResponseWriter; but it records instead of responding)
	rr := httptest.NewRecorder()
	// Dummy request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Run the "ping-handler"
	ping(rr, r)
	// Get Recorder's result i.e. a http.Response
	rs := rr.Result()

	// Check 200 status
	assert.Equal(t, rs.StatusCode, http.StatusOK)
	// Check body of response equal "OK"
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

func TestPingServer(t *testing.T) {
	app := newTestApplication(t)
	// create test HTTPS server
	ts := newTestServer(t, app.routes()) // handler parameter is our REAL FULL Routes
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

// Test sending GET request to an URL
func TestSnippetView(t *testing.T) {
	// Mock application
	app := newTestApplication(t)
	// Mock server
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string // substring of Response's body
	}{
		{
			name:     "Valid ID", // will be canonically modified to "Valid_ID"
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Send a GET request
			code, _, body := ts.get(t, tt.urlPath)
			// check if status is OK
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" { // "tt" might NOT have "wantBody"; in that case "tt.wantBody" == ""
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

// 303 See Other: good
// 400 Bad Request: form without CSRF token
// 422 Unprocessable Entity: display Signup form again; when name, email, password empty, duplicated etc.
func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Send dummy GET request to get signup form
	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)
	// t.Logf("CSRFToken is: %q", validCSRFToken)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com", // check mocks/users.go to understand
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy form
			form := url.Values{} // url.Values is a map from "string"->"[]string"
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			// Send dummy POST request with that dummy form
			code, _, body := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, code, tt.wantCode)
			// wantFormTag: what we expects Server to send us when /user/signup page was rendered AGAIN (in case of wrong password format etc.)
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

// Check what happens if Client accesses `GET /snippet/create`
// Authenticated: show the Create Snippet Form
// Unauthenticated: redirect to the Login Form
func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusSeeOther)              // status is "Redirect"
		assert.Equal(t, headers.Get("Location"), "/user/login") // redirect to "login page"
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Make a GET to get CSRFToken from Login form
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)
		// Make a POST to login
		form := url.Values{} // map from "string" to "[]string"
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)

		// Now we can GET the "/snippet/create" form
		code, _, body := ts.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})
}
