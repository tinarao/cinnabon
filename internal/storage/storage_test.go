package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserCRUD(t *testing.T) {
	SetupTestDB()
	defer TeardownTestDB()

	ctx := context.Background()

	// Test CreateUser
	userID, err := Q.CreateUser(ctx, CreateUserParams{
		Email:     "test@example.com",
		Username:  sql.NullString{String: "testuser", Valid: true},
		FirstName: sql.NullString{String: "Test", Valid: true},
		LastName:  sql.NullString{String: "User", Valid: true},
		Password:  "hashedpassword",
	})
	require.NoError(t, err)
	assert.NotZero(t, userID)

	// Test GetUserByID
	user, err := Q.GetUserByID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username.String)
	assert.Equal(t, "Test", user.FirstName.String)
	assert.Equal(t, "User", user.LastName.String)

	// Test GetUserByEmail
	userByEmail, err := Q.GetUserByEmail(ctx, "test@example.com")
	require.NoError(t, err)
	assert.Equal(t, userID, userByEmail.ID)
}

func TestSessionCRUD(t *testing.T) {
	SetupTestDB()
	defer TeardownTestDB()

	ctx := context.Background()

	// First create a user
	userID, err := Q.CreateUser(ctx, CreateUserParams{
		Email:    "session@example.com",
		Password: "hashedpassword",
	})
	require.NoError(t, err)

	// Test CreateSession
	hash := "test-session-hash"
	expiresAt := time.Now().Add(24 * time.Hour)
	sessionID, err := Q.CreateSession(ctx, CreateSessionParams{
		Hash:      hash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	require.NoError(t, err)
	assert.NotZero(t, sessionID)

	// Test GetSessionByHash
	session, err := Q.GetSessionByHash(ctx, hash)
	require.NoError(t, err)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, hash, session.Hash)

	// Test DeleteSessionByHash
	err = Q.DeleteSessionByHash(ctx, hash)
	require.NoError(t, err)

	// Verify session is deleted
	_, err = Q.GetSessionByHash(ctx, hash)
	require.Error(t, err)
}
