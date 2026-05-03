package board

import (
	"TicTacToe/internal/domain/consts"
	"fmt"
	"math"
)

type Board [consts.Rows][consts.Columns]uint

func (b *Board) CheckWinner(mark uint) bool {
	//TODO implement me
	panic("implement me")
}

func (b *Board) Draw() {
	var char rune
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch b[i][j] {
			case consts.Empty:
				char = ' '
			case consts.Cross:
				char = 'X'
			case consts.Zero:
				char = 'O'
			}
			fmt.Printf(" %c", char)
			if j != consts.Columns-1 {
				fmt.Printf(" |")
			}
		}
		fmt.Printf("\n")
		if i != consts.Rows-1 {
			fmt.Printf("---+---+---\n")
		}
	}
}

func Sum(sl [3]uint) int {
	res := 0
	for i := 0; i < 3; i++ {
		res += int(sl[i])
	}
	return res
}

func (b *Board) FindEmptyCells() [][]uint {
	var notFiledPoss [][]uint
	for i := 0; i < consts.Rows; i++ {
		for j := 0; j < consts.Columns; j++ {
			if b[i][j] == consts.Empty {
				notFiledPoss = append(notFiledPoss, []uint{uint(i), uint(j)})
			}
		}
	}
	return notFiledPoss
}

func (b *Board) MakeMove(row, col uint, player uint) error {
	op := "domain.MakeMove"
	if b[row][col] != 0 {
		return fmt.Errorf("%s: cell already occupied", op)
	}
	if player != 1 && player != 2 {
		return fmt.Errorf("%s: invalid player", op)
	}
	b[row][col] = uint(player)
	return nil
}

func (b *Board) IsFinished() bool {
	return b.CheckWin(consts.Cross) || b.CheckWin(consts.Zero) || b.BoardFiled()
}

func (b *Board) CheckWin(mark uint) bool {
	for i := 0; i < 3; i++ {
		if b[i][0] == mark && b[i][1] == mark && b[i][2] == mark {
			return true
		}
	}
	for i := 0; i < 3; i++ {
		if b[0][i] == mark && b[1][i] == mark && b[2][i] == mark {
			return true
		}
	}
	if b[0][0] == mark && b[1][1] == mark && b[2][2] == mark {
		return true
	}
	if b[0][2] == mark && b[1][1] == mark && b[2][0] == mark {
		return true
	}
	return false
}

func (b *Board) GetAvailableMoves() [][2]int {
	moves := make([][2]int, 0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == 0 {
				moves = append(moves, [2]int{i, j})
			}
		}
	}
	return moves
}

func (b *Board) IsFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

func (b *Board) Copy() Board {
	newBoard := Board{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			newBoard[i][j] = b[i][j]
		}
	}
	return newBoard
}

func (b *Board) IsEmpty() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] != 0 {
				return false
			}
		}
	}
	return true
}

func (b *Board) BoardFiled() bool {
	for i := 0; i < consts.Rows; i++ {
		for j := 0; j < consts.Columns; j++ {
			if b[i][j] == consts.Empty {
				return false
			}
		}
	}
	return true
}

func (b *Board) IsDraw() bool {
	return !b.CheckWin(consts.Zero) && !b.CheckWin(consts.Cross) && b.IsFull()
}

func (b *Board) BoardToSymbols() [consts.Rows][consts.Columns]string {
	var board [consts.Rows][consts.Columns]string
	for i := 0; i < consts.Rows; i++ {
		for j := 0; j < consts.Columns; j++ {
			board[i][j] = NumToSymbol(b[i][j])
		}
	}
	return board
}

func NumToSymbol(num uint) string {
	switch num {
	case consts.Zero:
		return "O"
	case consts.Cross:
		return "X"
	case consts.Empty:
		return " "
	default:
		return ""
	}
}

func SymbolToNum(symbol string) uint {
	switch symbol {
	case "O":
		return consts.Zero
	case "X":
		return consts.Cross
	case " ":
		return consts.Empty
	default:
		return math.MaxInt
	}
}
