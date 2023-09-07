package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	errs "github.com/muhammadswa/personal-library/internal/errors"
)

func (bc *Controllers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		errs.ClientError(w, "Error parsing id")
		return
	}

	err = bc.repos.DeleteBook(r.Context(), id)

	if err != nil {
		errs.ServerError(w, "Error deleting book")
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
