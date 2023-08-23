package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	errs "github.com/muhammadswa/personal-library/cmd/errors"
	"github.com/muhammadswa/personal-library/cmd/models"
	"github.com/muhammadswa/personal-library/cmd/repositories"
	"github.com/muhammadswa/personal-library/cmd/templates"
	"github.com/muhammadswa/personal-library/internal/database"
	"github.com/muhammadswa/personal-library/internal/validator"
)

type webBooksController struct {
	booksRespsitory *repositories.BooksRepository
	session         *scs.SessionManager
}

func NewBooksController(booksRespository *repositories.BooksRepository, session *scs.SessionManager) *webBooksController {
	return &webBooksController{
		booksRespsitory: booksRespository,
		session:         session,
	}
}

func (bc *webBooksController) CreateBook(w http.ResponseWriter, r *http.Request) {
	templateData := templates.NewTemplateData(bc.session, r)
	templateData.Form = &models.BookForm{}
	templates.Render(w, "create_book", templateData)
}

func (bc *webBooksController) CreateBookPost(w http.ResponseWriter, r *http.Request) {
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

	if !form.Valid() {
		data := templates.NewTemplateData(bc.session, r)
		data.Form = form
		templates.Render(w, "create_book", data)
		// c.templateCache.Render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// check if isbn is valid

	userId := bc.session.GetInt32(r.Context(), "userId")
	// for i := 0; i < 100; i++ {
	// 	_, _ = bc.booksRespsitory.CreateBook(r.Context(), database.CreateBookParams{
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
	id, err := bc.booksRespsitory.CreateBook(r.Context(), database.CreateBookParams{
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
	fmt.Fprintln(w, id)
}

func (bc *webBooksController) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	pageStr := httprouter.ParamsFromContext(r.Context()).ByName("page")
	// query := r.URL.Query().Get("q")
	if pageStr == "" {
		pageStr = "1"
	}
	// fmt.Println(query)
	fmt.Println(pageStr)

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errs.WebClientErr(w, "Error parsing offset")
		return
	}
	// TODO: check maximum offset can be reached, sql rows?

	booksLen, err := bc.booksRespsitory.GetBooksLength(r.Context())
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

	books, err := bc.booksRespsitory.GetBooks(r.Context(), int(offset))
	if err != nil {
		errs.WebServerErr(w, "err getting books")
		return
	}

	data := templates.NewTemplateData(bc.session, r)
	data.Books = &books
	data.IsPageNext = nextPage
	data.NextPage = page + 1

	if r.Header.Get("HX-Trigger") == "load-more-btn" {
		templates.RenderFragment(w, "books_list", data)
		return
	}
	templates.Render(w, "books", data)
}
