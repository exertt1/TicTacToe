package game

import (
	"TicTacToe/internal/domain/consts"
	"TicTacToe/internal/pkg/errors"
	"math"
	"time"
)

type Move struct {
	Row   int
	Col   int
	Score int
}

func (g *GameVersusBot) ComputerMove() error {
	// Проверяем, можно ли делать ход
	if g.Status != consts.Active {
		return gameErrors.ErrGameFinished
	}

	if g.CurrentTurn != consts.Computer {
		return gameErrors.ErrNotComputerTurn
	}

	bestMove := g.minimax(g.ComputerMark)

	if bestMove.Row == -1 && bestMove.Col == -1 {
		return gameErrors.ErrNoMovesLeft
	}

	g.GameBoard[bestMove.Row][bestMove.Col] = g.ComputerMark
	g.LastComputerMove = Move{Row: bestMove.Row, Col: bestMove.Col}
	g.UpdatedAt = time.Now()

	if g.checkWin(g.ComputerMark) {
		g.Winner = consts.Computer
		g.Status = consts.Finished
		return nil
	}

	if g.isDraw() {
		g.Status = consts.Draw
		return nil
	}

	g.CurrentTurn = consts.Player

	return nil
}

func (g *GameVersusBot) minimax(currentPlayer uint) Move {
	bestMove := Move{Row: -1, Col: -1, Score: 0}

	if g.checkWin(g.ComputerMark) {
		return Move{Row: -1, Col: -1, Score: 10}
	} else if g.checkWin(g.PlayerMark) {
		return Move{Row: -1, Col: -1, Score: -10}
	} else if g.isDraw() {
		return Move{Row: -1, Col: -1, Score: 0}
	}

	if currentPlayer == g.ComputerMark {
		bestMove.Score = -math.MaxInt32

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if g.GameBoard[i][j] == consts.Empty {
					g.GameBoard[i][j] = g.ComputerMark

					move := g.minimax(g.PlayerMark)

					g.GameBoard[i][j] = consts.Empty

					if move.Score > bestMove.Score {
						bestMove.Score = move.Score
						bestMove.Row = i
						bestMove.Col = j
					}
				}
			}
		}
	} else {
		// Если ход игрока
		bestMove.Score = math.MaxInt32

		// Перебираем все пустые клетки
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if g.GameBoard[i][j] == consts.Empty {
					g.GameBoard[i][j] = g.PlayerMark

					move := g.minimax(g.ComputerMark)

					g.GameBoard[i][j] = consts.Empty

					if move.Score < bestMove.Score {
						bestMove.Score = move.Score
						bestMove.Row = i
						bestMove.Col = j
					}
				}
			}
		}
	}

	return bestMove
}

// checkWin - проверяет победу для заданной метки
func (g *GameVersusBot) checkWin(mark uint) bool {
	// Проверка строк
	return g.GameBoard.CheckWin(mark)
}

// isDraw - проверяет ничью
func (g *GameVersusBot) isDraw() bool {
	return g.GameBoard.IsDraw()
}
