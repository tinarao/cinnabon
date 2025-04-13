package main

import (
	"cinnabon/config"
	"cinnabon/internal/server"
	"cinnabon/internal/storage"
	"log/slog"
	"os"
)

func main() {
	c := config.New("Secret", ":8080")
	c.Load()

	if err := storage.Init(); err != nil {
		slog.Error("failed to initialize storage", "error", err.Error())
		os.Exit(1)
	}

	defer storage.Conn.Close()

	server := server.New(c.Port)
	server.Start()
}
