package interfaces

import (
	"context"

	"github.com/hinha/echo-casbin-ddd-app/internal/application/dto"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
)

// UserUseCase defines the interface for user-related business logic
type UserUseCase interface {
	// Register registers a new user
	Register(ctx context.Context, input dto.RegisterInput) (*dto.RegisterOutput, error)

	// Login authenticates a user
	Login(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, error)

	// GetUserByID gets a user by ID
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)

	// UpdateUser updates a user
	UpdateUser(ctx context.Context, input dto.UpdateUserInput) (*entity.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error

	// ListUsers lists users with pagination
	ListUsers(ctx context.Context, input dto.ListUsersInput) (*dto.ListUsersOutput, error)
}
