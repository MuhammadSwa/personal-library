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
	// err := r.ParseForm()
	// if err != nil {
	// 	errs.WebClientErr(w, "Error parsing form")
	// 	return
	// }
	// change from map[string][]string to map[string]string
	// formMap := make(map[string]string, len(r.Form))
	// for k, v := range r.Form {
	// 	formMap[k] = v[0]
	// }

	// yearOfPublishing, _ := strconv.Atoi(formMap["year_of_publishing"])
	// numOfPages, _ := strconv.Atoi(formMap["number_of_pages"])
	// form := &models.BookForm{
	// 	Isbn:             r.Form.Get("isbn"),
	// 	Title:            r.Form.Get("title"),
	// 	Author:           r.Form.Get("author"),
	// 	Category:         r.Form.Get("category"),
	// 	Publisher:        r.Form.Get("publisher"),
	// 	Img:              r.Form.Get("img"),
	// 	YearOfPublishing: int32(yearOfPublishing),
	// 	NumberOfPages:    int32(numOfPages),
	// }
	//
	// form.Isbn = r.Form.Get("isbn")
	// form.Title = r.Form.Get("title")
	// form.Author = r.Form.Get("author")
	// templateData := templates.NewTemplateData(bc.session, r)
	// templateData.Form = form
	// book := database.Book{
	// 	Isbn:             form.Isbn,
	// 	Title:            form.Title,
	// 	Author:           form.Author,
	// 	Category:         form.Category,
	// 	Publisher:        form.Publisher,
	// 	YearOfPublishing: form.YearOfPublishing,
	// 	Img:              form.Img,
	// 	NumberOfPages:    form.NumberOfPages,
	// 	PersonalRating:   form.PersonalRating,
	// 	PersonalNotes:    form.PersonalNotes,
	// 	ReadStatus:       form.ReadStatus,
	// 	ReadDate:         form.ReadDate,
	// }
	//
	// templateData.Book = &book

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
	// TODO: try to remove it ? do we need it?
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
		errs.WebServerErr(w, "Error creating book")
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/book/%d", id), http.StatusSeeOther)
}

func (bc *Controllers) FetchByIsbn(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Query().Get("isbn")
	form, err := getBookByIsbn(isbn)
	if err != nil {
		errs.WebServerErr(w, "Error fetching book")
		return
	}
	data := templates.NewTemplateData(bc.session, r)
	data.Form = form

	templates.RenderFragment(w, "book_form", data)

	// construct query string
	// q := url.Values{
	// 	"isbn":               {isbn},
	// 	"title":              {book.Title},
	// 	"author":             {book.Authors[0].Name},
	// 	"number_of_pages":    {strconv.Itoa(book.NumberOfPages)},
	// 	"year_of_publishing": {book.PublishDate},
	// 	"img":                {book.Cover.Large},
	// 	"publisher":          {book.Publishers[0].Name},
	// }
	// createUrl := fmt.Sprintf("/create?%s", q.Encode())
	//
	// http.Redirect(w, r, createUrl, http.StatusSeeOther)
}

// put this is it's own file and package? in pkg/api?
func getBookByIsbn(isbn string) (*models.BookForm, error) {

	// TODO: make sure isbn is valid (13 digits) and not empty , no - or spaces ?
	// maybe - is okay? ? use regex client side
	openLibraryUrl := fmt.Sprintf("https://openlibrary.org/api/books?bibkeys=ISBN:%s&jscmd=data&format=json", isbn)
	resp, err := http.Get(openLibraryUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: use form/playground??

	// I'm using this weird way of decoding because the stupid way the json from openlibrary is structured
	openLibRes := models.OpenLibResponse{}
	err = json.NewDecoder(resp.Body).Decode(&openLibRes)
	if err != nil {
		return nil, err
	}

	key := "ISBN:" + isbn
	jsonBook := openLibRes[key]

	if err != nil {
		return nil, err
	}
	// Dec 07, 2018
	splits := strings.Split(jsonBook.PublishDate, ", ")[1]
	publishDate, err := strconv.Atoi(splits)
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
