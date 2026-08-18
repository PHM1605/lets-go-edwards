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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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

// Authentication
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// Empty user signup form data
	var form userSignupForm
	// Parse form data
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate form contents - 2nd and 3rd parameters are for Validator Error dictionary
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// If there is any error(s), re-render Signup page with red color here and there AND filled email/password etc.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	// Form format data is good; now we check for DB insertion
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		// read "Insert": it returns "ErrDuplicateEmail" when email exists
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		}
	}
	// Put "flash" message to session of User
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	// After Sign up successfully => return to Login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

// NOTE: If login successfully => add "authenticatedUserID" to Session to know that that User has logged in
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check if user fills in valid Email/Password FORMAT
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Now login info format is correct; we check what User enters
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		// if User enters gibberish => render Login Form with red color here and there
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Renew Session (getting new Session ID ONLY, KEEP data) after User successfully logged in
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// NOTE: add "authenticatedUserID" field to Session; so that next time when he calls, Server knows that he has logged in
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect logged in User to Snippet Create Page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Refresh token ID
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// Remove "authenticatedUserID" from Session => next time the request is of anonymous User
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
