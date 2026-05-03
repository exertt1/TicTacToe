package di

import (
	"TicTacToe/internal/domain/repository"
	"TicTacToe/internal/handler"
	"TicTacToe/server/router"
	"TicTacToe/services/authorization"
	"TicTacToe/storage"
	"TicTacToe/storage/config"
	"context"
	"log"
	"net/http"
	"time"

	"go.uber.org/fx"
)

type App struct {
	Router *router.Router
}

func NewApp(router *router.Router) *App {
	return &App{
		Router: router,
	}
}

func (a *App) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    "localhost:8000",
		Handler: a.Router,
	}

	go func() {
		log.Println("Server starting on http://localhost:8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func (a *App) Stop(ctx context.Context) error {
	log.Println("App stopped")
	return nil
}

var Module = fx.Module("tictactoe",
	fx.Provide(
		config.NewDBConnection,
		fx.Annotate(
			storage.NewStorage,
			fx.As(new(repository.GameRepository)),
		),
		handler.NewGameHandler,
		router.NewRouter,
		authorization.NewUserService,
		NewApp,
	),
	fx.Invoke(func(lc fx.Lifecycle, app *App) {
		lc.Append(
			fx.Hook{
				OnStart: func(ctx context.Context) error {
					return app.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return app.Stop(ctx)
				},
			},
		)
	}),
)
