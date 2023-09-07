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

func (bc *Controllers) EditBook(w http.ResponseWriter, r *http.Request) {
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
	data.Form = &models.BookForm{
		ID:               book.ID,
		Isbn:             book.Isbn,
		Title:            book.Title,
		Author:           book.Author,
		Category:         book.Category,
		Publisher:        book.Publisher,
		YearOfPublishing: book.YearOfPublishing,
		Img:              book.Img,
		NumberOfPages:    book.NumberOfPages,
		PersonalRating:   book.PersonalRating,
		PersonalNotes:    book.PersonalNotes,
		ReadStatus:       book.ReadStatus,
		ReadDate:         book.ReadDate,
	}
	templates.Render(w, "edit_book", data)
}

func (bc *Controllers) EditBookPut(w http.ResponseWriter, r *http.Request) {
	idStr := httprouter.ParamsFromContext(r.Context()).ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errs.ClientError(w, "Error parsing id")
		return
	}
	err = r.ParseForm()
	if err != nil {
		errs.ClientError(w, "Error parsing form")
		return
	}
	form := &models.BookForm{}

	err = models.DecodePostForm(r, &form)
	if err != nil {
		errs.ClientError(w, "Error decoding form")
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
		errs.ServerError(w, "Error updating book")
		return
	}
	// rediret to the book details
	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)

}
