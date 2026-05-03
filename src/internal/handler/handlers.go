package handler

import (
	"TicTacToe/internal/domain/board"
	"TicTacToe/internal/domain/consts"
	"TicTacToe/internal/domain/game"
	"TicTacToe/internal/domain/helpers"
	"TicTacToe/internal/domain/repository"
	"TicTacToe/services/authorization"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

type GameHandler struct {
	gameRepo    repository.GameRepository
	userService *authorization.UserService
}

func NewGameHandler(gameRepo repository.GameRepository, userService *authorization.UserService) *GameHandler {
	return &GameHandler{
		gameRepo:    gameRepo,
		userService: userService,
	}
}

type GameRequest struct {
	Mode string `json:"mode"`
}

type ComputerGameResponse struct {
	GameID       uuid.UUID `json:"id"`
	PlayerID     uuid.UUID `json:"player_id"`
	PlayerMark   string    `json:"player_mark"`
	ComputerMark string    `json:"computer_mark"`
	FirstTurn    string    `json:"first_turn"`
}

type PlayerGameResponse struct {
	GameID      uuid.UUID `json:"id"`
	Player1ID   uuid.UUID `json:"player1_id"`
	Player1Mark string    `json:"player1_mark"`
	Player2ID   uuid.UUID `json:"player2_id,omitempty"`
	Player2Mark string    `json:"player2_mark"`
	FirstTurn   string    `json:"first_turn"`
}

func (h *GameHandler) createGameVersusBot(ctx context.Context, playerID uuid.UUID) (*game.GameVersusBot, error) {
	op := "gameHandler.createGameVersusBot"

	newGame := game.InitGameVersusBot(playerID)

	playersGoFirst := rand.Intn(2) == 0
	playerMark := consts.Cross
	if !playersGoFirst {
		playerMark = consts.Zero
	}
	computerMark := consts.Zero
	if playerMark == consts.Zero {
		computerMark = consts.Cross
	}
	newGame.PlayerMark = playerMark
	newGame.ComputerMark = computerMark
	if playerMark == consts.Cross {
		newGame.CurrentTurn = consts.Player
	} else {
		newGame.CurrentTurn = consts.Computer
	}

	err := h.gameRepo.SaveGame(ctx, newGame)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !playersGoFirst {
		err := newGame.ComputerMove()
		if err != nil {
			return nil, err
		}
		_ = h.gameRepo.UpdateGame(ctx, newGame)
	}
	return newGame, nil
}

func (h *GameHandler) createGameVersusPlayer(ctx context.Context, playerID uuid.UUID) (*game.GameVersusPlayer, error) {
	op := "gameHandler.createGameVersusPlayer"
	log.Println(playerID)
	newGame := game.InitGameVersusPlayer(playerID)
	log.Println(2)
	log.Println(newGame.FirstPlayerID)
	var player1Mark, player2Mark uint
	playersGoFirst := rand.Intn(2) == 0
	if playersGoFirst {
		player1Mark = consts.Cross
		player2Mark = consts.Zero
	} else {
		player1Mark = consts.Zero
		player2Mark = consts.Cross
	}
	newGame.Player1Mark = player1Mark
	newGame.Player2Mark = player2Mark
	if player1Mark == consts.Cross {
		newGame.CurrentTurn = consts.Player1
	} else {
		newGame.CurrentTurn = consts.Player2
	}

	err := h.gameRepo.SaveGame(ctx, newGame)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return newGame, nil
}

func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req GameRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "incorrect request fro creating game", http.StatusBadRequest)
	}
	signUp, err := authorization.BasicToUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ok, err := h.userService.Authenticate(r.Context(), signUp.Login, signUp.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !ok {
		http.Error(w, "Incorrect password or login", http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUser(r.Context(), signUp.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Mode {
	case consts.BotMode:
		newGame, err := h.createGameVersusBot(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var resp ComputerGameResponse
		resp.GameID = newGame.ID
		resp.PlayerID = newGame.PlayerID
		resp.FirstTurn = helpers.ConvertPlayer(newGame.CurrentTurn)
		resp.PlayerMark = board.NumToSymbol(newGame.PlayerMark)
		resp.ComputerMark = board.NumToSymbol(newGame.ComputerMark)
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	case consts.PlayerMode:
		newGame, err := h.createGameVersusPlayer(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(newGame.ID)
		var resp PlayerGameResponse
		resp.GameID = newGame.ID
		resp.Player1ID = newGame.FirstPlayerID
		resp.FirstTurn = helpers.ConvertPlayer(newGame.CurrentTurn)
		resp.Player1Mark = board.NumToSymbol(newGame.Player1Mark)
		resp.Player2Mark = board.NumToSymbol(newGame.Player2Mark)

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type JoinGameResponse struct {
	GameID      uuid.UUID `json:"id"`
	Player1ID   uuid.UUID `json:"player1_id"`
	Player1Mark string    `json:"player1_mark"`
	Player2ID   uuid.UUID `json:"player2_id"`
	Player2Mark string    `json:"player2_mark"`
	CurrentTurn string    `json:"current_turn"`
}

func (h *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.Context().Value("id")

	signUp, err := authorization.BasicToUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	ok, err := h.userService.Authenticate(r.Context(), signUp.Login, signUp.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !ok {
		http.Error(w, "Incorrect login or password", http.StatusUnauthorized)
	}
	user, err := h.userService.GetUser(r.Context(), signUp.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	g, err := h.gameRepo.GetGameByID(r.Context(), gameID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var resp JoinGameResponse
	switch gm := g.(type) {
	case *game.GameVersusPlayer:
		gm.SecondPlayerID = user.ID
		gm.Status = consts.Active
		resp.GameID = gm.ID
		resp.Player1ID = gm.FirstPlayerID
		resp.Player2ID = gm.SecondPlayerID
		resp.Player1Mark = board.NumToSymbol(gm.Player1Mark)
		resp.Player2Mark = board.NumToSymbol(gm.Player2Mark)
		resp.CurrentTurn = helpers.ConvertPlayer(gm.CurrentTurn)
		err := h.gameRepo.UpdateGame(r.Context(), gm)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	case *game.GameVersusBot:
		http.Error(w, "you cant join to game versus bot", http.StatusBadRequest)
		return
	}
}

type MakeMoveRequest struct {
	Row uint `json:"row"`
	Col uint `json:"col"`
}

type MakeMoveResponse struct {
	Board        [consts.Rows][consts.Columns]string `json:"board"`
	CurrentTurn  string                              `json:"current_turn"`
	Player1Mark  string                              `json:"player1_mark,omitempty"`
	Player2Mark  string                              `json:"player2_mark,omitempty"`
	PlayerMark   string                              `json:"player_mark,omitempty"`
	ComputerMark string                              `json:"computer_mark,omitempty"`
	Status       string                              `json:"status"`
	Winner       string                              `json:"winner"`
	ComputerMove *MoveInfo                           `json:"computer_move,omitempty"`
	Message      string                              `json:"message"`
}

type MoveInfo struct {
	Row uint
	Col uint
}

func (h *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	op := "gameHandler.MakeMove"

	gameID := r.Context().Value("id")

	signUp, err := authorization.BasicToUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.userService.Authenticate(r.Context(), signUp.Login, signUp.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var req MakeMoveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	g, err := h.gameRepo.GetGameByID(r.Context(), gameID.(string))
	if err != nil {
		http.Error(w, "Game not founded", http.StatusNotFound)
		return
	}

	if req.Row < 1 || req.Row > 3 || req.Col < 1 || req.Col > 3 {
		http.Error(w, "invalid value for position", http.StatusBadRequest)
	}
	row, col := req.Row-1, req.Col-1

	switch gm := g.(type) {
	case *game.GameVersusBot:
		if gm.Status != consts.Active {
			http.Error(w, "game is finished", http.StatusBadRequest)
			return
		}

		if gm.CurrentTurn != consts.Player {
			http.Error(w, "computer turn", http.StatusBadRequest)
		}

		if err := gm.PlayerMove(row, col); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := MakeMoveResponse{
			Board:        gm.GameBoard.BoardToSymbols(),
			CurrentTurn:  helpers.ConvertPlayer(gm.CurrentTurn),
			PlayerMark:   board.NumToSymbol(gm.PlayerMark),
			ComputerMark: board.NumToSymbol(gm.ComputerMark),
			Status:       helpers.ConvertStatus(gm.Status),
			Message:      gm.Message,
		}

		if gm.Status != consts.Active {
			if gm.Winner == consts.Player {
				gm.Message = "You won"
			} else if gm.Winner == consts.Nobody {
				gm.Message = "Draw"
			}
			resp.Winner = helpers.ConvertWinner(gm.Winner)
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(&resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		if err := gm.ComputerMove(); err != nil {
			http.Error(w, fmt.Sprintf("%s: %v", op, err), http.StatusBadRequest)
			return
		}
		if gm.Status != consts.Active {
			if gm.Winner == consts.Computer {
				gm.Message = "Computer won"
			} else if gm.Winner == consts.Nobody {
				gm.Message = "Draw"
			}
			resp.Winner = helpers.ConvertWinner(gm.Winner)
		} else {
			resp.Message = "Your turn"
		}
		resp.Board = gm.GameBoard.BoardToSymbols()
		resp.CurrentTurn = helpers.ConvertPlayer(gm.CurrentTurn)
		resp.Status = helpers.ConvertStatus(gm.Status)
		resp.Message = gm.Message

		err := h.gameRepo.UpdateGame(r.Context(), gm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	case *game.GameVersusPlayer:
		if gm.Status != consts.Active {
			http.Error(w, "game is finished", http.StatusBadRequest)
			return
		}

		user, err := h.userService.GetUser(r.Context(), signUp.Login)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if user.ID != gm.FirstPlayerID && user.ID != gm.SecondPlayerID {
			http.Error(w, "It's not your game", http.StatusUnauthorized)
			return
		}

		if gm.CurrentTurn == consts.Player1 && user.ID != gm.FirstPlayerID {
			http.Error(w, "Player1 turn", http.StatusBadRequest)
		}

		if gm.CurrentTurn == consts.Player2 && user.ID != gm.SecondPlayerID {
			http.Error(w, "Player1 turn", http.StatusBadRequest)
		}

		if gm.CurrentTurn == consts.Player1 {
			err := gm.MakeMove(row, col, gm.Player1Mark)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			gm.CurrentTurn = consts.Player2
		} else {
			err := gm.MakeMove(row, col, gm.Player2Mark)
			log.Println(1312312312312)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			gm.CurrentTurn = consts.Player1
		}

		resp := MakeMoveResponse{
			Board:       gm.GameBoard.BoardToSymbols(),
			CurrentTurn: helpers.ConvertPlayer(gm.CurrentTurn),
			Player1Mark: board.NumToSymbol(gm.Player1Mark),
			Player2Mark: board.NumToSymbol(gm.Player2Mark),
			Status:      helpers.ConvertStatus(gm.Status),
			Message:     gm.Message,
			Winner:      helpers.ConvertPlayer(gm.Winner),
		}

		if gm.Status != consts.Active {
			if gm.Winner == consts.Player1 {
				gm.Message = "Player1 won"
			} else if gm.Winner == consts.Nobody {
				gm.Message = "Draw"
			} else {
				gm.Message = "Player2 won"
			}
			resp.Winner = helpers.ConvertWinner(gm.Winner)
			err := json.NewEncoder(w).Encode(&resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}

		err = h.gameRepo.UpdateGame(r.Context(), gm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

type GetComputerGameResponse struct {
	GameID       uuid.UUID                           `json:"id"`
	Board        [consts.Rows][consts.Columns]string `json:"board"`
	PlayerMark   string                              `json:"player_mark"`
	ComputerMark string                              `json:"computer_mark"`
	CurrentTurn  string                              `json:"current_turn"`
}

type GetPlayerGameResponse struct {
	GameID      uuid.UUID                           `json:"id"`
	Board       [consts.Rows][consts.Columns]string `json:"board"`
	Player1Mark string                              `json:"player_mark"`
	Player2Mark string                              `json:"computer_mark"`
	CurrentTurn string                              `json:"current_turn"`
}

func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	op := "gameHandler.GetGame"

	gameID := r.Context().Value("id")
	if gameID == nil {
		http.Error(w, "Game ID not found in context", http.StatusBadRequest)
		return
	}

	g, err := h.gameRepo.GetGameByID(r.Context(), gameID.(string))
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %w", op, err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch gm := g.(type) {
	case *game.GameVersusBot:
		resp := GetComputerGameResponse{
			GameID:       gm.ID,
			Board:        g.BoardToSymbols(),
			PlayerMark:   board.NumToSymbol(gm.PlayerMark),
			ComputerMark: board.NumToSymbol(gm.ComputerMark),
			CurrentTurn:  helpers.ConvertPlayer(gm.CurrentTurn),
		}

		err := json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	case *game.GameVersusPlayer:
		resp := GetPlayerGameResponse{
			GameID:      gm.ID,
			Board:       g.BoardToSymbols(),
			Player1Mark: board.NumToSymbol(gm.Player1Mark),
			Player2Mark: board.NumToSymbol(gm.Player2Mark),
			CurrentTurn: helpers.ConvertPlayer(gm.CurrentTurn),
		}
		err := json.NewEncoder(w).Encode(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
