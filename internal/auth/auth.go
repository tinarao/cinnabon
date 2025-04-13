package auth

import (
	"cinnabon/internal/sessions"
	"cinnabon/internal/storage"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Username  *string `json:"username,omitempty" validate:"omitempty,min=3,max=32,alphanum"`
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=50,alpha"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=50,alpha"`
	Password  string  `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Hash   string `json:"hash"`
	UserID int64  `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

const SESSION_COOKIE_NAME = "session"
const SESSION_COOKIE_PATH = "/"
const SESSION_AUTH_HEADER = "Authorization"

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := storage.Q.GetUserByEmail(context.Background(), req.Email)
	if err == nil {
		sendError(w, "Email already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Convert optional fields to sql.NullString
	var username sql.NullString
	if req.Username != nil {
		username = sql.NullString{
			String: *req.Username,
			Valid:  true,
		}
	}

	var firstName sql.NullString
	if req.FirstName != nil {
		firstName = sql.NullString{
			String: *req.FirstName,
			Valid:  true,
		}
	}

	var lastName sql.NullString
	if req.LastName != nil {
		lastName = sql.NullString{
			String: *req.LastName,
			Valid:  true,
		}
	}

	userID, err := storage.Q.CreateUser(context.Background(), storage.CreateUserParams{
		Email:     req.Email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
	})
	if err != nil {
		sendError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	hash, err := sessions.New(userID)
	if err != nil {
		sendError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    hash,
		Path:     SESSION_COOKIE_PATH,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessions.MAX_AGE.Seconds()),
	})

	sendJSON(w, AuthResponse{
		UserID: userID,
		Hash:   hash,
	}, http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := storage.Q.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	hash, err := sessions.New(user.ID)
	if err != nil {
		sendError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	sendJSON(w, AuthResponse{
		UserID: user.ID,
		Hash:   hash,
	}, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionIdFromHeader(r)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = storage.Q.DeleteSessionByHash(context.Background(), sessionID)
	if err != nil {
		sendError(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func getSessionIdFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get(SESSION_AUTH_HEADER)
	if header == "" {
		return "", errors.New("session ID is required")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid session ID format")
	}

	return parts[1], nil
}
