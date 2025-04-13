package sessions

import (
	"cinnabon/internal/storage"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const MAX_AGE = time.Hour * 24 * 30 // 30 days

func generateSessionID() string {
	uuid := uuid.New()
	return uuid.String()
}

func New(userID int64) (hash string, err error) {
	h := generateSessionID()
	_, e := storage.Q.CreateSession(context.Background(), storage.CreateSessionParams{
		UserID:    userID,
		ExpiresAt: time.Now().Add(MAX_AGE),
		Hash:      h,
	})
	if e != nil {
		slog.Error("failed to create session", "error", e)
		return "", errors.New("failed to create session")
	}

	return h, nil
}
