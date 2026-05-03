package storage

import (
	"TicTacToe/internal/domain/board"
	"TicTacToe/internal/domain/consts"
	"TicTacToe/internal/domain/game"
	_interface "TicTacToe/internal/domain/game/interface"
	"TicTacToe/internal/domain/helpers"
	gameErrors "TicTacToe/internal/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

type GameDBModel struct {
	ID                    uuid.UUID  `db:"id"`
	Mode                  string     `db:"mode"`
	Board                 string     `db:"board"` // JSON
	CurrentTurn           string     `db:"current_turn"`
	Winner                string     `db:"winner"`
	Status                string     `db:"status"`
	Message               string     `db:"message"`
	Err                   string     `db:"err"`
	PlayerID              uuid.UUID  `db:"player_id"`
	Player1ID             uuid.UUID  `db:"player1_id"` // X
	Player2ID             *uuid.UUID `db:"player2_id"` // O (для player режима)
	Player1Mark           string     `db:"player1_mark"`
	Player2Mark           string     `db:"player2_mark"`
	PlayerMark            string     `db:"player_mark"`   // для bot режима
	ComputerMark          string     `db:"computer_mark"` // для bot режима
	ComputerLastMoveRow   *int       `db:"computer_last_move_row"`
	ComputerLastMoveCol   *int       `db:"computer_last_move_col"`
	ComputerLastMoveScore *int       `db:"computer_last_move_score"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

type UserDbModel struct {
	ID       uuid.UUID `db:"id"`
	Login    uuid.UUID `db:"login"`
	Password uuid.UUID `db:"password"`
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	storage := &Storage{
		pool: pool,
	}
	err := storage.InitStorage()
	if err != nil {
		panic(err)
	}
	return storage
}

func (s *Storage) gameToBotModel(g *game.GameVersusBot) (*GameDBModel, error) {
	gameBoard, err := json.Marshal(g.GameBoard)
	if err != nil {
		return nil, err
	}
	var lastMoveRow, lastMoveCol, lastMoveScore *int
	if g.LastComputerMove.Row != 0 && g.LastComputerMove.Col != 0 {
		lastMoveRow = &g.LastComputerMove.Row
		lastMoveCol = &g.LastComputerMove.Col
		lastMoveScore = &g.LastComputerMove.Score
	}

	return &GameDBModel{
		ID:                    g.ID,
		Mode:                  consts.BotMode,
		Board:                 string(gameBoard),
		CurrentTurn:           helpers.ConvertPlayer(g.CurrentTurn),
		Winner:                helpers.ConvertPlayer(g.Winner),
		Status:                helpers.ConvertStatus(g.Status),
		Message:               g.Message,
		PlayerID:              g.PlayerID,
		Player1ID:             g.PlayerID,
		PlayerMark:            board.NumToSymbol(g.PlayerMark),
		ComputerMark:          board.NumToSymbol(g.ComputerMark),
		ComputerLastMoveRow:   lastMoveRow,
		ComputerLastMoveCol:   lastMoveCol,
		ComputerLastMoveScore: lastMoveScore,
		CreatedAt:             g.CreatedAt,
		UpdatedAt:             g.UpdatedAt,
	}, nil
}

func (s *Storage) gameToPlayerModel(g *game.GameVersusPlayer) (*GameDBModel, error) {
	gameBoard, err := json.Marshal(g.GameBoard)
	if err != nil {
		return nil, err
	}

	var player2ID *uuid.UUID
	if g.SecondPlayerID != uuid.Nil {
		player2ID = &g.SecondPlayerID
	}

	return &GameDBModel{
		ID:          g.ID,
		Mode:        consts.PlayerMode,
		Board:       string(gameBoard),
		CurrentTurn: helpers.ConvertPlayer(g.CurrentTurn),
		Winner:      helpers.ConvertPlayer(g.Winner),
		Status:      helpers.ConvertStatus(g.Status),
		PlayerID:    g.FirstPlayerID,
		Player1ID:   g.FirstPlayerID,
		Player2ID:   player2ID,
		Player1Mark: board.NumToSymbol(g.Player1Mark),
		Player2Mark: board.NumToSymbol(g.Player2Mark),
		Err:         g.Err,
		Message:     g.Message,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}, nil
}

func (s *Storage) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return s.pool.Exec(ctx, sql, args...)
}

func (s *Storage) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return s.pool.Query(ctx, sql, args...)
}

func (s *Storage) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return s.pool.QueryRow(ctx, sql, args...)
}

func (s *Storage) SaveGame(ctx context.Context, g _interface.Game) error {
	op := "storage.SaveGame"
	var dbModel *GameDBModel
	switch gm := g.(type) {
	case *game.GameVersusBot:
		dbModelBot, err := s.gameToBotModel(gm)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
		dbModel = dbModelBot
	case *game.GameVersusPlayer:
		dbModelPlayer, err := s.gameToPlayerModel(gm)
		if err != nil {
			return fmt.Errorf("%s:%w", op, err)
		}
		dbModel = dbModelPlayer
	default:
		return fmt.Errorf("incorrect type for Game interface")
	}

	query := `
        INSERT INTO games (
            id, mode, board, current_turn, winner, status, err, message,
			player_id, player1_id, player2_id,   
            player_mark, computer_mark, player1_mark, player2_mark, computer_last_move_row,
            computer_last_move_col, computer_last_move_score,
            created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7,
            $8, $9, $10, $11, $12, $13, $14,
           	$15, $16, $17, $18, $19, $20
        )
    `

	_, err := s.pool.Exec(ctx, query,
		dbModel.ID, dbModel.Mode, dbModel.Board, dbModel.CurrentTurn, dbModel.Winner,
		dbModel.Status, dbModel.Err, dbModel.Message, dbModel.PlayerID, dbModel.Player1ID, dbModel.Player2ID,
		dbModel.PlayerMark, dbModel.ComputerMark, dbModel.Player1Mark, dbModel.Player2Mark,
		dbModel.ComputerLastMoveRow, dbModel.ComputerLastMoveCol, dbModel.ComputerLastMoveScore,
		dbModel.CreatedAt, dbModel.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (s *Storage) DeleteGame(ctx context.Context, gameID string) error {
	_, err := s.GetGameByID(ctx, gameID)
	if err != nil {
		return err
	}
	query := `
	DELETE * FROM tictactoe_db
	WHERE id = $1
`
	tag, err := s.pool.Exec(ctx, query, gameID)

	if tag.RowsAffected() == 0 {
		return gameErrors.ErrGameNotFound
	}

	return nil
}

func (s *Storage) CheckLogin(ctx context.Context, login string) (bool, error) {
	query := `
	SELECT COUNT(*) FROM players
	WHERE login = $1
`
	var count int
	err := s.QueryRow(ctx, query, login).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Storage) GetGameByID(ctx context.Context, gameID string) (_interface.Game, error) {
	query := `
        SELECT *
        FROM games 
        WHERE id = $1
    `

	rows, err := s.pool.Query(ctx, query, gameID)
	if err != nil {
		return nil, err
	}

	dbModel, err := pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[GameDBModel],
	)

	if err != nil {
		return nil, err
	}

	var gameBoard board.Board
	err = json.Unmarshal([]byte(dbModel.Board), &gameBoard)
	if err != nil {
		return nil, err
	}

	switch dbModel.Mode {
	case consts.BotMode:
		domainGame := &game.GameVersusBot{
			ID:           dbModel.ID,
			GameBoard:    gameBoard,
			CurrentTurn:  helpers.InverseConvertPlayer(dbModel.CurrentTurn),
			PlayerMark:   board.SymbolToNum(dbModel.PlayerMark),
			ComputerMark: board.SymbolToNum(dbModel.ComputerMark),
			Status:       helpers.InverseConvertStatus(dbModel.Status),
			Winner:       helpers.InverseConvertPlayer(dbModel.Winner),
			Message:      dbModel.Message,
			Err:          dbModel.Err,
			CreatedAt:    dbModel.CreatedAt,
			UpdatedAt:    dbModel.UpdatedAt,
		}
		if dbModel.ComputerLastMoveRow != nil && dbModel.ComputerLastMoveCol != nil && dbModel.ComputerLastMoveScore != nil {
			domainGame.LastComputerMove = game.Move{
				*dbModel.ComputerLastMoveRow,
				*dbModel.ComputerLastMoveCol,
				*dbModel.ComputerLastMoveScore,
			}
		}
		return domainGame, nil

	case consts.PlayerMode:
		var secondPlayerID uuid.UUID
		if dbModel.Player2ID != nil {
			secondPlayerID = *dbModel.Player2ID
		} else {
			secondPlayerID = uuid.Nil
		}

		domainGame := &game.GameVersusPlayer{
			ID:             dbModel.ID,
			GameBoard:      gameBoard,
			Winner:         helpers.InverseConvertPlayer(dbModel.Winner),
			Status:         helpers.InverseConvertStatus(dbModel.Status),
			CurrentTurn:    helpers.InverseConvertPlayer(dbModel.CurrentTurn),
			FirstPlayerID:  dbModel.Player1ID,
			SecondPlayerID: secondPlayerID,
			Player1Mark:    board.SymbolToNum(dbModel.Player1Mark),
			Player2Mark:    board.SymbolToNum(dbModel.Player2Mark),
			Err:            dbModel.Err,
			Message:        dbModel.Message,
			CreatedAt:      dbModel.CreatedAt,
			UpdatedAt:      dbModel.UpdatedAt,
		}
		return domainGame, nil
	default:
		return nil, fmt.Errorf("incorrect type of interface Game")
	}
}

func (s *Storage) UpdateGame(ctx context.Context, g _interface.Game) error {
	op := "storage.UpdateGame"

	var dbModel *GameDBModel

	switch gm := g.(type) {
	case *game.GameVersusBot:
		dbModelBot, err := s.gameToBotModel(gm)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		dbModel = dbModelBot
	case *game.GameVersusPlayer:
		dbModelPlayer, err := s.gameToPlayerModel(gm)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		dbModel = dbModelPlayer
	}

	query := `
	UPDATE games
	SET 
	    board = $1,
		current_turn = $2,
		winner = $3,
		status = $4,
		computer_last_move_row = $5,
		computer_last_move_col = $6,
		computer_last_move_score = $7,
		updated_at = $8,
		err = $9,
		message = $10,
-- 		player_id = $11,
-- 		player1_id = $12,
		player2_id = $11,
		player_mark = $12,
		player1_mark = $13,
		player2_mark = $14,
		computer_mark = $15
	WHERE id = $16
`

	tag, err := s.Exec(ctx, query, dbModel.Board, dbModel.CurrentTurn, dbModel.Winner, dbModel.Status,
		dbModel.ComputerLastMoveRow, dbModel.ComputerLastMoveCol, dbModel.ComputerLastMoveScore, time.Now(), dbModel.Err, dbModel.Message,
		dbModel.Player2ID, dbModel.PlayerMark, dbModel.Player1Mark, dbModel.Player2Mark, dbModel.ComputerMark, dbModel.ID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if tag.RowsAffected() == 0 {
		return gameErrors.ErrGameNotFound
	}

	return nil
}

func (s *Storage) InitStorage() error {
	query1 := `
	CREATE TABLE IF NOT EXISTS players (
		id UUID PRIMARY KEY,
		login TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)
`
	query2 := `
	CREATE TABLE IF NOT EXISTS games (
	    id UUID PRIMARY KEY,
	    mode TEXT NOT NULL,
	    board TEXT NOT NULL,
	    current_turn TEXT NOT NULL,
	    winner TEXT,
	    status TEXT NOT NULL,
	    message TEXT,
	    err TEXT,
	    player_id UUID REFERENCES players (id),
	    player1_id UUID REFERENCES players (id),
	    player2_id UUID REFERENCES players (id),
	    player1_mark TEXT,
	    player2_mark TEXT,
	    player_mark TEXT,
	    computer_mark TEXT,
	    computer_last_move_row INT,
	    computer_last_move_col INT,
	    computer_last_move_score INT,
	    created_at TIME,
	    updated_at TIME
	)
`
	query3 := `ALTER TABLE games ALTER COLUMN player_id DROP NOT NULL;`
	query4 := `ALTER TABLE games ALTER COLUMN player1_id DROP NOT NULL;`
	query5 := `ALTER TABLE games ALTER COLUMN player2_id DROP NOT NULL;`
	_, err := s.Exec(context.Background(), query1)
	if err != nil {
		return err
	}

	_, err = s.Exec(context.Background(), query2)
	if err != nil {
		return err
	}

	_, err = s.Exec(context.Background(), query3)
	if err != nil {
		return err
	}
	_, err = s.Exec(context.Background(), query4)
	if err != nil {
		return err
	}

	_, err = s.Exec(context.Background(), query5)
	if err != nil {
		return err
	}

	return nil
}
