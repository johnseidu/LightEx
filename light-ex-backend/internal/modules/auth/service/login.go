package service

import "github.com/light-group/light-ex-backend/internal/modules/auth/domain"

// LoginResult is returned after a successful login.
type LoginResult struct {
	User *domain.User

	AccessToken  string
	RefreshToken string

	ExpiresIn int64
}