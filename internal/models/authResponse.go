package models

type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func NewAuthResponse(user *User, accessToken, refreshToken string) *AuthResponse {
	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func NewRegisterResponse(user *User, message string) *RegisterResponse {
	return &RegisterResponse{
		Email:   user.Email,
		Message: message,
	}
}
