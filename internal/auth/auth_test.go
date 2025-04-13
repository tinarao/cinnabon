package auth

import (
	"bytes"
	"cinnabon/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestAuth(t *testing.T) {
	storage.SetupTestDB()
}

func teardownTestAuth(t *testing.T) {
	storage.TeardownTestDB()
}

func TestRegister(t *testing.T) {
	setupTestAuth(t)
	defer teardownTestAuth(t)

	username := "testuser"
	firstName := "Test"
	lastName := "User"

	reqBody := RegisterRequest{
		Email:     "test@example.com",
		Username:  &username,
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  "securepassword123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.UserID)

	// duplicate email
	w = httptest.NewRecorder()
	Register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// invalid email
	reqBody.Email = "invalid-email"
	body, _ = json.Marshal(reqBody)
	req = httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	Register(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin(t *testing.T) {
	setupTestAuth(t)
	defer teardownTestAuth(t)

	// First register a user
	username := "loginuser"
	reqBody := RegisterRequest{
		Email:    "login@example.com",
		Username: &username,
		Password: "securepassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, req)

	// Test valid login
	loginReq := LoginRequest{
		Email:    "login@example.com",
		Password: "securepassword123",
	}
	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	Login(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Hash)

	// Test invalid password
	loginReq.Password = "wrongpassword"
	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	Login(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test non-existent user
	loginReq.Email = "nonexistent@example.com"
	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	Login(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogout(t *testing.T) {
	setupTestAuth(t)
	defer teardownTestAuth(t)

	// First register and login to get a session
	reqBody := RegisterRequest{
		Email:    "logout@example.com",
		Password: "securepassword123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	Register(w, req)

	var resp AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Test valid logout
	req = httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Hash)
	w = httptest.NewRecorder()

	Logout(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	// invalid session
	req = httptest.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer invalid-session")
	w = httptest.NewRecorder()

	Logout(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
