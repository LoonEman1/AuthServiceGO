package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("ошибка создания инстанса миграций: %w", err)
	}

	migrationTime := time.Now()
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Миграции не требуются, схема базы данных актуальна")
			return nil
		}
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	log.Printf("Миграции успешно применены!, время выполнения: %v\n", time.Since(migrationTime))
	return nil
}
