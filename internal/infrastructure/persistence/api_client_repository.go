package persistence

import (
	"context"
	"errors"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/repository"
	"github.com/hinha/echo-casbin-ddd-app/internal/infrastructure/persistence/models"
	"gorm.io/gorm"
)

// APIClientRepository is the implementation of repository.APIClientRepository
type APIClientRepository struct {
	db *gorm.DB
}

// NewAPIClientRepository creates a new APIClientRepository
func NewAPIClientRepository(db *gorm.DB) repository.APIClientRepository {
	return &APIClientRepository{
		db: db,
	}
}

// Create creates a new API client
func (r *APIClientRepository) Create(ctx context.Context, client *entity.APIClient) error {
	model := &models.APIClient{}
	model.FromEntity(client)
	model.ID = 0 // Ensure ID is not set for creation

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return result.Error
	}

	client.ID = model.ID
	return nil
}

// GetByID retrieves an API client by ID
func (r *APIClientRepository) GetByID(ctx context.Context, id uint) (*entity.APIClient, error) {
	var model models.APIClient
	result := r.db.WithContext(ctx).First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// GetByAPIKey retrieves an API client by API key
func (r *APIClientRepository) GetByAPIKey(ctx context.Context, apiKey string) (*entity.APIClient, error) {
	var model models.APIClient
	result := r.db.WithContext(ctx).Where("api_key = ?", apiKey).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// Update updates an API client
func (r *APIClientRepository) Update(ctx context.Context, client *entity.APIClient) error {
	model := &models.APIClient{}
	model.FromEntity(client)
	model.ID = client.ID

	result := r.db.WithContext(ctx).Save(model)
	return result.Error
}

// Delete soft deletes an API client
func (r *APIClientRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.APIClient{}, id)
	return result.Error
}

// List retrieves all API clients with pagination
func (r *APIClientRepository) List(ctx context.Context, offset, limit int) ([]*entity.APIClient, int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.APIClient{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var models []models.APIClient
	result := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	clients := make([]*entity.APIClient, len(models))
	for i, model := range models {
		clients[i] = model.ToEntity()
	}

	return clients, count, nil
}

// GetDeletedByID retrieves a soft deleted API client by ID
func (r *APIClientRepository) GetDeletedByID(ctx context.Context, id uint) (*entity.APIClient, error) {
	var model models.APIClient
	result := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.ToEntity(), nil
}

// ListDeleted retrieves all soft deleted API clients with pagination
func (r *APIClientRepository) ListDeleted(ctx context.Context, offset, limit int) ([]*entity.APIClient, int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Unscoped().Model(&models.APIClient{}).Where("deleted_at IS NOT NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var models []models.APIClient
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Offset(offset).Limit(limit).Find(&models)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	clients := make([]*entity.APIClient, len(models))
	for i, model := range models {
		clients[i] = model.ToEntity()
	}

	return clients, count, nil
}

// Restore restores a soft deleted API client
func (r *APIClientRepository) Restore(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.APIClient{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// PermanentDelete permanently deletes an API client
func (r *APIClientRepository) PermanentDelete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.APIClient{}, id)
	return result.Error
}
