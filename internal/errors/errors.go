package errs

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func ServerError(w http.ResponseWriter, err error) {
	log.Debug().Stack().Err(err).Msg("stack trace")
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
