package server

import (
	"cinnabon/internal/auth"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Port string
}

func New(port string) *Server {
	return &Server{
		Port: port,
	}
}

func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", s.healthcheck).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/register", auth.Register).Methods("POST")
	apiRouter.HandleFunc("/login", auth.Login).Methods("POST")
	apiRouter.HandleFunc("/logout", auth.Logout).Methods("POST")

	if err := http.ListenAndServe(s.Port, router); err != nil {
		slog.Error("failed to start server", "error", err.Error())
	}
}

func (s *Server) healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
