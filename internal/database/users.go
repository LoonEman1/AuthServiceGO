package database

import (
	"AuthService/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) Register(input models.User) (*models.User, error) {
	var user models.User

	query := `
	INSERT INTO users (nickname, password_hash, email, real_name, birth_date, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7)
	returning id, nickname, password_hash, email, real_name, birth_date, created_at, updated_at;
	`

	err := u.db.QueryRowx(
		query,
		input.Nickname,
		input.PasswordHash,
		input.Email,
		input.RealName,
		input.BirthDate,
		input.CreatedAt,
		input.UpdatedAt,
	).StructScan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) GetUserByIdentifier(identifier string) (*models.User, error) {
	var user models.User

	query := `
	SELECT id, nickname, email, real_name, birth_date, created_at, updated_at, password_hash
	FROM users
	WHERE nickname = $1 OR email = $1
	LIMIT 1;
	`

	err := u.db.QueryRowx(query, identifier).StructScan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserStore) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3);
		`

	_, err := u.db.Exec(query, userID, token, expiresAt)
	return err
}

func (u *UserStore) GetUserByRefreshToken(refreshToken string) (*models.User, error) {
	query := `
	SELECT id, nickname, email, real_name, birth_date, created_at, updated_at
	FROM users
	WHERE id = (
		SELECT user_id
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > $2
		LIMIT 1
	);
	`

	var user models.User
	err := u.db.QueryRowx(query, refreshToken, time.Now()).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) DeleteRefreshToken(refreshToken string) error {
	query := `
	DELETE FROM refresh_tokens WHERE token = $1;
	`
	_, err := u.db.Exec(query, refreshToken)
	return err
}
