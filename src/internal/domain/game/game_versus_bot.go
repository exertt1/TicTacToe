package game

import (
	"TicTacToe/internal/domain/board"
	"TicTacToe/internal/domain/consts"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type GameVersusBot struct {
	ID               uuid.UUID
	Mode             uint
	GameBoard        board.Board
	PlayerMark       uint
	ComputerMark     uint
	CurrentTurn      uint
	LastComputerMove Move
	Status           uint
	Winner           uint
	Message          string
	Err              string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	PlayerID         uuid.UUID
}

type DataGame struct {
	ID                    string    `db:"id"`
	GameBoard             string    `db:"board"`
	PlayerMark            uint      `db:"player_mark"`
	ComputerMark          uint      `db:"computer_mark"`
	CurrentTurn           uint      `db:"current_turn"`
	LastComputerMoveRow   *int      `db:"last_computer_move_row"`
	LastComputerMoveCol   *int      `db:"last_computer_move_col"`
	LastComputerMoveScore *int      `db:"last_computer_move_score"`
	Status                int       `db:"status"`
	Winner                int       `db:"winner"`
	Message               string    `db:"msg"`
	Err                   string    `db:"err"`
	CreatedAt             time.Time `db:"created_at"`
	UpdatedAt             time.Time `db:"updated_at"`
	PlayerID              string    `db:"player_id"`
	Mode                  int       `db:"mode"`
}

func InitGameVersusBot(playerID uuid.UUID) *GameVersusBot {
	id, _ := uuid.NewV7()
	return &GameVersusBot{
		ID:        id,
		GameBoard: board.Board{},
		Status:    consts.Active,
		Winner:    consts.Nobody,
		PlayerID:  playerID,
	}
}

func InitGameVersusPlayer(playerID uuid.UUID) *GameVersusPlayer {
	id, _ := uuid.NewV7()
	return &GameVersusPlayer{
		ID:            id,
		GameBoard:     board.Board{},
		Status:        consts.AwaitingSecondPlayer,
		Winner:        consts.Nobody,
		FirstPlayerID: playerID,
	}
}

func (g *GameVersusBot) Game() {
	var (
		pos uint
		row uint
		col uint
	)
	for !g.IsFinished() {
		g.GameBoard.Draw()
		fmt.Println("Ваш ход: ")
		_, _ = fmt.Scan(&pos)
		_ = findXY(pos, &row, &col)
		err := g.MakeMove(row, col, consts.Cross)
		if err != nil {
			fmt.Println("ошибка!")
			continue
		}
		notFilledPoss := g.FindEmptyCells()
		var ran int
		if len(notFilledPoss) > 1 {
			ran = rand.Intn(len(notFilledPoss))
		} else {
			ran = 0
		}
		randY, randX := notFilledPoss[ran][0], notFilledPoss[ran][1]
		g.GameBoard[randY][randX] = consts.Zero
	}
	g.GameBoard.Draw()
	fmt.Println("игра закончена")
}

func (g *GameVersusBot) FindEmptyCells() [][]uint {
	return g.GameBoard.FindEmptyCells()
}

func (g *GameVersusBot) MakeMove(row, col uint, player uint) error {
	return g.GameBoard.MakeMove(row, col, player)
}

func (g *GameVersusBot) PlayerMove(row, col uint) error {
	err := g.MakeMove(row, col, g.PlayerMark)
	if err != nil {
		return err
	}
	g.CurrentTurn = consts.Computer
	return err
}

func (g *GameVersusBot) IsFinished() bool {
	if res := g.GameBoard.IsFinished(); res {
		g.Status = consts.Finished
	}
	return g.GameBoard.IsFinished()
}

func (g *GameVersusBot) CheckWinner(mark uint) bool {
	return g.GameBoard.CheckWin(mark)
}

func (g *GameVersusBot) IsDraw() bool {
	return g.GameBoard.IsDraw()
}

func (g *GameVersusBot) BoardToSymbols() [consts.Rows][consts.Columns]string {
	return g.GameBoard.BoardToSymbols()
}

func findXY(pos uint, row, col *uint) error {
	op := "game.findXY"
	if pos < 1 || pos > 9 {
		return fmt.Errorf("%s: pos in [1, 9]", op)
	}
	for i := 0; i < consts.Rows; i++ {
		for j := 0; j < consts.Columns; j++ {
			if uint(i)*consts.Rows+uint(j) == pos-1 {
				*row, *col = uint(i), uint(j)
			}
		}
	}
	return nil
}
