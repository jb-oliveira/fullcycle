package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/dto"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/infra/database"
)

type UserHandler struct {
	userDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{
		userDB: userDB,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	inserInput := dto.CreateUserInput{}
	json.NewDecoder(r.Body).Decode(&inserInput)
	user, err := entity.NewUser(inserInput.Name, inserInput.Email, inserInput.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.userDB.Create(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.UserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	})
}
