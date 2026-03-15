package database

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveVerificationCode(userID int, code string, ttl time.Duration) error {
	query := `
	INSERT INTO email_verification_codes (user_id, code, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id)
	DO UPDATE SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at
	`

	expiresAt := time.Now().Add(ttl)

	_, err := s.db.Exec(query, userID, code, expiresAt)
	return err
}

func (s *Store) VerifyAndActivateUser(userID int, inputCode string) (bool, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return false, errors.New("Ошибка начала транзакции")
	}
	defer tx.Rollback()

	var storedCode string
	var expiresAt time.Time

	query := `SELECT code, expires_at FROM email_verification_codes
	WHERE user_id = $1`
	err = tx.QueryRow(query, userID).Scan(&storedCode, &expiresAt)
	if err != nil {
		return false, errors.New("Не найден код в базе данных")
	}

	if storedCode != inputCode {
		return false, errors.New("Введен неверный код")
	}

	if time.Now().After(expiresAt) {
		return false, errors.New("Код недействителен")
	}

	_, err = tx.Exec(
		`
	UPDATE users SET is_verified = true 
	WHERE id = $1
	`, userID)
	if err != nil {
		return false, errors.New("Ошибка активации пользователя")
	}

	_, err = tx.Exec(`
	DELETE FROM email_verification_codes 
	WHERE user_id = $1
	`, userID)
	if err != nil {
		return false, errors.New("Ошибка удаления из бд верификационного кода, после активации пользователя")
	}

	return true, tx.Commit()
}
