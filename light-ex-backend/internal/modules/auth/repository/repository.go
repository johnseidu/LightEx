// Package repository defines persistence contracts for the authentication
// module. It contains interfaces only and is independent of any database
// implementation.
package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/light-group/light-ex-backend/internal/modules/auth/domain"
)

// Repository defines the persistence operations required by the
// authentication module.
//
// Implementations may use PostgreSQL, MySQL, or another storage engine
// without affecting the service layer.
type Repository interface {

	// Create stores a new user.
	Create(ctx context.Context, user *domain.User) error

	// Update persists changes to an existing user.
	Update(ctx context.Context, user *domain.User) error

	// FindByID returns a user by their unique identifier.
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// FindByEmail returns a user by their normalized email address.
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// ExistsByEmail reports whether a normalized email is already registered.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Delete removes a user.
	//
	// This is intended for development, testing, or administrative
	// operations only. In production, LightEx should generally deactivate
	// or suspend users rather than permanently deleting them.
	Delete(ctx context.Context, id uuid.UUID) error
}
