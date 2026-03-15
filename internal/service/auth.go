package service

import (
	"AuthService/internal/database"
	"AuthService/internal/hashPassword"
	jwtPckg "AuthService/internal/jwt"
	"AuthService/internal/models"
	"AuthService/internal/queue"
	"AuthService/internal/utils"
	"context"
	"errors"
	"log"
	"time"
)

type AuthService struct {
	userStore    *database.UserStore
	tokenManager *jwtPckg.TokenManager
	codesStore   *database.CodesStore
	producer     *queue.KafkaProducer
}

func NewAuthService(store *database.UserStore, tokenManager *jwtPckg.TokenManager, codesStore *database.CodesStore, producer *queue.KafkaProducer) *AuthService {
	return &AuthService{userStore: store, tokenManager: tokenManager, codesStore: codesStore, producer: producer}
}

func (s *AuthService) Register(input models.RegisterUserInput) (*models.User, error) {
	hashedPassword, err := hashPassword.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	userModel := input.ToUser(hashedPassword)

	user, err := s.userStore.Register(*userModel)

	if err != nil {
		return nil, err
	}

	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	err = s.codesStore.SaveVerificationCode(user.ID, code, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	go s.notifyUserCreated(context.Background(), user, code)

	// accessToken, err := s.tokenManager.NewAccessToken(user.ID)

	// if err != nil {
	// 	return nil, err
	// }

	// refreshToken, err := s.tokenManager.NewRefreshToken()
	// if err != nil {
	// 	return nil, err
	// }

	// err = s.userStore.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(time.Hour*24*30))
	// if err != nil {
	// 	return nil, err
	// }

	return user, nil
}

func (s *AuthService) notifyUserCreated(ctx context.Context, user *models.User, code string) {
	task := models.NewEmailTask(user.Email, code, "Registration")

	sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.producer.SendEmailTask(sendCtx, *task); err != nil {
		log.Printf("Ошибка отправки email task в кафку для юзера %d: %v", user.ID, err)
	}

}

func (s *AuthService) NewEmailConfirmationCode(input models.GenerateNewCodeInput) error {

	user, err := s.userStore.GetUserByMail(input.Email)
	if err != nil {
		return errors.New("Пользователь с такой почтой не найден")
	}

	if user.IsVerified == true {
		return errors.New("Почта пользователя уже привязана")
	}

	code, err := utils.GenerateVerificationCode()
	if err != nil {
		return err
	}

	err = s.codesStore.SaveVerificationCode(user.ID, code, 15*time.Minute)
	if err != nil {
		return err
	}

	go s.notifyUserCreated(context.Background(), user, code)

	return err
}

func (s *AuthService) Verify(input models.VerifyInput) (*models.AuthResponse, error) {

	user, err := s.userStore.GetUserByMail(input.Email)
	if err != nil {
		return nil, err
	}

	if err := s.codesStore.VerifyAndActivateUser(user.ID, input.Code); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	err = s.userStore.SaveRefreshToken(user.ID, *refreshToken, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return nil, err
	}

	return models.NewAuthResponse(user, *accessToken, *refreshToken), nil
}

func (s *AuthService) Refresh(input models.RefreshInput) (*models.AuthResponse, error) {

	user, err := s.userStore.GetUserByRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New("Неверный или просроченный токен")
	}

	newAccess, newRefresh, err := s.tokenManager.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	err = s.userStore.SaveRefreshToken(user.ID, *newRefresh, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return nil, err
	}

	return models.NewAuthResponse(user, *newAccess, *newRefresh), nil
}

func (s *AuthService) Login(input models.LoginUserInput) (*models.AuthResponse, error) {
	user, err := s.userStore.GetUserByIdentifier(input.Identifier)
	if err != nil {
		return nil, err
	}

	if user.IsVerified == false {
		return nil, errors.New("Для входа в аккаунт необходимо подтвердить почту")
	}

	err = hashPassword.CheckPasswordHash(input.Password, user.PasswordHash)
	if err != nil {
		return nil, errors.New("Неверный пароль")
	}

	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	err = s.userStore.SaveRefreshToken(user.ID, *refreshToken, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return nil, err
	}

	return models.NewAuthResponse(user, *accessToken, *refreshToken), nil
}

func (s *AuthService) Logout(input models.RefreshInput) error {
	err := s.userStore.DeleteRefreshToken(input.RefreshToken)

	return err
}
