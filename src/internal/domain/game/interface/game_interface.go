package _interface

import "TicTacToe/internal/domain/consts"

type Game interface {
	MakeMove(row, col uint, player uint) error

	FindEmptyCells() [][]uint

	IsFinished() bool

	CheckWinner(mark uint) bool

	IsDraw() bool

	BoardToSymbols() [consts.Rows][consts.Columns]string
}
