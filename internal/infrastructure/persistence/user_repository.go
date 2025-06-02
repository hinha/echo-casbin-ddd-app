package persistence

import (
	"context"
	"errors"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/repository"
	"github.com/hinha/echo-casbin-ddd-app/internal/infrastructure/persistence/models"
	"gorm.io/gorm"
)

// UserRepository is the implementation of repository.UserRepository
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	model := &models.User{}
	model.FromEntity(user)
	model.ID = 0 // Ensure ID is not set for creation

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return result.Error
	}

	user.ID = model.ID
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var model models.User
	result := r.db.WithContext(ctx).First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var model models.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var model models.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	model := &models.User{}
	model.FromEntity(user)
	model.ID = user.ID

	result := r.db.WithContext(ctx).Save(model)
	return result.Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	return result.Error
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var models []models.User
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	users := make([]*entity.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	return users, count, nil
}

// GetDeletedByID retrieves a soft deleted user by ID
func (r *UserRepository) GetDeletedByID(ctx context.Context, id uint) (*entity.User, error) {
	var model models.User
	result := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// ListDeleted retrieves all soft deleted users with pagination
func (r *UserRepository) ListDeleted(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("deleted_at IS NOT NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var models []models.User
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	users := make([]*entity.User, len(models))
	for i, model := range models {
		users[i] = model.ToEntity()
	}

	return users, count, nil
}

// Restore restores a soft deleted user
func (r *UserRepository) Restore(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// PermanentDelete permanently deletes a user
func (r *UserRepository) PermanentDelete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, id)
	return result.Error
}
