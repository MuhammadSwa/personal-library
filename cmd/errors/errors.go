package errs

import "net/http"

func WebServerErr(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

func WebClientErr(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}
