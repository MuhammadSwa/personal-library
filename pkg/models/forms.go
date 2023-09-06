package models

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/form"
	"github.com/muhammadswa/personal-library/internal/validator"
)

type LoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-`
}

type RegisterForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	Username            string `form:"username"`
	validator.Validator `form:"-`
}

type BookForm struct {
	ID               int32     `form:"id"`
	Isbn             string    `form:"isbn"`
	Title            string    `form:"title"`
	Author           string    `form:"author"`
	Category         string    `form:"category"`
	Publisher        string    `form:"publisher"`
	YearOfPublishing int32     `form:"year_of_publishing"`
	Img              string    `form:"img"`
	NumberOfPages    int32     `form:"number_of_pages"`
	PersonalRating   float64   `form:"personal_rating"`
	PersonalNotes    string    `form:"personal_notes"`
	ReadStatus       bool      `form:"read_status"`
	ReadDate         time.Time `form:"read_date"`
	// database.CreateBookParams
	validator.Validator `form:"-`
}

func DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	decoder := form.NewDecoder()
	decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		return time.Parse("2006-01-02", vals[0])
	}, time.Time{})

	if r.PostForm.Get("read_date") == "" {
		r.PostForm.Set("read_date", time.Now().Format("2006-01-02"))
	}
	err = decoder.Decode(&dst, r.PostForm)

	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, we return them as normal.
		fmt.Println(err)
		return err
	}
	return nil
}

type OpenLibResponse map[string]JSONBook

type JSONBook struct {
	URL     string `json:"url"`
	Key     string `json:"key"`
	Title   string `json:"title"`
	Authors []struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"authors"`
	NumberOfPages int    `json:"number_of_pages"`
	Pagination    string `json:"pagination"`
	Weight        string `json:"weight"`
	ByStatement   string `json:"by_statement"`
	Identifiers   struct {
		Amazon       []string `json:"amazon"`
		Google       []string `json:"google"`
		Librarything []string `json:"librarything"`
		Goodreads    []string `json:"goodreads"`
		Isbn10       []string `json:"isbn_10"`
		Isbn13       []string `json:"isbn_13"`
		Lccn         []string `json:"lccn"`
		Oclc         []string `json:"oclc"`
		Openlibrary  []string `json:"openlibrary"`
	} `json:"identifiers"`
	Classifications struct {
		LcClassifications []string `json:"lc_classifications"`
		DeweyDecimalClass []string `json:"dewey_decimal_class"`
	} `json:"classifications"`
	Publishers []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	PublishPlaces []struct {
		Name string `json:"name"`
	} `json:"publish_places"`
	PublishDate string `json:"publish_date"`
	Subjects    []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"subjects"`
	Notes           string `json:"notes"`
	TableOfContents []struct {
		Level   int    `json:"level"`
		Label   string `json:"label"`
		Title   string `json:"title"`
		Pagenum string `json:"pagenum"`
	} `json:"table_of_contents"`
	Links []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"links"`
	Ebooks []struct {
		PreviewURL   string `json:"preview_url"`
		Availability string `json:"availability"`
		Formats      struct {
		} `json:"formats"`
		BorrowURL  string `json:"borrow_url"`
		Checkedout bool   `json:"checkedout"`
	} `json:"ebooks"`
	Cover struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"cover"`
}
