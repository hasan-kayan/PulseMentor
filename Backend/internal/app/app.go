package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/config"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/domain/users"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/http/routes"
	"github.com/hasan-kayan/PulseMentor/tree/main/Backend/internal/infra/db/postgres"
)

type App struct {
	Router *chi.Mux
	db     *postgres.DB
}

func New(cfg config.Config) (*App, error) {
	ctx := context.Background()

	// Initialize database
	db, err := postgres.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	usersRepo := postgres.NewUsersRepository(db)

	// Initialize services
	userService := users.NewService(usersRepo, cfg)

	// Setup router
	router := routes.SetupRouter(userService)

	return &App{
		Router: router,
		db:     db,
	}, nil
}

func (a *App) Close(ctx context.Context) error {
	if a.db != nil {
		a.db.Close()
	}
	return nil
}

