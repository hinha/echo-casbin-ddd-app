package interfaces

import (
	"context"

	"github.com/hinha/echo-casbin-ddd-app/internal/application/dto"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
)

// APIClientUseCase defines the interface for API client-related business logic
type APIClientUseCase interface {
	// Create creates a new API client
	Create(ctx context.Context, input dto.CreateAPIClientInput) (*dto.CreateAPIClientOutput, error)

	// GetByID gets an API client by ID
	GetByID(ctx context.Context, input dto.GetAPIClientByIDInput) (*entity.APIClient, error)

	// GetByAPIKey gets an API client by API key
	GetByAPIKey(ctx context.Context, input dto.GetAPIClientByAPIKeyInput) (*entity.APIClient, error)

	// Update updates an API client
	Update(ctx context.Context, input dto.UpdateAPIClientInput) (*entity.APIClient, error)

	// RegenerateAPIKey regenerates an API key for an API client
	RegenerateAPIKey(ctx context.Context, input dto.RegenerateAPIKeyInput) (*entity.APIClient, error)

	// SetActive sets an API client's active status
	SetActive(ctx context.Context, input dto.SetAPIClientActiveInput) (*entity.APIClient, error)

	// Delete deletes an API client
	Delete(ctx context.Context, input dto.DeleteAPIClientInput) error

	// List lists API clients with pagination
	List(ctx context.Context, input dto.ListAPIClientsInput) (*dto.ListAPIClientsOutput, error)
}
