package controllers

// type BooksController struct {
// 	repos *repositories.
// }
//
// func NewBooksController(booksRespository *repositories.BooksRepository) *BooksController {
// 	return &BooksController{
// 		booksRespsitory: booksRespository,
// 	}
// }
//
// // TODO: Check if book exist
// // TODO: change to isbn
// func (bc *BooksController) GetBookByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	// isbn := ps.ByName("isbn")
// 	id := 1
// 	book, err := bc.booksRespsitory.GetBookByID(r.Context(), id)
// 	if err != nil {
// 		ErrReponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	respondWithJSON(w, http.StatusOK, book)
// }
//
// func (bc *BooksController) CreateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	var book database.CreateBookParams
// 	err := json.NewDecoder(r.Body).Decode(&book)
// 	if err != nil {
// 		ErrReponse(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	_, err = bc.booksRespsitory.CreateBook(r.Context(), book)
// 	if err != nil {
// 		ErrReponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// }
//
// func (bc *BooksController) CreateBookWithISBN(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	isbn := ps.ByName("isbn")
// 	// https://openlibrary.org/books/OL7353617M.json
//
// 	// https://openlibrary.org/api/books?bibkeys=ISBN:9780980200447&jscmd=data&format=json
// 	openLibraryUrl := fmt.Sprintf("https://openlibrary.org/isbn/%s.json", isbn)
// 	resp, err := http.Get(openLibraryUrl)
// 	if err != nil {
// 		ErrReponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	defer resp.Body.Close()
// 	openLibBook := struct {
// 		Isbn        []string `json:"isbn_13"`
// 		PagesNumber int      `json:"number_of_pages"`
// 		Title       string   `json:"title"`
// 		PublishDate string   `json:"publish_date"`
// 		Publishers  []string `json:"publishers"`
// 		Authors     []struct {
// 			Key string `json:"key"`
// 		} `json:"authors"`
// 	}{}
// 	err = json.NewDecoder(resp.Body).Decode(&openLibBook)
// 	if err != nil {
// 		ErrReponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	// TODO: parse authors
// 	// var authors []string
// 	publishDate, _ := strconv.ParseInt(openLibBook.PublishDate, 10, 64)
// 	book := database.CreateBookParams{
// 		Isbn:             isbn,
// 		Title:            openLibBook.Title,
// 		Author:           openLibBook.Authors[0].Key,
// 		Category:         "",
// 		Publisher:        arrToStr(openLibBook.Publishers),
// 		YearOfPublishing: int32(publishDate),
// 		Img:              "",
// 		NumberOfPages:    int32(openLibBook.PagesNumber),
// 		PersonalNotes:    "",
// 		PersonalRating:   0,
// 		ReadStatus:       false,
// 		ReadDate:         time.Now(),
// 	}
// 	_, err = bc.booksRespsitory.CreateBook(r.Context(), book)
// 	if err != nil {
// 		ErrReponse(w, http.StatusInternalServerError, err.Error())
// 		return
// 	}
// }
//
// func mapValueTostr(m map[string]string) string {
// 	var str string
// 	for _, v := range m {
// 		str += v + ","
// 	}
// 	return str
// }
//
// func arrToStr(arr []string) string {
// 	var str string
// 	for i, v := range arr {
// 		if i == len(arr)-1 {
// 			str += v
// 			break
// 		}
// 		str += v + ","
// 	}
// 	return str
// }
