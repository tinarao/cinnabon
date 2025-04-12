package sessions

import "time"

type Session struct {
	ID        uint      `json:"id"`
	ForeignID string    `json:"foreign_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func New() *Session {
	return &Session{}
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
