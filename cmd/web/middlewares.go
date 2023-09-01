package main

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Middleware struct {
	session *scs.SessionManager
}

func NewMiddleware(session *scs.SessionManager) *Middleware {
	return &Middleware{
		session: session,
	}
}

func (m *Middleware) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.session.Exists(r.Context(), "authenticatedUserID") {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) isAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.session.Exists(r.Context(), "authenticatedUserID") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
