package repository

import (
	"context"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
)

// APIClientRepository defines the interface for API client repository
type APIClientRepository interface {
	// Create creates a new API client
	Create(ctx context.Context, client *entity.APIClient) error

	// GetByID retrieves an API client by ID
	GetByID(ctx context.Context, id uint) (*entity.APIClient, error)

	// GetByAPIKey retrieves an API client by API key
	GetByAPIKey(ctx context.Context, apiKey string) (*entity.APIClient, error)

	// Update updates an API client
	Update(ctx context.Context, client *entity.APIClient) error

	// Delete soft deletes an API client
	Delete(ctx context.Context, id uint) error

	// List retrieves all API clients with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.APIClient, int64, error)

	// GetDeletedByID retrieves a soft deleted API client by ID
	GetDeletedByID(ctx context.Context, id uint) (*entity.APIClient, error)

	// ListDeleted retrieves all soft deleted API clients with pagination
	ListDeleted(ctx context.Context, offset, limit int) ([]*entity.APIClient, int64, error)

	// Restore restores a soft deleted API client
	Restore(ctx context.Context, id uint) error

	// PermanentDelete permanently deletes an API client
	PermanentDelete(ctx context.Context, id uint) error
}
