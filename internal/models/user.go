package models

import (
	"errors"
	"time"
)

type User struct {
	ID           int        `json:"id" db:"id"`
	Nickname     string     `json:"nickname" db:"nickname"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Email        string     `json:"email" db:"email"`
	RealName     *string    `json:"real_name" db:"real_name"`
	BirthDate    *time.Time `json:"birth_date" db:"birth_date"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

func NewUser(nickname, email, passwordHash string, realName *string, birthDate *time.Time) *User {
	now := time.Now()
	return &User{
		Nickname:     nickname,
		Email:        email,
		PasswordHash: passwordHash,
		RealName:     realName,
		BirthDate:    birthDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type RegisterUserInput struct {
	Nickname  string     `json:"nickname"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	RealName  *string    `json:"real_name"`
	BirthDate *time.Time `json:"birth_date"`
}

func (i RegisterUserInput) ToUser(hashedPassword string) *User {
	return NewUser(i.Nickname, i.Email, hashedPassword, i.RealName, i.BirthDate)
}

func (i RegisterUserInput) Validate() error {
	if i.Nickname == "" {
		return errors.New("Имя не может быть пустым")
	}
	if i.Password == "" {
		return errors.New("Пароль не может быть пустым")
	}
	if i.Email == "" {
		return errors.New("Почта не может быть пустой")
	}
	return nil
}

type LoginUserInput struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (i LoginUserInput) Validate() error {
	if i.Identifier == "" {
		return errors.New("Идентификатор не может быть пустым")
	}
	if i.Password == "" {
		return errors.New("Пароль не может быть пуст")
	}
	return nil
}
