package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/muhammadswa/personal-library/internal/database"
	"github.com/muhammadswa/personal-library/pkg/controllers"
	"github.com/muhammadswa/personal-library/pkg/repositories"
)

type httpServer struct {
	booksController *controllers.BooksController
	booksRepository *repositories.BooksRepository
	router          *http.Handler
	port            string
}

func InitHttpServer(conn *sql.DB, port string) *httpServer {
	dbQueries := database.New(conn)

	// initialize a sessionManager
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(conn)
	// TODO: make an check box (sign in for 30 days) when logging in
	sessionManager.Lifetime = 12 * time.Hour

	// api
	booksRespository := repositories.NewBooksRepository(dbQueries)
	bookController := controllers.NewBooksController(booksRespository, sessionManager)
	// web
	webBooksController := controllers.NewBooksController(booksRespository, sessionManager)
	usersRepository := repositories.NewUsersRepository(dbQueries)
	usersController := controllers.NewUsersController(usersRepository, sessionManager)
	staticController := controllers.NewStaticController(booksRespository, sessionManager)

	middleware := NewMiddleware(sessionManager)
	dynamicMiddleware := alice.New(sessionManager.LoadAndSave)
	protectedRoutes := dynamicMiddleware.Append(middleware.authMiddleware)

	router := httprouter.New()
	// serve static files
	router.ServeFiles("/static/*filepath", http.Dir("./web/static/"))

	// api routes
	// router.HandlerFunc()
	// router.GET("/v1/book/:isbn", bookController.GetBookByID)
	// router.POST("/v1/book", bookController.CreateBook)
	// router.POST("/v1/book/:isbn", bookController.CreateBookWithISBN)

	// web routes
	router.GET("/", staticController.Home)
	router.GET("/login", usersController.Login)
	router.POST("/login", usersController.LoginPost)
	router.GET("/signup", usersController.SignUp)
	router.POST("/signup", usersController.SignupPost)
	router.POST("/logout", usersController.LogoutPost)
	//// protected routes
	router.Handler(http.MethodGet, "/create", protectedRoutes.ThenFunc(webBooksController.CreateBook))
	router.Handler(http.MethodGet, "/book/:id", protectedRoutes.ThenFunc(webBooksController.GetBookByID))
	router.Handler(http.MethodDelete, "/book/:id", protectedRoutes.ThenFunc(webBooksController.DeleteBook))
	router.Handler(http.MethodGet, "/edit/book/:id", protectedRoutes.ThenFunc(webBooksController.EditBook))
	router.Handler(http.MethodPut, "/edit/book/:id", protectedRoutes.ThenFunc(webBooksController.EditBookPut))
	router.Handler(http.MethodGet, "/books/:page", protectedRoutes.ThenFunc(webBooksController.GetAllBooks))
	router.Handler(http.MethodGet, "/books", protectedRoutes.ThenFunc(webBooksController.GetAllBooks))
	router.Handler(http.MethodPost, "/create", protectedRoutes.ThenFunc(webBooksController.CreateBookPost))
	router.Handler(http.MethodGet, "/profile", protectedRoutes.ThenFunc(staticController.Profile))

	mainMiddleware := sessionManager.LoadAndSave(router)

	return &httpServer{
		booksController: bookController,
		booksRepository: booksRespository,
		router:          &mainMiddleware,
		port:            port,
	}
}

func (hs *httpServer) Run() error {
	err := http.ListenAndServe(hs.port, *hs.router)
	if err != nil {
		return err
	}
	return nil
}