package jwtPckg

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secretKey string
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{secretKey: secret}
}

func (m *TokenManager) NewAccessToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString([]byte(m.secretKey))
}

func (m *TokenManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// returning access token, refresh token and return error if can't generate tokens
func (m *TokenManager) GenerateTokens(userID int) (*string, *string, error) {
	accessToken, err := m.NewAccessToken(userID)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := m.NewRefreshToken()
	if err != nil {
		return nil, nil, err
	}
	return &accessToken, &refreshToken, nil
}
