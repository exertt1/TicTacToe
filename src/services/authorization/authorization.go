package authorization

import (
	"TicTacToe/internal/domain/repository"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	Message string `json:"msg"`
}

type UserService struct {
	storage repository.GameRepository
}

func NewUserService(storage repository.GameRepository) *UserService {
	if storage == nil {
		panic("storage is nil in NewUserService")
	}
	return &UserService{storage: storage}
}

func (s *UserService) CheckLogin(ctx context.Context, login string) bool {
	query := `
	SELECT COUNT(*) FROM players
	WHERE login = ?
`
	var count int

	err := s.storage.QueryRow(ctx, query, login).Scan(&count)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if count != 0 {
		return false
	}
	return true

}

func (s *UserService) SaveUser(ctx context.Context, user User) error {
	query := `
	INSERT INTO players (
		id, login, password
	)
	VALUES ($1, $2, $3)
`
	_, err := s.storage.Exec(ctx, query, user.ID, user.Login, user.Password)

	return err
}

func (s *UserService) GetUser(ctx context.Context, login string) (*User, error) {
	op := "userService.GetUser"

	var user User

	query := `
	SELECT * FROM players
	WHERE login = $1
`
	rows, err := s.storage.Query(ctx, query, login)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	user, err = pgx.CollectOneRow(
		rows,
		pgx.RowToStructByName[User],
	)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	log.Printf("userID=%s", user.ID.String())
	return &user, nil
}

func (s *UserService) Registration(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	exists, err := s.storage.CheckLogin(r.Context(), req.Login)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s:%s", "s.storage.CheckLogin", err.Error()), http.StatusInternalServerError)
		return
	}

	response := RegistrationResponse{}

	if exists {
		response.Message = "this login already exists"
		json.NewEncoder(w).Encode(&response)
		http.Error(w, "this login already exists", http.StatusBadRequest)
		return
	}
	playerID, err := uuid.NewV7()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{
		ID:       playerID,
		Login:    req.Login,
		Password: req.Password,
	}

	err = s.SaveUser(r.Context(), user)

	if err != nil {
		http.Error(w, fmt.Sprintf("%s:%s", "s.SaveUser", err.Error()), http.StatusBadRequest)
		return
	}
	response.Message = fmt.Sprintf("Your ID: %s", playerID)
	json.NewEncoder(w).Encode(&response)

}

type AuthorizationResponse struct {
	ID      uuid.UUID `json:"id"`
	Message string    `json:"msg"`
}

func (s *UserService) Authorization(w http.ResponseWriter, r *http.Request) {
	signUp, err := BasicToUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := s.GetUser(r.Context(), signUp.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := AuthorizationResponse{}
	if user.Password != signUp.Password {
		resp.Message = "The password wrong"
		json.NewEncoder(w).Encode(&resp)
		http.Error(w, "The password wrong", http.StatusUnauthorized)
		return
	}
	resp.ID = user.ID
	resp.Message = "Authorization completed success!"

	json.NewEncoder(w).Encode(&resp)

}

func BasicToUser(w http.ResponseWriter, r *http.Request) (*SignUpRequest, error) {
	auth := r.Header.Get("Authorization")

	decode, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
	if err != nil {
		return nil, err
	}
	logPass := strings.Split(string(decode), ":")
	login, password := logPass[0], logPass[1]
	return &SignUpRequest{
		Login:    login,
		Password: password,
	}, nil
}

func (s *UserService) Authenticate(ctx context.Context, login, password string) (bool, error) {
	user, err := s.GetUser(ctx, login)
	if err != nil {
		return false, err
	}
	if user.Password != password {
		return false, fmt.Errorf("Incorrect password")
	}
	return true, nil

}
