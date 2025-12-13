package middleware

import (
	"net/http"
	"strings"

	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/users"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/http/httpx"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/context"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/shared/errors"
)

type AuthMiddleware struct {
	userService *users.Service
}

func NewAuthMiddleware(userService *users.Service) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.Error(w, http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpx.Error(w, http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		userID, err := m.userService.ValidateToken(parts[1])
		if err != nil {
			httpx.Error(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithUserID(r.Context(), userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

