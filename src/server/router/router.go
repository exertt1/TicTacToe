package router

import (
	"TicTacToe/internal/handler"
	"TicTacToe/services/authorization"
	"context"
	"log"
	"net/http"
	"strings"
)

type Router struct {
	routes      map[string]map[string]http.HandlerFunc
	gameHandler *handler.GameHandler
	userService *authorization.UserService
}

func NewRouter(gameHandler *handler.GameHandler, userService *authorization.UserService) *Router {
	r := &Router{
		routes:      make(map[string]map[string]http.HandlerFunc),
		gameHandler: gameHandler,
		userService: userService,
	}
	r.registerRoutes()
	return r
}

func (r *Router) registerRoutes() {
	r.HandleFunc("POST", "tictactoe/games", r.gameHandler.CreateGame)
	r.HandleFunc("GET", "tictactoe/games/{id}", r.gameHandler.GetGame)
	r.HandleFunc("POST", "tictactoe/games/{id}", r.gameHandler.MakeMove)
	r.HandleFunc("POST", "tictactoe/reg", r.userService.Registration)
	r.HandleFunc("GET", "tictactoe/auth", r.userService.Authorization)
	r.HandleFunc("POST", "tictactoe/join/{id}", r.gameHandler.JoinGame)
}

func (r *Router) HandleFunc(method, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) matchPath(pattern, path string) bool {
	patternPart := strings.Split(strings.Trim(pattern, "/"), "/")
	pathPart := strings.Split(strings.Trim(path, "/"), "/")

	if len(pathPart) != len(patternPart) {
		return false
	}

	for i := 0; i < len(pathPart); i++ {
		if strings.HasPrefix(patternPart[i], "{") && strings.HasSuffix(patternPart[i], "}") {
			continue
		}
		if pathPart[i] != patternPart[i] {
			return false
		}
	}
	return true
}

func (r *Router) addPathParams(req *http.Request, pattern, path string) {
	patternPart := strings.Split(strings.Trim(pattern, "/"), "/")
	pathPart := strings.Split(strings.Trim(path, "/"), "/")

	ctx := req.Context()
	for i := 0; i < len(patternPart); i++ {
		if strings.HasPrefix(patternPart[i], "{") && strings.HasSuffix(patternPart[i], "}") {
			paramName := strings.Trim(patternPart[i], "{}")
			paramValue := pathPart[i]

			ctx = context.WithValue(ctx, paramName, paramValue)
			log.Printf("Added param: %s = %s\n", paramName, paramValue)
		}
	}

	*req = *req.WithContext(ctx)
}

type contextKey string

func GetPathParam(r *http.Request, key string) string {
	if val := r.Context().Value(contextKey(key)); val != nil {
		return val.(string)
	}
	return ""
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := req.URL.Path
	log.Println(path)
	if handlers, ok := r.routes[req.Method]; ok {
		log.Println(handlers)
		for routePath, handlerFunc := range handlers {
			log.Println("Проверяем путь запроса с шаблонным")
			if r.matchPath(routePath, path) {
				log.Println("проверка прошла")
				r.addPathParams(req, routePath, path)
				handlerFunc(w, req)
				return
			}
		}
	}

	http.NotFound(w, req)
}
