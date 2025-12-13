package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/users"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/http/handlers"
	httpMiddleware "github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/http/middleware"
)

func SetupRouter(userService *users.Service) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httpMiddleware.CORS)
	r.Use(httpMiddleware.Recover)

	// Auth handler
	authHandler := handlers.NewAuthHandler(userService)
	authMiddleware := httpMiddleware.NewAuthMiddleware(userService)

	// Public routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.RefreshToken)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/auth/me", authHandler.Me)
		})
	})

	return r
}

