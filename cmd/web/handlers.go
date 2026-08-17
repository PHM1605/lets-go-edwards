package main

import (
	"errors"
	"fmt"
	"lets-go-edwards/internal/models"
	"lets-go-edwards/internal/validator"
	"net/http"
	"strconv"
)

// grouping what will be returned to Client
// NOTE: "struct embedding" in last entry
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// "home" now can access Global properties of "application" struct
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// panic("oops! something went wrong")
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
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

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
	data := app.newTemplateData(r) // "Flash" already inside
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

// Display the form to fill in
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, r, http.StatusOK, "create.html", data)
}

// Eventually create snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Filling data to a Form; the formDecoder will use field tags `form:xxx` to fill "form"
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// // NOTE: for the case of Checkbox (many values - we don't have here)
	// // e.g. <input type="checkbox" name="items" value="foo"> Foo
	// // <input type="checkbox" name="items" value="bar"> Bar
	// for i, item := range r.PostForm["items"] { // list of strings, each is "foo"/"bar"
	// 	fmt.Fprintf(w, "%d: Item %s\n", i, item)
	// }

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValues(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// Check if any Error exists
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// insert to our DB model
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Put {"flash": "xxxx"} to Session data
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// return Client to relevant page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
