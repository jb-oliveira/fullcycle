package handlers

import "net/http"

func ReturnHttpError(w http.ResponseWriter, error string, code int) {
	http.Error(w, error, code)
}
