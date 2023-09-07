package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/muhammadswa/personal-library/internal/database"
	errs "github.com/muhammadswa/personal-library/internal/errors"
	"github.com/muhammadswa/personal-library/internal/validator"
	"github.com/muhammadswa/personal-library/pkg/models"
	"github.com/muhammadswa/personal-library/pkg/templates"
)

func (bc *Controllers) CreateBook(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "create_book", nil)
}

func (bc *Controllers) CreateBookPost(w http.ResponseWriter, r *http.Request) {
	// TODO: make a helper function? form(r)
	// parse form
	err := r.ParseForm()
	if err != nil {
		errs.ClientError(w, "Error parsing form")
		return
	}
	// TODO: try to remove it ? do we need it?
	form := &models.BookForm{}

	err = models.DecodePostForm(r, &form)
	if err != nil {
		errs.ClientError(w, "Error decoding form")
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
		return
	}

	// check if isbn is valid
	// TODO: look up for the isbn of the title?

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
		errs.ServerError(w, "Error creating book")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)
}

func (bc *Controllers) FetchByIsbn(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	form, err := getBookByIsbn(isbn)
	if err != nil {
		errs.ServerError(w, "Error fetching book")
		return
	}
	data := templates.NewTemplateData(bc.session, r)
	data.Form = form

	templates.RenderFragment(w, "book_form", data)
}

// put this is it's own file and package? in pkg/api?
func getBookByIsbn(isbn string) (*models.BookForm, error) {

	// TODO: make sure isbn is valid (13 digits) and not empty , no - or spaces ?
	// maybe - is okay? ? use regex client side
	openLibraryUrl := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&jscmd=data&format=json", isbn)
	resp, err := http.Get(openLibraryUrl)
	if err != nil {
		return nil, fmt.Errorf("Error fetching book from open library: %v", err)
	}
	defer resp.Body.Close()

	// TODO: use form/playground??

	openLibRes := models.OpenLibResponse{}
	err = json.NewDecoder(resp.Body).Decode(&openLibRes)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("ISBN:%v", isbn)
	jsonBook := openLibRes[key]

	if err != nil {
		return nil, err
	}
	// Dec 07, 2018
	splits := strings.Split(jsonBook.PublishDate, ", ")[1]
	publishDate, err := strconv.Atoi(splits)
	if err != nil {
		return nil, fmt.Errorf("Error converting publish date to int: %v", err)
	}
	form := models.BookForm{
		Isbn:             isbn,
		Title:            jsonBook.Title,
		Publisher:        jsonBook.Publishers[0].Name,
		YearOfPublishing: int32(publishDate),
		Img:              jsonBook.Cover.Large,
		NumberOfPages:    int32(jsonBook.NumberOfPages),
		// TODO: contactenate authors
		Author: jsonBook.Authors[0].Name,
		// TODO: []string of all categories ??
		Category: jsonBook.Subjects[0].Name,
	}

	return &form, nil
}
