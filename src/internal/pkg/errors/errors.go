package gameErrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

var (
	// Игровые ошибки (400)
	ErrGameFinished    = New("GAME_FINISHED", "Game is already finished", http.StatusBadRequest)
	ErrNotPlayerTurn   = New("NOT_PLAYER_TURN", "It's not player's turn", http.StatusBadRequest)
	ErrNotComputerTurn = New("NOT_COMPUTER_TURN", "It's not computer's turn", http.StatusBadRequest)
	ErrInvalidMove     = New("INVALID_MOVE", "Invalid move coordinates", http.StatusBadRequest)
	ErrCellOccupied    = New("CELL_OCCUPIED", "Cell is already occupied", http.StatusBadRequest)
	ErrNoMovesLeft     = New("NO_MOVES_LEFT", "No moves left", http.StatusBadRequest)

	// Ошибки репозитория (404)
	ErrGameNotFound = New("GAME_NOT_FOUND", "Game not found", http.StatusNotFound)

	// Внутренние ошибки (500)
	ErrInternalServer = New("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrFailedToSave   = New("FAILED_TO_SAVE", "Failed to save game", http.StatusInternalServerError)
	ErrFailedToUpdate = New("FAILED_TO_UPDATE", "Failed to update game", http.StatusInternalServerError)
)

// Ошибки валидации
var (
	ErrInvalidGameID   = New("INVALID_GAME_ID", "Invalid game ID", http.StatusBadRequest)
	ErrInvalidPlayerID = New("INVALID_PLAYER_ID", "Invalid player ID", http.StatusBadRequest)
	ErrInvalidMark     = New("INVALID_MARK", "Invalid mark, must be X or O", http.StatusBadRequest)
)
