package controllers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	errs "github.com/muhammadswa/personal-library/cmd/errors"
	"github.com/muhammadswa/personal-library/cmd/repositories"
	"github.com/muhammadswa/personal-library/cmd/templates"
)

type StaticController struct {
	session         *scs.SessionManager
	booksRepository *repositories.BooksRepository
}

func NewStaticController(booksRepository *repositories.BooksRepository, session *scs.SessionManager) *StaticController {
	return &StaticController{
		session:         session,
		booksRepository: booksRepository,
	}
}

func (sc *StaticController) Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !sc.session.Exists(r.Context(), "authenticatedUserID") {
		data := templates.NewTemplateData(sc.session, r)
		templates.Render(w, "home_not_logged", data)
		return
	}

	offset := 0
	books, err := sc.booksRepository.GetBooks(r.Context(), offset)
	if err != nil {
		errs.WebServerErr(w, "err getting books")
		return
	}

	data := templates.NewTemplateData(sc.session, r)
	data.Books = &books
	templates.Render(w, "home", data)
}

func (sc *StaticController) Profile(w http.ResponseWriter, r *http.Request) {
	if !sc.session.Exists(r.Context(), "authenticatedUserID") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Header().Add("Cache-Control", "no-store")
		return
	}
	data := templates.NewTemplateData(sc.session, r)
	templates.Render(w, "profile", data)
}
