package usecase

import (
	"context"
	"errors"

	"github.com/hinha/echo-casbin-ddd-app/internal/application/dto"
	"github.com/hinha/echo-casbin-ddd-app/internal/application/interfaces"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/repository"
	"github.com/hinha/echo-casbin-ddd-app/internal/infrastructure/auth"
)

// UserUseCaseImpl handles user-related business logic
// It implements the interfaces.UserUseCase interface
type UserUseCaseImpl struct {
	userRepository repository.UserRepository
	jwtService     *auth.JWTService
	casbinService  *auth.CasbinService
}

// NewUserUseCase creates a new UserUseCaseImpl
func NewUserUseCase(
	userRepository repository.UserRepository,
	jwtService *auth.JWTService,
	casbinService *auth.CasbinService,
) interfaces.UserUseCase {
	return &UserUseCaseImpl{
		userRepository: userRepository,
		jwtService:     jwtService,
		casbinService:  casbinService,
	}
}

// Register registers a new user
func (uc *UserUseCaseImpl) Register(ctx context.Context, input dto.RegisterInput) (*dto.RegisterOutput, error) {
	// Check if user with the same username already exists
	existingUser, err := uc.userRepository.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if user with the same email already exists
	existingUser, err = uc.userRepository.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Create new user
	user, err := entity.NewUser(input.Username, input.Email, input.Password, input.Role)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := uc.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	// Add role for user in Casbin
	if _, err := uc.casbinService.AddRoleForUser(user.Username, user.Role, "default"); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterOutput{
		User:  user,
		Token: token,
	}, nil
}

// Login authenticates a user
func (uc *UserUseCaseImpl) Login(ctx context.Context, input dto.LoginInput) (*dto.LoginOutput, error) {
	// Get user by username
	user, err := uc.userRepository.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Validate password
	if !user.ValidatePassword(input.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("user is inactive")
	}

	// Generate JWT token
	token, err := uc.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginOutput{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID gets a user by ID
func (uc *UserUseCaseImpl) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return uc.userRepository.GetByID(ctx, id)
}

// UpdateUser updates a user
func (uc *UserUseCaseImpl) UpdateUser(ctx context.Context, input dto.UpdateUserInput) (*entity.User, error) {
	// Get user by ID
	user, err := uc.userRepository.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Update user
	if err := user.UpdateProfile(input.Username, input.Email); err != nil {
		return nil, err
	}

	// Save user to database
	if err := uc.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (uc *UserUseCaseImpl) ChangePassword(ctx context.Context, input dto.ChangePasswordInput) error {
	// Get user by ID
	user, err := uc.userRepository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Validate old password
	if !user.ValidatePassword(input.OldPassword) {
		return errors.New("invalid old password")
	}

	// Change password
	if err := user.ChangePassword(input.NewPassword); err != nil {
		return err
	}

	// Save user to database
	return uc.userRepository.Update(ctx, user)
}

// ListUsers lists users with pagination
func (uc *UserUseCaseImpl) ListUsers(ctx context.Context, input dto.ListUsersInput) (*dto.ListUsersOutput, error) {
	// Calculate offset
	offset := (input.Page - 1) * input.Limit

	// Get users from database
	users, count, err := uc.userRepository.List(ctx, offset, input.Limit)
	if err != nil {
		return nil, err
	}

	return &dto.ListUsersOutput{
		Users:      users,
		TotalCount: count,
	}, nil
}
