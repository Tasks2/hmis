package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tasks2/hmis/internal/services"
)

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthHandler struct {
	authService *services.AuthService
	jwtService  *services.JWTService
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authService *services.AuthService, jwtService *services.JWTService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" ||
		req.Password == "" ||
		req.FirstName == "" ||
		req.LastName == "" {
		http.Error(
			w,
			"email, password, first_name and last_name are required",
			http.StatusBadRequest,
		)
		return
	}

	err := h.authService.CreatePatient(
		r.Context(),
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
	)

	if err != nil {
		http.Error(
			w,
			"could not create patient",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "patient registered successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(
			w,
			"email and password are required",
			http.StatusBadRequest,
		)
		return
	}

	userID, role, err := h.authService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			http.Error(
				w,
				"invalid credentials",
				http.StatusUnauthorized,
			)
			return
		}

		http.Error(
			w,
			"could not login",
			http.StatusInternalServerError,
		)
		return
	}

	token, err := h.jwtService.GenerateToken(
		userID,
		role,
	)

	if err != nil {
		http.Error(
			w,
			"could not generate token",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
