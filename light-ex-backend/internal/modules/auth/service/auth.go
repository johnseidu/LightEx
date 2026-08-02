package service

import (
	"context"
	"errors"
	"strings"

	"github.com/light-group/light-ex-backend/internal/modules/auth/domain"
	"github.com/light-group/light-ex-backend/internal/modules/auth/repository"
	jwt "github.com/light-group/light-ex-backend/internal/platform/security/jwt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountLocked      = errors.New("account is locked")
	ErrAccountSuspended   = errors.New("account is suspended")
)

type Service struct {
	repository repository.Repository
	password   PasswordService
	jwt        jwt.Service
}

func New(
	repo repository.Repository,
	password PasswordService,
	jwtService jwt.Service,
) *Service {

	return &Service{
		repository: repo,
		password:   password,
		jwt:        jwtService,
	}
}

// Register registers a new user.
func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (*domain.User, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	exists, err := s.repository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := s.password.Hash(password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(email, hash)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user.
func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (*LoginResult, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	ok, err := s.password.Verify(
		password,
		user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrInvalidCredentials
	}

	switch user.Status {

	case domain.UserStatusLocked:
		return nil, ErrAccountLocked

	case domain.UserStatusSuspended:
		return nil, ErrAccountSuspended
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwt.AccessTokenTTL().Seconds()),
	}, nil
}
