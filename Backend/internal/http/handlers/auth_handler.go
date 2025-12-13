package handlers

import (
	"net/http"

	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/users"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/http/httpx"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/context"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
)

type AuthHandler struct {
	userService *users.Service
}

func NewAuthHandler(userService *users.Service) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	User  *users.User    `json:"user"`
	Token *users.TokenPair `json:"token,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := httpx.BindJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.Register(r.Context(), users.CreateUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrInvalidInput || err == errors.ErrAlreadyExists {
			status = http.StatusBadRequest
		}
		httpx.Error(w, status, err)
		return
	}

	httpx.JSON(w, http.StatusCreated, AuthResponse{User: user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := httpx.BindJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.userService.Login(r.Context(), users.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrUnauthorized {
			status = http.StatusUnauthorized
		}
		httpx.Error(w, status, err)
		return
	}

	httpx.JSON(w, http.StatusOK, AuthResponse{Token: token})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := context.GetUserID(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, errors.ErrUnauthorized)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, err)
		return
	}

	httpx.JSON(w, http.StatusOK, AuthResponse{User: user})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := httpx.BindJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, errors.ErrInvalidInput)
		return
	}

	token, err := h.userService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrUnauthorized {
			status = http.StatusUnauthorized
		}
		httpx.Error(w, status, err)
		return
	}

	httpx.JSON(w, http.StatusOK, AuthResponse{Token: token})
}

