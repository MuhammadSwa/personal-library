package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) initServer() error {
	router := httprouter.New()
	srv := http.Server{
		Addr:    app.cfg.Port,
		Handler: router,
	}

	router.HandlerFunc(http.MethodGet, "/api/v1/health", app.healthCheckHandler)
	log.Printf("Starting server on port %s", srv.Addr)
	return srv.ListenAndServe()
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
