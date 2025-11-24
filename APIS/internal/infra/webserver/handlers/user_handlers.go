package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/dto"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/infra/database"
)

type UserHandler struct {
	userDB        database.UserInterface
	jwtAuth       *jwtauth.JWTAuth
	jwtExpiration int
}

func NewUserHandler(userDB database.UserInterface, jwtAuth *jwtauth.JWTAuth, expiration int) *UserHandler {
	return &UserHandler{
		userDB:        userDB,
		jwtAuth:       jwtAuth,
		jwtExpiration: expiration,
	}
}

func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userLogin := &dto.LoginInput{}
	err := json.NewDecoder(r.Body).Decode(userLogin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userDB.FindByEmail(userLogin.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !user.ValidatePassword(userLogin.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{}

	_, token, err := h.jwtAuth.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"eml": user.Email,
		"exp": time.Now().Add(time.Second * time.Duration(h.jwtExpiration)).Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accessToken.AccessToken = token

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
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
