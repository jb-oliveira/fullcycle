package handlers

import "net/http"

func ReturnHttpError(w http.ResponseWriter, error error, code int) {
	http.Error(w, error.Error(), code)
}
