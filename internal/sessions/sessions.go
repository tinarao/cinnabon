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

func GetSessionByHash(hash string) (userID int64, err error) {
	session, err := storage.Q.GetSessionByHash(context.Background(), hash)
	if err != nil {
		return 0, errors.New("failed to get session")
	}

	return session.UserID, nil
}

func GetUserBySessionHash(hash string) (userID int64, err error) {
	userID, err = GetSessionByHash(hash)
	if err != nil {
		return 0, errors.New("failed to get user by session hash")
	}

	user, err := storage.Q.GetUserByID(context.Background(), userID)
	if err != nil {
		return 0, errors.New("failed to get user by session hash")
	}

	return user.ID, nil
}
