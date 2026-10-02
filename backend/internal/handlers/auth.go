package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Tasks2/hmis/internal/response"
	"github.com/Tasks2/hmis/internal/services"
	"github.com/Tasks2/hmis/internal/validation"
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

	if !validation.Required(req.Email) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"EMAIL_REQUIRED",
			"email is required",
		)
		return
	}

	if !validation.Email(req.Email) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"INVALID_EMAIL",
			"email is invalid",
		)
		return
	}

	if !validation.Required(req.Password) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"PASSWORD_REQUIRED",
			"password is required",
		)
		return
	}

	if !validation.Required(req.FirstName) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"FIRST_NAME_REQUIRED",
			"first_name is required",
		)
		return
	}

	if !validation.Required(req.LastName) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"LAST_NAME_REQUIRED",
			"last_name is required",
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

	if errors.Is(err, services.ErrEmailExists) {
		response.JSONError(
			w,
			http.StatusConflict,
			"EMAIL_EXISTS",
			"email already exists",
		)
		return
	}

	if err != nil {
		http.Error(
			w,
			"could not create patient",
			http.StatusInternalServerError,
		)
		return
	}

	response.JSON(
		w,
		http.StatusCreated,
		map[string]string{
			"message": "patient registered successfully",
		},
	)
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

	if !validation.Required(req.Email) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"EMAIL_REQUIRED",
			"email is required",
		)
		return
	}

	if !validation.Email(req.Email) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"INVALID_EMAIL",
			"email is invalid",
		)
		return
	}

	if !validation.Required(req.Password) {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"PASSWORD_REQUIRED",
			"password is required",
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
