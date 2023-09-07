package controllers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	errs "github.com/muhammadswa/personal-library/internal/errors"
	"github.com/muhammadswa/personal-library/pkg/templates"
)

func (bc *Controllers) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	pageStr := httprouter.ParamsFromContext(r.Context()).ByName("page")
	query := r.URL.Query().Get("q")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs.ClientError(w, "Error parsing offset")
		return
	}
	// TODO: check maximum offset can be reached, sql rows?

	booksLen, err := bc.repos.GetBooksLength(r.Context())
	if err != nil {
		errs.ServerError(w, "err getting books")
		return
	}

	// if offset is 1, then offset = 0 * 10 = 0 => first 10 books
	// if offset is 2, then offset = 1 * 10 = 10 => next 10 books
	offset := (page - 1) * 10

	if offset >= booksLen {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	nextPage := true
	if offset+10 >= booksLen {
		nextPage = false
	}

	userId := bc.session.GetInt32(r.Context(), "userId")
	books, err := bc.repos.GetBooks(r.Context(), userId, query, int(offset))
	if err != nil {
		errs.ServerError(w, "err getting books")
		return
	}
	if len(books) < 10 {
		nextPage = false
	}

	data := templates.NewTemplateData(bc.session, r)
	data.Books = &books
	data.IsPageNext = nextPage
	data.NextPage = page + 1
	data.Query = query

	hxTrigger := r.Header.Get("HX-Trigger")
	if hxTrigger == "search-books" || hxTrigger == "load-more-btn" {
		templates.RenderFragment(w, "books_list", data)
		return
	}
	templates.Render(w, "books", data)
}

func (bc *Controllers) GetBookByID(w http.ResponseWriter, r *http.Request) {
	// isbn := ps.ByName("isbn")
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errs.ClientError(w, "Error parsing id")
		return
	}
	book, err := bc.repos.GetBookByID(r.Context(), id)
	if err != nil {
		errs.ServerError(w, "Error getting book")
		return
	}
	data := templates.NewTemplateData(bc.session, r)
	data.Book = book
	templates.Render(w, "book_details", data)
}
