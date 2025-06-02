package auth

import (
	"context"
	"errors"

	"github.com/hinha/echo-casbin-ddd-app/internal/config"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/repository"
)

// APIKeyService handles API key authentication
type APIKeyService struct {
	config     *config.Config
	repository repository.APIClientRepository
}

// NewAPIKeyService creates a new APIKeyService
func NewAPIKeyService(config *config.Config, repository repository.APIClientRepository) *APIKeyService {
	return &APIKeyService{
		config:     config,
		repository: repository,
	}
}

// ValidateAPIKey validates an API key
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, apiKey string) (bool, uint, error) {
	if apiKey == "" {
		return false, 0, errors.New("API key is required")
	}

	client, err := s.repository.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return false, 0, err
	}

	if client == nil {
		return false, 0, errors.New("invalid API key")
	}

	if !client.Active {
		return false, 0, errors.New("API client is inactive")
	}

	return true, client.ID, nil
}

// GetAPIKeyFromHeader extracts the API key from the request header
func (s *APIKeyService) GetAPIKeyFromHeader(header string) string {
	return header
}
