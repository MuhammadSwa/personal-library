package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/muhammadswa/personal-library/internal/database"
	errs "github.com/muhammadswa/personal-library/internal/errors"
	"github.com/muhammadswa/personal-library/internal/validator"
	"github.com/muhammadswa/personal-library/pkg/models"
	"github.com/muhammadswa/personal-library/pkg/templates"
)

func (bc *Controllers) CreateBook(w http.ResponseWriter, r *http.Request) {
	templateData := templates.NewTemplateData(bc.session, r)
	templateData.Form = &models.BookForm{}
	templates.Render(w, "create_book", templateData)
}

func (bc *Controllers) CreateBookPost(w http.ResponseWriter, r *http.Request) {
	// TODO: make a helper function? form(r)
	// parse form
	err := r.ParseForm()
	if err != nil {
		errs.WebClientErr(w, "Error parsing form")
		return
	}
	form := &models.BookForm{}

	err = models.DecodePostForm(r, &form)
	if err != nil {
		errs.WebClientErr(w, "Error decoding form")
		return
	}

	// validate form
	form.CheckField((validator.NotBlank(form.Title)), "title", "This field can't be blank")

	book := database.Book{
		Isbn:             form.Isbn,
		Title:            form.Title,
		Author:           form.Author,
		Category:         form.Category,
		Publisher:        form.Publisher,
		YearOfPublishing: form.YearOfPublishing,
		Img:              form.Img,
		NumberOfPages:    form.NumberOfPages,
		PersonalRating:   form.PersonalRating,
		PersonalNotes:    form.PersonalNotes,
		ReadStatus:       form.ReadStatus,
		ReadDate:         form.ReadDate,
	}

	if !form.Valid() {
		data := templates.NewTemplateData(bc.session, r)
		data.Book = &book
		data.Form = form
		templates.Render(w, "create_book", data)
		// c.templateCache.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// check if isbn is valid
	// TODO: look up for the isbn of the title?

	// for i := 0; i < 100; i++ {
	// 	_, _ = bc.repos.CreateBook(r.Context(), database.CreateBookParams{
	// 		UserID:           userId,
	// 		Isbn:             form.Isbn,
	// 		Title:            form.Title,
	// 		Author:           form.Author,
	// 		Category:         form.Category,
	// 		Publisher:        form.Publisher,
	// 		YearOfPublishing: form.YearOfPublishing,
	// 		Img:              form.Img,
	// 		NumberOfPages:    form.NumberOfPages,
	// 		PersonalRating:   form.PersonalRating,
	// 		PersonalNotes:    form.PersonalNotes,
	// 		ReadStatus:       form.ReadStatus,
	// 		ReadDate:         form.ReadDate,
	// 	})
	// }
	userId := bc.session.GetInt32(r.Context(), "userId")
	id, err := bc.repos.CreateBook(r.Context(), database.CreateBookParams{
		UserID:           userId,
		Isbn:             form.Isbn,
		Title:            form.Title,
		Author:           form.Author,
		Category:         form.Category,
		Publisher:        form.Publisher,
		YearOfPublishing: form.YearOfPublishing,
		Img:              form.Img,
		NumberOfPages:    form.NumberOfPages,
		PersonalRating:   form.PersonalRating,
		PersonalNotes:    form.PersonalNotes,
		ReadStatus:       form.ReadStatus,
		ReadDate:         form.ReadDate,
	})
	if err != nil {
		fmt.Println(err)
		errs.WebServerErr(w, "Error creating book")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)
}

func (bc *Controllers) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	pageStr := httprouter.ParamsFromContext(r.Context()).ByName("page")
	query := r.URL.Query().Get("q")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs.WebClientErr(w, "Error parsing offset")
		return
	}
	// TODO: check maximum offset can be reached, sql rows?

	booksLen, err := bc.repos.GetBooksLength(r.Context())
	if err != nil {
		errs.WebServerErr(w, "err getting books")
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
		errs.WebServerErr(w, "err getting books")
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
		errs.WebClientErr(w, "Error parsing id")
		return
	}
	book, err := bc.repos.GetBookByID(r.Context(), id)
	if err != nil {
		errs.WebServerErr(w, "Error getting book")
		return
	}
	data := templates.NewTemplateData(bc.session, r)
	data.Book = book
	templates.Render(w, "book_details", data)
}

func (bc *Controllers) EditBook(w http.ResponseWriter, r *http.Request) {
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errs.WebClientErr(w, "Error parsing id")
		return
	}
	book, err := bc.repos.GetBookByID(r.Context(), id)
	if err != nil {
		errs.WebServerErr(w, "Error getting book")
		return
	}
	data := templates.NewTemplateData(bc.session, r)
	data.Book = book
	data.Form = &models.BookForm{}
	templates.Render(w, "edit_book", data)
}

func (bc *Controllers) EditBookPut(w http.ResponseWriter, r *http.Request) {
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errs.WebClientErr(w, "Error parsing id")
		return
	}
	err = r.ParseForm()
	if err != nil {
		errs.WebClientErr(w, "Error parsing form")
		return
	}
	form := &models.BookForm{}

	err = models.DecodePostForm(r, &form)
	if err != nil {
		errs.WebClientErr(w, "Error decoding form")
		return
	}

	// validate form
	form.CheckField((validator.NotBlank(form.Title)), "title", "This field can't be blank")

	book := database.Book{
		ID:               int32(id),
		Isbn:             form.Isbn,
		Title:            form.Title,
		Author:           form.Author,
		Category:         form.Category,
		Publisher:        form.Publisher,
		YearOfPublishing: form.YearOfPublishing,
		Img:              form.Img,
		NumberOfPages:    form.NumberOfPages,
		PersonalRating:   form.PersonalRating,
		PersonalNotes:    form.PersonalNotes,
		ReadStatus:       form.ReadStatus,
		ReadDate:         form.ReadDate,
	}

	if !form.Valid() {
		data := templates.NewTemplateData(bc.session, r)
		data.Book = &book
		data.Form = form
		templates.Render(w, "edit_book", data)
		// c.templateCache.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	// update the book
	err = bc.repos.UpdateBook(r.Context(), database.UpdateBookParams{
		ID:               int32(id),
		Isbn:             form.Isbn,
		Title:            form.Title,
		Author:           form.Author,
		Category:         form.Category,
		Publisher:        form.Publisher,
		YearOfPublishing: form.YearOfPublishing,
		Img:              form.Img,
		NumberOfPages:    form.NumberOfPages,
		PersonalRating:   form.PersonalRating,
		PersonalNotes:    form.PersonalNotes,
		ReadStatus:       form.ReadStatus,
		ReadDate:         form.ReadDate,
	})
	if err != nil {
		errs.WebServerErr(w, "Error updating book")
		return
	}
	// rediret to the book details
	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)

}

func (bc *Controllers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		errs.WebClientErr(w, "Error parsing id")
		return
	}

	err = bc.repos.DeleteBook(r.Context(), id)

	if err != nil {
		errs.WebServerErr(w, "Error deleting book")
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
