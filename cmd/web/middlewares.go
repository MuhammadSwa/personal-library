package main

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/muhammadswa/personal-library/internal/logger"
	"github.com/rs/zerolog/hlog"
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

// func (m *Middleware) requestLogger(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
//
// 		defer func() {
// 			logger.Log.Info().
// 				Str("method", r.Method).
// 				Str("url", r.URL.RequestURI()).
// 				Str("user_agent", r.UserAgent()).
// 				Dur("elapsed_ms", time.Since(start)).
// 				Msg("incoming request")
// 		}()
//
// 		next.ServeHTTP(w, r)
// 	})
// }

func (m *Middleware) requestLogger(next http.Handler) http.Handler {

	h := hlog.NewHandler(logger.Log)

	accessHandler := hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status_code", status).
				Int("response_size_bytes", size).
				Dur("elapsed_ms", duration).
				Msg("incoming request")
		},
	)

	userAgentHandler := hlog.UserAgentHandler("http_user_agent")

	return h(accessHandler(userAgentHandler(next)))
}
