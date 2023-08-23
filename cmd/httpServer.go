package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	controllers "github.com/muhammadswa/personal-library/cmd/api"
	"github.com/muhammadswa/personal-library/cmd/repositories"
	webControllers "github.com/muhammadswa/personal-library/cmd/web/controllers"
	"github.com/muhammadswa/personal-library/internal/database"
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
	bookController := controllers.NewBooksController(booksRespository)
	// web
	webBooksController := webControllers.NewBooksController(booksRespository, sessionManager)
	usersRepository := repositories.NewUsersRepository(dbQueries)
	usersController := webControllers.NewUsersController(usersRepository, sessionManager)
	staticController := webControllers.NewStaticController(booksRespository, sessionManager)

	middleware := NewMiddleware(sessionManager)
	dynamicMiddleware := alice.New(sessionManager.LoadAndSave)
	protectedRoutes := dynamicMiddleware.Append(middleware.authMiddleware)

	router := httprouter.New()
	// serve static files
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static/"))

	// api routes
	// router.HandlerFunc()
	router.GET("/v1/book/:isbn", bookController.GetBookByID)
	router.POST("/v1/book", bookController.CreateBook)
	router.POST("/v1/book/:isbn", bookController.CreateBookWithISBN)

	// web routes
	router.GET("/", staticController.Home)
	router.GET("/login", usersController.Login)
	router.POST("/login", usersController.LoginPost)
	router.GET("/signup", usersController.SignUp)
	router.POST("/signup", usersController.SignupPost)
	router.POST("/logout", usersController.LogoutPost)
	//// protected routes
	router.Handler(http.MethodGet, "/create", protectedRoutes.ThenFunc(webBooksController.CreateBook))
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
