package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	errs "github.com/muhammadswa/personal-library/internal/errors"
	"github.com/muhammadswa/personal-library/pkg/templates"
)

func (sc *Controllers) Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !sc.session.Exists(r.Context(), "authenticatedUserID") {
		data := templates.NewTemplateData(sc.session, r)
		templates.Render(w, "home_not_logged", data)
		return
	}

	offset := 0
	books, err := sc.repos.GetBooks(r.Context(), 0, "", offset)
	if err != nil {
		errs.WebServerErr(w, "err getting books")
		return
	}

	data := templates.NewTemplateData(sc.session, r)
	data.Books = &books
	templates.Render(w, "home", data)
}

func (sc *Controllers) Profile(w http.ResponseWriter, r *http.Request) {
	if !sc.session.Exists(r.Context(), "authenticatedUserID") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Header().Add("Cache-Control", "no-store")
		return
	}
	data := templates.NewTemplateData(sc.session, r)
	templates.Render(w, "profile", data)
}
