package repository

import (
	"context"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
)

// UserRepository defines the interface for user repository
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uint) (*entity.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uint) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)

	// GetDeletedByID retrieves a soft deleted user by ID
	GetDeletedByID(ctx context.Context, id uint) (*entity.User, error)

	// ListDeleted retrieves all soft deleted users with pagination
	ListDeleted(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)

	// Restore restores a soft deleted user
	Restore(ctx context.Context, id uint) error

	// PermanentDelete permanently deletes a user
	PermanentDelete(ctx context.Context, id uint) error
}
