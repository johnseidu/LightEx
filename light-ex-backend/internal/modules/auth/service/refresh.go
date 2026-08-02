package service

import (
	"context"
)

// RefreshResult contains newly generated tokens.
type RefreshResult struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int64
}

// Refresh validates a refresh token and issues new tokens.
func (s *Service) Refresh(
	ctx context.Context,
	refreshToken string,
) (*RefreshResult, error) {

	claims, err := s.jwt.Validate(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GenerateAccessToken(
		claims.UserID,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(
		claims.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwt.AccessTokenTTL().Seconds()),
	}, nil
}
