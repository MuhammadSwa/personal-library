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

func InitHttpServer(conn *sql.DB, port string) error {
	dbQueries := database.New(conn)

	// initialize a sessionManager
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(conn)
	// TODO: make an check box (sign in for 30 days) when logging in
	sessionManager.Lifetime = 12 * time.Hour

	repos := repositories.New(dbQueries)
	controllers := controllers.New(repos, sessionManager)

	middleware := NewMiddleware(sessionManager)
	dynamicMiddleware := alice.New(sessionManager.LoadAndSave)
	protectedRoutes := dynamicMiddleware.Append(middleware.authMiddleware)
	isAuthenticated := dynamicMiddleware.Append(middleware.isAuthenticated)

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(controllers.NotFound)
	// serve static files
	router.ServeFiles("/static/*filepath", http.Dir("./web/static/"))

	// unprotected routes
	router.Handler(http.MethodGet, "/login", isAuthenticated.ThenFunc(controllers.Login))
	router.Handler(http.MethodGet, "/register", isAuthenticated.ThenFunc(controllers.Register))
	router.GET("/", controllers.Home)
	router.POST("/login", controllers.LoginPost)
	router.POST("/register", controllers.RegisterPost)
	router.POST("/logout", controllers.LogoutPost)
	//// protected routes
	router.Handler(http.MethodGet, "/create", protectedRoutes.ThenFunc(controllers.CreateBook))
	router.Handler(http.MethodGet, "/book/:id", protectedRoutes.ThenFunc(controllers.GetBookByID))
	router.Handler(http.MethodDelete, "/book/:id", protectedRoutes.ThenFunc(controllers.DeleteBook))
	router.Handler(http.MethodGet, "/edit/book/:id", protectedRoutes.ThenFunc(controllers.EditBook))
	router.Handler(http.MethodPut, "/edit/book/:id", protectedRoutes.ThenFunc(controllers.EditBookPut))
	router.Handler(http.MethodGet, "/books/:page", protectedRoutes.ThenFunc(controllers.GetAllBooks))
	router.Handler(http.MethodGet, "/books", protectedRoutes.ThenFunc(controllers.GetAllBooks))
	router.Handler(http.MethodPost, "/create", protectedRoutes.ThenFunc(controllers.CreateBookPost))
	router.HandlerFunc(http.MethodGet, "/fetchByIsbn", controllers.FetchByIsbn)
	// router.Handler(http.MethodPost, "/fetchByIsbn", protectedRoutes.ThenFunc(controllers.CreateBookWithISBN))
	router.Handler(http.MethodGet, "/profile", protectedRoutes.ThenFunc(controllers.Profile))

	mainMiddleware := sessionManager.LoadAndSave(router)

	return http.ListenAndServe(port, mainMiddleware)

}
