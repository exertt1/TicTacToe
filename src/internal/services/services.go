package services

import (
	"myproject/internal/models"
)

func StartGame() {
	game := models.InitGameField()
	game.Game()
}
