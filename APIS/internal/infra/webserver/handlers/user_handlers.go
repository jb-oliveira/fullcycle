package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jb-oliveira/fullcycle/APIS/pkg/log"

	"github.com/go-chi/jwtauth"
	"github.com/jb-oliveira/fullcycle/APIS/internal/dto"
	"github.com/jb-oliveira/fullcycle/APIS/internal/entity"
	"github.com/jb-oliveira/fullcycle/APIS/internal/infra/database"
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

// Auth godoc
// @Summary      Authenticate user
// @Description  Authenticate user with email and password and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginInput true "User login credentials"
// @Success      200 {object} dto.AuthResponse "Authentication successful"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /users/auth [post]
func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userLogin := &dto.LoginInput{}
	err := json.NewDecoder(r.Body).Decode(userLogin)
	if err != nil {
		log.Error(err.Error())
		ReturnHttpError(w, errors.New("Invalid request body"), http.StatusBadRequest)
		return
	}

	user, err := h.userDB.FindByEmail(userLogin.Email)
	if err != nil {
		log.Error(err.Error())
		ReturnHttpError(w, errors.New("User not found"), http.StatusNotFound)
		return
	}

	if !user.ValidatePassword(userLogin.Password) {
		log.Error("Invalid credentials")
		ReturnHttpError(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	accessToken := dto.AuthResponse{}

	_, token, err := h.jwtAuth.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"eml": user.Email,
		"exp": time.Now().Add(time.Second * time.Duration(h.jwtExpiration)).Unix(),
	})
	if err != nil {
		log.Error(err.Error())
		ReturnHttpError(w, errors.New("Failed to generate token"), http.StatusInternalServerError)
		return
	}

	accessToken.Token = token

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create User godoc
// @Summary      Create a new user
// @Description  Create a new user with name, email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUserInput true "User creation data"
// @Success      201 {object} dto.UserOutput "User created successfully"
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	inserInput := dto.CreateUserInput{}
	err := json.NewDecoder(r.Body).Decode(&inserInput)
	if err != nil {
		log.Error(err.Error())
		ReturnHttpError(w, errors.New("Invalid request body"), http.StatusBadRequest)
		return
	}
	user, err := entity.NewUser(inserInput.Name, inserInput.Email, inserInput.Password)
	if err != nil {
		log.Error(err.Error())
		ReturnHttpError(w, err, http.StatusBadRequest)
		return
	}
	err = h.userDB.Create(user)
	if err != nil {
		log.Error(err.Error())
		ReturnHttpError(w, errors.New("Failed to create user"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.UserOutput{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	})
}
