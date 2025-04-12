package main

import (
	"cinnabon/config"
	"cinnabon/internal/storage"
	"context"
	"fmt"
	"log/slog"
)

func main() {
	c := config.New("Secret", ":8080")
	c.Load()

	storage.Init()
	defer storage.Conn.Close()

	u, err := storage.Q.GetUserByID(context.Background(), 1)
	if err != nil {
		slog.Error("failed to get user", "error", err.Error())
	}

	fmt.Println(u)
}
