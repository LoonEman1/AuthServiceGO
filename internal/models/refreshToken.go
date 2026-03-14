package models

import (
	"errors"
	"time"
)

type RefreshToken struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

func (i RefreshInput) Validate() error {

	if i.RefreshToken == "" {
		return errors.New("Токен не может быть пустым")
	}

	if len(i.RefreshToken) < 12 || len(i.RefreshToken) > 200 {
		return errors.New("Невалидный токен")
	}
	return nil
}
