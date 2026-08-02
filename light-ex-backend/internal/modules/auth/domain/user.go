package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusPending   UserStatus = "PENDING"
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusLocked    UserStatus = "LOCKED"
)

var (
	ErrUserAlreadyVerified = errors.New("user is already verified")
	ErrInvalidUserEmail    = errors.New("invalid user email")
)

// User represents an authenticated identity within LightEx.
//
// It intentionally contains only authentication-related data.
// Wallets, KYC, trading limits, balances and profile information
// belong to other modules.
type User struct {
	ID             uuid.UUID
	Email          string
	PasswordHash   string

	EmailVerified bool
	Status        UserStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user in the pending state.
func NewUser(email, passwordHash string) (*User, error) {

	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, ErrInvalidUserEmail
	}

	now := time.Now().UTC()

	return &User{
		ID:             uuid.New(),
		Email:          email,
		PasswordHash:   passwordHash,
		EmailVerified:  false,
		Status:         UserStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// VerifyEmail marks the user's email as verified.
func (u *User) VerifyEmail() error {

	if u.EmailVerified {
		return ErrUserAlreadyVerified
	}

	u.EmailVerified = true
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now().UTC()

	return nil
}

// Suspend disables the account.
func (u *User) Suspend() {
	u.Status = UserStatusSuspended
	u.UpdatedAt = time.Now().UTC()
}

// Lock temporarily locks the account.
func (u *User) Lock() {
	u.Status = UserStatusLocked
	u.UpdatedAt = time.Now().UTC()
}

// Activate activates the account.
func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now().UTC()
}