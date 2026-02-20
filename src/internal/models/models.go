package models

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/google/uuid"
)

const (
	Empty = iota
	Cross
	Zero
)

const (
	Rows = 3
	Columns
)

const (
	Negative = iota - 1
	Neutral
	Positive
)

type Step struct {
	PosY             int
	PosX             int
	SuccessOfTheMove int
}

type GameField struct {
	Field    [Rows][Columns]uint
	IsFilled bool
	ID       uuid.UUID
}

func InitGameField() *GameField {
	id, _ := uuid.NewV7()
	return &GameField{
		IsFilled: false,
		ID:       id,
	}
}

func (gf *GameField) Draw() {
	var char rune
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			switch gf.Field[i][j] {
			case Empty:
				char = ' '
			case Cross:
				char = 'X'
			case Zero:
				char = 'O'
			}
			fmt.Printf(" %c", char)
			if j != Columns-1 {
				fmt.Printf(" |")
			}
		}
		fmt.Printf("\n")
		if i != Rows-1 {
			fmt.Printf("---+---+---\n")
		}
	}
}

func (gf *GameField) IsFinished() bool {
	isFinished := true
	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns; j++ {
			if gf.Field[i][j] == 0 {
				isFinished = false
			}
		}
	}
	var sl []uint
	for i := 0; i < Rows; i++ {
		sl = append(sl, gf.Field[i][i])
	}
	if (slices.Compare(sl, []uint{1, 1, 1}) == 0) || (slices.Compare(sl, []uint{2, 2, 2}) == 0) {
		return true
	}
	sl = sl[:0]
	for i := Rows - 1; i >= 0; i-- {
		sl = append(sl, gf.Field[i][i])
	}
	if (slices.Compare(sl, []uint{1, 1, 1}) == 0) || (slices.Compare(sl, []uint{2, 2, 2}) == 0) {
		return true
	}
	sl = sl[:0]

	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns; j++ {
			sl = append(sl, gf.Field[i][j])
		}
		if (slices.Compare(sl, []uint{1, 1, 1}) == 0) || (slices.Compare(sl, []uint{2, 2, 2}) == 0) {
			return true
		}
		sl = sl[:0]
	}
	var slT []uint
	for i := 0; i < Columns; i++ {
		for j := 0; j < Rows; j++ {
			slT = append(slT, gf.Field[j][i])
		}
		if (slices.Compare(slT, []uint{1, 1, 1}) == 0) || (slices.Compare(slT, []uint{2, 2, 2}) == 0) {
			return true
		}
		slT = slT[:0]
	}
	return isFinished
}

func (gf *GameField) Game() {
	var pos int
	for !gf.IsFinished() {
		gf.Draw()
		fmt.Println("Ваш ход: ")
		_, _ = fmt.Scan(&pos)
		if pos < 1 || pos > 9 {
			fmt.Println("Вы указали неверную позицию")
			continue
		}
		//[i][j] = [i*Rows + j] = pos - 1
		var (
			y int
			x int
		)
		findedXY := false
		for i := 0; i < Rows; i++ {
			for j := 0; j < Columns; j++ {
				if i*Rows+j == pos-1 {
					y, x = i, j
					findedXY = true
					break
				}
			}
			if findedXY == true {
				break
			}
		}
		if gf.Field[y][x] == Empty {
			gf.Field[y][x] = Cross
		} else {
			fmt.Println("Эта позиция уже занята")
			continue
		}
		notFilledPoss := gf.FindEmptyCells()
		var ran int
		if len(notFilledPoss) > 1 {
			ran = rand.Intn(len(notFilledPoss))
		} else {
			ran = 0
		}
		randY, randX := notFilledPoss[ran][0], notFilledPoss[ran][1]
		gf.Field[randY][randX] = Zero
	}
	gf.Draw()
	fmt.Println("игра закончена")
}

func Sum(sl [3]uint) int {
	res := 0
	for i := 0; i < 3; i++ {
		res += int(sl[i])
	}
	return res
}

func (gf *GameField) MiniMax() {
	notFilledPoss := gf.FindEmptyCells()
	
}

func (gf *GameField) FindEmptyCells() [][]uint {
	var notFiledPoss [][]uint
	for i := 0; i < Rows; i++ {
		for j := 0; j < Columns; j++ {
			if gf.Field[i][j] == Empty {
				notFiledPoss = append(notFiledPoss, []uint{uint(i), uint(j)})
			}
		}
	}
	return notFiledPoss
}
