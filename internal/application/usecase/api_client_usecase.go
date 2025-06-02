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

// APIClientUseCaseImpl handles API client-related business logic
// It implements the interfaces.APIClientUseCase interface
type APIClientUseCaseImpl struct {
	apiClientRepository repository.APIClientRepository
	casbinService       *auth.CasbinService
}

// NewAPIClientUseCase creates a new APIClientUseCaseImpl
func NewAPIClientUseCase(
	apiClientRepository repository.APIClientRepository,
	casbinService *auth.CasbinService,
) interfaces.APIClientUseCase {
	return &APIClientUseCaseImpl{
		apiClientRepository: apiClientRepository,
		casbinService:       casbinService,
	}
}

// Create creates a new API client
func (uc *APIClientUseCaseImpl) Create(ctx context.Context, input dto.CreateAPIClientInput) (*dto.CreateAPIClientOutput, error) {
	// Create new API client
	client, err := entity.NewAPIClient(input.Name, input.Description)
	if err != nil {
		return nil, err
	}

	// Save API client to database
	if err := uc.apiClientRepository.Create(ctx, client); err != nil {
		return nil, err
	}

	// Add policy for API client in Casbin
	if _, err := uc.casbinService.AddPolicy(client.Name, "api", "/api/*", "GET"); err != nil {
		return nil, err
	}

	return &dto.CreateAPIClientOutput{
		APIClient: client,
	}, nil
}

// GetByID gets an API client by ID
func (uc *APIClientUseCaseImpl) GetByID(ctx context.Context, input dto.GetAPIClientByIDInput) (*entity.APIClient, error) {
	return uc.apiClientRepository.GetByID(ctx, input.ID)
}

// GetByAPIKey gets an API client by API key
func (uc *APIClientUseCaseImpl) GetByAPIKey(ctx context.Context, input dto.GetAPIClientByAPIKeyInput) (*entity.APIClient, error) {
	return uc.apiClientRepository.GetByAPIKey(ctx, input.APIKey)
}

// Update updates an API client
func (uc *APIClientUseCaseImpl) Update(ctx context.Context, input dto.UpdateAPIClientInput) (*entity.APIClient, error) {
	// Get API client by ID
	client, err := uc.apiClientRepository.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("API client not found")
	}

	// Update API client
	if err := client.UpdateInfo(input.Name, input.Description); err != nil {
		return nil, err
	}

	// Save API client to database
	if err := uc.apiClientRepository.Update(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// RegenerateAPIKey regenerates an API key for an API client
func (uc *APIClientUseCaseImpl) RegenerateAPIKey(ctx context.Context, input dto.RegenerateAPIKeyInput) (*entity.APIClient, error) {
	// Get API client by ID
	client, err := uc.apiClientRepository.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("API client not found")
	}

	// Regenerate API key
	if err := client.RegenerateAPIKey(); err != nil {
		return nil, err
	}

	// Save API client to database
	if err := uc.apiClientRepository.Update(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// SetActive sets an API client's active status
func (uc *APIClientUseCaseImpl) SetActive(ctx context.Context, input dto.SetAPIClientActiveInput) (*entity.APIClient, error) {
	// Get API client by ID
	client, err := uc.apiClientRepository.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("API client not found")
	}

	// Set active status
	client.SetActive(input.Active)

	// Save API client to database
	if err := uc.apiClientRepository.Update(ctx, client); err != nil {
		return nil, err
	}

	return client, nil
}

// Delete deletes an API client
func (uc *APIClientUseCaseImpl) Delete(ctx context.Context, input dto.DeleteAPIClientInput) error {
	// Get API client by ID
	client, err := uc.apiClientRepository.GetByID(ctx, input.ID)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.New("API client not found")
	}

	// Remove policies for API client in Casbin
	if _, err := uc.casbinService.RemovePolicy(client.Name, "api", "/api/*", "GET"); err != nil {
		return err
	}

	// Delete API client from database
	return uc.apiClientRepository.Delete(ctx, input.ID)
}

// List lists API clients with pagination
func (uc *APIClientUseCaseImpl) List(ctx context.Context, input dto.ListAPIClientsInput) (*dto.ListAPIClientsOutput, error) {
	// Calculate offset
	offset := (input.Page - 1) * input.Limit

	// Get API clients from database
	clients, count, err := uc.apiClientRepository.List(ctx, offset, input.Limit)
	if err != nil {
		return nil, err
	}

	return &dto.ListAPIClientsOutput{
		APIClients: clients,
		TotalCount: count,
	}, nil
}
