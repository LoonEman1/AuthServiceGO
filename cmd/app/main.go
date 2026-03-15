package main

import (
	"AuthService/internal/database"
	"AuthService/internal/handlers"
	jwtPckg "AuthService/internal/jwt"
	"AuthService/internal/queue"
	"AuthService/internal/service"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret-key-change-me"
	}

	log.Printf("Сервер начинает запуск")

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	emailProducer := queue.NewKafkaProducer([]string{kafkaBrokers}, "email-notifications")
	defer emailProducer.Close()

	err = emailProducer.EnsureTopicExists(context.Background(), []string{kafkaBrokers}, "email-notifications")
	if err != nil {
		log.Printf("Предупреждение по Кафке: %v", err)
	}

	if err := database.RunMigrations(databaseURL); err != nil {
		log.Fatalf("Критическая ошибка миграций: %v", err)
	}

	log.Println("Успешно подключено к бд")

	userStore := database.NewUserStore(db)
	codesStore := database.NewCodeStore(db)

	tokenManager := jwtPckg.NewTokenManager(jwtSecret)

	authService := service.NewAuthService(userStore, tokenManager, codesStore, emailProducer)

	handler := handlers.NewHandler(authService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", handler.Register)
	mux.HandleFunc("POST /api/v1/auth/refresh", handler.Refresh)
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", handler.Logout)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Сервер запущен на порту %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}

}
