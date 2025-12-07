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

// Auth godoc
// @Summary      Authenticate user
// @Description  Authenticate user with email and password and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginInput true "User login credentials"
// @Success      200 {object} object{access_token=string} "Authentication successful"
// @Failure      400 {string} string "Invalid request body"
// @Failure      401 {string} string "Invalid credentials"
// @Failure      404 {string} string "User not found"
// @Failure      500 {string} string "Failed to generate token"
// @Router       /users/auth [post]
func (h *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	userLogin := &dto.LoginInput{}
	err := json.NewDecoder(r.Body).Decode(userLogin)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userDB.FindByEmail(userLogin.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !user.ValidatePassword(userLogin.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	accessToken.AccessToken = token

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
// @Failure      400 {string} string "Invalid request body or validation error"
// @Failure      500 {string} string "Failed to create user"
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	inserInput := dto.CreateUserInput{}
	err := json.NewDecoder(r.Body).Decode(&inserInput)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := entity.NewUser(inserInput.Name, inserInput.Email, inserInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.userDB.Create(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
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
