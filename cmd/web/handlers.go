package main

import (
	"errors"
	"fmt"
	"lets-go-edwards/internal/models"
	"net/http"
	"strconv"
)

// "home" now can access Global properties of "application" struct
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Add new Header field
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create data to render with Snippets information & current year
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.html", data)
}

// View 1 Snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	// GET Snippet with that ID
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Create template data to render, with Snippet & current Year
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

// Display the form to fill in
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}

// Eventually create snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	// insert to our DB model
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// return Client to relevant page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
