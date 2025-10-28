package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Vadim-Makhnev/snippetbox/internal/models"
	"github.com/Vadim-Makhnev/snippetbox/internal/validator"
)

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

// Home godoc
// @Summary      Get home page with latest snippets
// @Description  Retrieve the latest snippets and render the home page
// @Tags         pages
// @Produce      html
// @Success      200 {string} string "HTML page"
// @Router       / [get]
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// snippetView godoc
// @Summary      Get snippet by id
// @Description  Retrieve snippet by snippet id
// @Tags         snippets
// @Produce      html
// @Param        id path int true "Snippet ID"
// @Success      200 {string} string "HTML page"
// @Failure      404 {string} string "Snippet not found"
// @Failure      500 {string} string "Internal server error"
// @Router       /snippet/view/{id} [get]
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

// snippetCreate godoc
// @Summary      Show snippet creation form
// @Description  Display the form for creating a new code snippet
// @Tags         snippets
// @Produce      html
// @Success      200 {string} string "Snippet creation form"
// @Router       /snippet/create [get]
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// snippetCreatePost godoc
// @Summary      Create new snippet
// @Description  Create a new code snippet with validation
// @Tags         snippets
// @Accept       x-www-form-urlencoded
// @Produce      html
// @Param        title formData string true "Snippet title" minlength(1) maxlength(100)
// @Param        content formData string true "Snippet content" minlength(1)
// @Param        expires formData int true "Expiration in days" Enums(1, 7, 365)
// @Success      303 {string} string "Redirect to created snippet"
// @Failure      400 {string} string "Bad request - invalid form data"
// @Failure      422 {string} string "Unprocessable entity - validation failed"
// @Failure      500 {string} string "Internal server error"
// @Router       /snippet/create [post]
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// userSignup godoc
// @Summary      Show user registration form
// @Description  Display the form for new user registration
// @Tags         auth
// @Produce      html
// @Success      200 {string} string "User registration form"
// @Router       /user/signup [get]
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

// userSignupPost godoc
// @Summary      Register new user
// @Description  Create a new user account with email and password validation. Checks for duplicate emails.
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      html
// @Param        name formData string true "User's full name" minlength(1) maxlength(255)
// @Param        email formData string true "User's email address" format(email)
// @Param        password formData string true "User's password" minlength(8)
// @Success      303 {string} string "Redirect to login page with success message"
// @Failure      400 {string} string "Bad request - invalid form data"
// @Failure      422 {string} string "Unprocessable entity - validation failed or duplicate email"
// @Failure      500 {string} string "Internal server error"
// @Router       /user/signup [post]
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// userLogin godoc
// @Summary      Show login form
// @Description  Display the form for user authentication
// @Tags         auth
// @Produce      html
// @Success      200 {string} string "User login form"
// @Router       /user/login [get]
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

// userLoginPost godoc
// @Summary      Authenticate user
// @Description  Verify user credentials and create session. On success, redirects to snippet creation page.
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      html
// @Param        email formData string true "User's email address" format(email)
// @Param        password formData string true "User's password"
// @Success      303 {string} string "Redirect to /snippet/create with active session"
// @Failure      400 {string} string "Bad request - invalid form data"
// @Failure      401 {string} string "Unauthorized - invalid credentials"
// @Failure      422 {string} string "Unprocessable entity - validation failed"
// @Failure      500 {string} string "Internal server error"
// @Router       /user/login [post]
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// userLogoutPost godoc
// @Summary      Logout user
// @Description  Terminate user session and redirect to home page with confirmation message
// @Tags         auth
// @Produce      html
// @Success      303 {string} string "Redirect to home page with logout confirmation"
// @Failure      500 {string} string "Internal server error"
// @Router       /user/logout [post]
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
