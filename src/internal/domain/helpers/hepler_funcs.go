package helpers

import (
	"TicTacToe/internal/domain/consts"
	"math"
)

func ConvertStatus(status uint) string {
	switch status {
	case consts.Active:
		return "Active"
	case consts.Finished:
		return "Finished"
	case consts.Draw:
		return "Draw"
	default:
		return ""
	}
}

func ConvertWinner(winner uint) string {
	switch winner {
	case consts.Player:
		return "Player"
	case consts.Computer:
		return "Computer"
	default:
		return ""
	}
}

func ConvertPlayer(winner uint) string {
	switch winner {
	case consts.Player:
		return "Player"
	case consts.Computer:
		return "Computer"
	case consts.Player1:
		return "Player1"
	case consts.Player2:
		return "Player2"
	default:
		return ""
	}
}

func InverseConvertStatus(status string) uint {
	switch status {
	case "Active":
		return consts.Active
	case "Finished":
		return consts.Finished
	case "Draw":
		return consts.Draw
	default:
		return math.MaxInt
	}
}

func InverseConvertPlayer(winner string) uint {
	switch winner {
	case "Player":
		return consts.Player
	case "Computer":
		return consts.Computer
	case "Player1":
		return consts.Player1
	case "Player2":
		return consts.Player2
	default:
		return math.MaxInt
	}
}
