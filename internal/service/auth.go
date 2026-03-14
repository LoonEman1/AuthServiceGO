package service

import (
	"AuthService/internal/database"
	"AuthService/internal/hashPassword"
	jwtPckg "AuthService/internal/jwt"
	"AuthService/internal/models"
	"errors"
	"time"
)

type AuthService struct {
	userStore    *database.UserStore
	tokenManager *jwtPckg.TokenManager
}

func NewAuthService(store *database.UserStore, tokenManager *jwtPckg.TokenManager) *AuthService {
	return &AuthService{userStore: store, tokenManager: tokenManager}
}

func (s *AuthService) Register(input models.RegisterUserInput) (*models.AuthResponse, error) {
	hashedPassword, err := hashPassword.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	userModel := input.ToUser(hashedPassword)

	user, err := s.userStore.Register(*userModel)

	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenManager.NewAccessToken(user.ID)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.userStore.SaveRefreshToken(user.ID, refreshToken, time.Now().Add(time.Hour*24*30))
	if err != nil {
		return nil, err
	}

	return models.NewAuthResponse(user, accessToken, refreshToken), nil
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

	err = hashPassword.CheckPasswordHash(input.Password, user.PasswordHash)
	if err != nil {
		return nil, errors.New("Неверный пароль")
	}

	accessToken, refreshToken, err := s.tokenManager.GenerateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	return models.NewAuthResponse(user, *accessToken, *refreshToken), nil
}

func (s *AuthService) Logout(input models.RefreshInput) error {
	err := s.userStore.DeleteRefreshToken(input.RefreshToken)

	return err
}
