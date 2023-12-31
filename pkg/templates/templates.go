package templates

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/muhammadswa/personal-library/internal/database"
	errs "github.com/muhammadswa/personal-library/internal/errors"
	"github.com/muhammadswa/personal-library/pkg/models"
)

type templateData struct {
	// ValidationErrors map[string]string
	IsAuthenticated bool
	Book            *database.Book
	Books           *[]database.Book
	Form            any
	IsPageNext      bool
	NextPage        int
	Query           string
}

func New(session *scs.SessionManager, r *http.Request) *templateData {
	return &templateData{
		IsAuthenticated: session.Exists(r.Context(), "authenticatedUserID"),
		Form:            models.BookForm{},
	}
}

func Render(w http.ResponseWriter, page string, data any) {
	ts, err := template.New(page).Funcs(functions).ParseGlob("./web/html/**/*.tmpl.html")
	if err != nil {
		errs.ServerError(w, err)
		return
	}
	files := []string{
		"./web/html/base.tmpl.html",
		fmt.Sprintf("./web/html/pages/%s.tmpl.html", page),
	}
	ts, err = ts.ParseFiles(files...)
	if err != nil {
		errs.ServerError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		errs.ServerError(w, err)
		return
	}
}

func RenderFragment(w http.ResponseWriter, page string, data any) {
	ts, err := template.ParseFiles(fmt.Sprintf("./web/html/fragments/%s.tmpl.html", page))
	if err != nil {
		errs.ServerError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, page, data)
	if err != nil {
		errs.ServerError(w, err)
		return
	}
}

func prevOffset(offset int64) int64 {
	offset = offset - 2
	if offset <= 0 {
		return 1
	}
	return offset
}

var functions = template.FuncMap{
	"prevOffset": prevOffset,
}

// type TemplateCache map[string]*template.Template
//
// func NewTemplateCache() (*TemplateCache, error) {
// 	cache := &TemplateCache{}
// 	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, page := range pages {
// 		// exract the file name (like home.tmpl) from the full path
// 		name := filepath.Base(page)
//
// 		// parse the base template
// 		ts, err := template.ParseFiles("./ui/html/base.tmpl")
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// parse all the partials and add them to the base template
// 		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// parse the page template and add it to the base and partials
// 		ts, err = ts.ParseFiles(page)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// add the template to the cache
// 		(*cache)[name] = ts
// 	}
//
// 	return cache, nil
// }
//
// func (tc *TemplateCache) Render(w http.ResponseWriter, status int, pageName string, data any) {
// 	// get the page template from the cache
// 	ts, ok := (*tc)[pageName]
// 	if !ok {
// 		helpers.WebServerErr(w, fmt.Sprintf("page %s does not exist", pageName))
// 		return
// 	}
//
// 	buf := new(bytes.Buffer)
// 	// Write the template to the buffer, instead of straight to the
// 	// http.ResponseWriter. If there's an error, call our serverError() helper
// 	// and then return.
// 	err := ts.ExecuteTemplate(buf, "base", data)
// 	if err != nil {
// 		helpers.WebServerErr(w, "err executing template")
// 		return
// 	}
//
// 	w.WriteHeader(status)
//
// 	buf.WriteTo(w)
// }
