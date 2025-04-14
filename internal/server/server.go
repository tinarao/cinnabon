package server

import (
	"cinnabon/internal/auth"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	Port   string
	router *mux.Router
}

func New(port string) *Server {
	return &Server{
		Port:   port,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() {
	s.router.Use(s.logger)
	s.router.HandleFunc("/healthcheck", s.healthcheck).Methods("GET")

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/register", auth.Register).Methods("POST")
	apiRouter.HandleFunc("/login", auth.Login).Methods("POST")
	apiRouter.HandleFunc("/logout", auth.Logout).Methods("POST")

	done := make(chan struct{})
	go func() {
		s.listen()
		close(done)
	}()

	slog.Info("Server started", "PORT", s.Port)
	<-done
}

func (s *Server) listen() {
	if err := http.ListenAndServe(s.Port, s.router); err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

func (s *Server) healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *Server) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "method", r.Method, "url", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
