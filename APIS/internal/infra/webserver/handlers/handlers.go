package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jb-oliveira/fullcycle/APIS/internal/dto"
)

func ReturnHttpError(w http.ResponseWriter, err error, code int) {
	ReturnHttpErrors(w, []error{err}, code)
}

func ReturnHttpErrors(w http.ResponseWriter, errors []error, code int) {
	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Error())
	}
	newVar := dto.ErrorResponse{
		Messages: messages,
		Code:     code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(newVar)
}
