package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrPasswordTooShort    = errors.New("password must be at least 12 characters")
)

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024 // 64 MB
	argonThreads uint8  = 4
	argonKeyLen  uint32 = 32
	saltLength          = 16
)

// PasswordService defines password hashing operations.
type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password, encodedHash string) (bool, error)
}

type Argon2PasswordService struct{}

var _ PasswordService = (*Argon2PasswordService)(nil)

// NewPasswordService creates a new password service.
func NewPasswordService() PasswordService {
	return &Argon2PasswordService{}
}

// Hash hashes a plaintext password using Argon2id.
func (s *Argon2PasswordService) Hash(password string) (string, error) {

	if len(password) < 12 {
		return "", ErrPasswordTooShort
	}

	salt := make([]byte, saltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonTime,
		argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

// Verify compares a plaintext password against an encoded hash.
func (s *Argon2PasswordService) Verify(password, encodedHash string) (bool, error) {

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false, ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(hash, expectedHash) != 1 {
		return false, nil
	}

	return true, nil
}
