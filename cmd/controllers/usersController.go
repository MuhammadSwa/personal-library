package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	errs "github.com/muhammadswa/personal-library/cmd/errors"
	"github.com/muhammadswa/personal-library/cmd/models"
	"github.com/muhammadswa/personal-library/cmd/repositories"
	"github.com/muhammadswa/personal-library/cmd/templates"
	"github.com/muhammadswa/personal-library/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type UsersController struct {
	usersRepository *repositories.UsersRepository
	session         *scs.SessionManager
}

func NewUsersController(usersRepository *repositories.UsersRepository, session *scs.SessionManager) *UsersController {
	return &UsersController{
		usersRepository: usersRepository,
		session:         session,
	}
}

// TODO hahdlers like home, about, contact, put them in a separate controller (staticController)
func (uc *UsersController) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	templateData := templates.NewTemplateData(uc.session, r)
	templateData.Form = &models.LoginForm{}
	// uc.session.Put(r.Context(), "flash", "")
	templates.Render(w, "login", templateData)
}

func (uc *UsersController) LoginPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO: parse login form
	err := r.ParseForm()
	if err != nil {
		errs.WebClientErr(w, "Error parsing form")
		return
	}
	form := &models.LoginForm{}

	err = models.DecodePostForm(r, &form)
	if err != nil {
		errs.WebClientErr(w, "Error decoding form")
		return
	}

	// TODO: validate login form
	// validate form
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.ValidateEmail(form.Email), "email", "Email isn't valid")

	data := templates.NewTemplateData(uc.session, r)
	if !form.Valid() {
		data.Form = form
		templates.Render(w, "login", data)
		return
	}

	// authenticate user
	// TODO: Where to put validation
	// TODO: Use repo layer instead of database directly
	// Get user from db by email
	user, err := uc.usersRepository.GetUserByEmail(r.Context(), form.Email)
	if err != nil {
		form.AddNonFieldError("Invalid login credentials")
		data.Form = form
		templates.Render(w, "login", data)
		return
	}

	// check for password match
	// template for error
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(form.Password))
	if err != nil {
		form.AddNonFieldError("Invalid login credentials")
		data.Form = form
		templates.Render(w, "login", data)
		return
	}

	err = uc.session.RenewToken(r.Context())
	if err != nil {
		errs.WebServerErr(w, "Error renewing session token")
		return
	}

	uc.session.Put(r.Context(), "authenticatedUserID", user.ID)
	uc.session.Put(r.Context(), "userId", user.ID)
	// uc.session.Put(r.Context(), "flash", "Login successful")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc *UsersController) SignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	data := templates.NewTemplateData(uc.session, r)
	data.Form = &models.SignupForm{}
	templates.Render(w, "signup", data)
}

func (uc *UsersController) SignupPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		errs.WebClientErr(w, "Error parsing form")
		return
	}
	form := &models.SignupForm{}

	err = models.DecodePostForm(r, &form)
	if err != nil {
		errs.WebClientErr(w, "Error decoding form")
		return
	}

	// TODO: validate form
	// validate form
	const PASSWORD_MIN_LENGTH = 8
	const USERNAME_MIN_LENGTH = 3

	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, PASSWORD_MIN_LENGTH), "password",
		"Password must be at least 8 characters long")

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidateEmail(form.Email), "email", "Email isn't valid")

	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Username, USERNAME_MIN_LENGTH), "username",
		"Username must be at least 3 characters long")

	if !form.Valid() {
		data := templates.NewTemplateData(uc.session, r)
		data.Form = form
		templates.Render(w, "signup", data)
		return
	}

	// create a new user
	id, err := uc.usersRepository.CreateUser(r.Context(), form.Email, form.Password, form.Username)
	fmt.Println("Id from signup", id)
	if err != nil {
		templateData := templates.NewTemplateData(uc.session, r)
		templateData.Form = form

		pqerr := err.(*pq.Error)

		if pqerr.Code == "23505" && strings.Contains(pqerr.Message, "users_username_key") {
			form.AddFieldError("username", "Username already exists")
		}
		if pqerr.Code == "23505" && strings.Contains(pqerr.Message, "users_email_key") {
			form.AddFieldError("email", "Email already exists")
		}

		templates.Render(w, "signup", templateData)
		return
	}

	err = uc.session.RenewToken(r.Context())
	if err != nil {
		errs.WebServerErr(w, "Error renewing session token")
		return
	}

	uc.session.Put(r.Context(), "authenticatedUserID", id)
	uc.session.Put(r.Context(), "flash", "Login successful")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc *UsersController) LogoutPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := uc.session.RenewToken(r.Context())
	if err != nil {
		errs.WebServerErr(w, "Error renewing session token")
		return
	}
	uc.session.Remove(r.Context(), "authenticatedUserID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc *UsersController) ForgotPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (uc *UsersController) ForgotPasswordPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}

func (uc *UsersController) authenticateUser(r *http.Request, form models.LoginForm) error {
	user, err := uc.usersRepository.GetUserByEmail(r.Context(), form.Email)
	if err != nil {
		return err
	}

	// check for password match
	// template for error
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(form.Password))
	if err != nil {
		return err
	}

	err = uc.session.RenewToken(r.Context())
	if err != nil {
		return err
	}
	return nil
}
