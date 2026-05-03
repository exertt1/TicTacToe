package game

import (
	"TicTacToe/internal/domain/board"
	"TicTacToe/internal/domain/consts"
	"time"

	"github.com/google/uuid"
)

type GameVersusPlayer struct {
	ID             uuid.UUID
	Mode           uint
	GameBoard      board.Board
	CurrentTurn    uint
	Player1Mark    uint
	Player2Mark    uint
	Status         uint
	Winner         uint
	Message        string
	Err            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FirstPlayerID  uuid.UUID
	SecondPlayerID uuid.UUID
}

type DataGameVersusPlayer struct {
	ID             uuid.UUID `db:"id"`
	GameBoard      string    `db:"board"`
	CurrentTurn    uint      `db:"current_turn"`
	Status         int       `db:"status"`
	Winner         int       `db:"winner"`
	Message        string    `db:"msg"`
	Err            string    `db:"err"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	FirstPlayerID  uuid.UUID `db:"first_player_id"`
	SecondPlayerID uuid.UUID `db:"second_player_id"`
	Mode           int       `db:"mode"`
}

func (g *GameVersusPlayer) CheckWinner(mark uint) bool {
	return g.GameBoard.CheckWin(mark)
}

func (g *GameVersusPlayer) IsDraw() bool {
	return g.GameBoard.IsDraw()
}

func (g *GameVersusPlayer) BoardToSymbols() [consts.Rows][consts.Columns]string {
	return g.GameBoard.BoardToSymbols()
}

func (g *GameVersusPlayer) FindEmptyCells() [][]uint {
	return g.GameBoard.FindEmptyCells()
}

func (g *GameVersusPlayer) MakeMove(row, col uint, player uint) error {
	return g.GameBoard.MakeMove(row, col, player)
}

func (g *GameVersusPlayer) IsFinished() bool {
	if res := g.GameBoard.IsFinished(); res {
		g.Status = consts.Finished
	}
	return g.GameBoard.IsFinished()
}
