package database

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type CodesStore struct {
	db *sqlx.DB
}

func NewCodeStore(db *sqlx.DB) *CodesStore {
	return &CodesStore{db: db}
}

func (s *CodesStore) SaveVerificationCode(userID int, code string, ttl time.Duration) error {
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

func (s *CodesStore) VerifyAndActivateUser(userID int, inputCode string) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return errors.New("Ошибка начала транзакции")
	}
	defer tx.Rollback()

	var storedCode string
	var expiresAt time.Time

	query := `SELECT code, expires_at FROM email_verification_codes
	WHERE user_id = $1`
	err = tx.QueryRow(query, userID).Scan(&storedCode, &expiresAt)
	if err != nil {
		return errors.New("Не найден код в базе данных")
	}

	if storedCode != inputCode {
		return errors.New("Введен неверный код")
	}

	if time.Now().After(expiresAt) {
		return errors.New("Код недействителен")
	}

	_, err = tx.Exec(
		`
	UPDATE users SET is_verified = true 
	WHERE id = $1
	`, userID)
	if err != nil {
		return errors.New("Ошибка активации пользователя")
	}

	_, err = tx.Exec(`
	DELETE FROM email_verification_codes 
	WHERE user_id = $1
	`, userID)
	if err != nil {
		return errors.New("Ошибка удаления из бд верификационного кода, после активации пользователя")
	}

	return tx.Commit()
}
