package repository

import (
	_interface "TicTacToe/internal/domain/game/interface"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type GameRepository interface {
	CheckLogin(ctx context.Context, login string) (bool, error)

	SaveGame(ctx context.Context, g _interface.Game) error

	GetGameByID(ctx context.Context, gameID string) (_interface.Game, error)

	DeleteGame(ctx context.Context, gameID string) error

	UpdateGame(ctx context.Context, g _interface.Game) error

	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)

	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
