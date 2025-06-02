package dto

import (
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
)

// User DTOs

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Username string
	Email    string
	Password string
	Role     string
}

// RegisterOutput represents the output for user registration
type RegisterOutput struct {
	User  *entity.User
	Token string
}

// LoginInput represents the input for user login
type LoginInput struct {
	Username string
	Password string
}

// LoginOutput represents the output for user login
type LoginOutput struct {
	User  *entity.User
	Token string
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	ID       uint
	Username string
	Email    string
}

// ChangePasswordInput represents the input for changing a user's password
type ChangePasswordInput struct {
	ID          uint
	OldPassword string
	NewPassword string
}

// ListUsersInput represents the input for listing users
type ListUsersInput struct {
	Page  int
	Limit int
}

// ListUsersOutput represents the output for listing users
type ListUsersOutput struct {
	Users      []*entity.User
	TotalCount int64
}
