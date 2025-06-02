package dto

import (
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
)

// API Client DTOs

// CreateAPIClientInput represents the input for creating an API client
type CreateAPIClientInput struct {
	Name        string
	Description string
}

// CreateAPIClientOutput represents the output for creating an API client
type CreateAPIClientOutput struct {
	APIClient *entity.APIClient
}

// GetAPIClientByIDInput represents the input for getting an API client by ID
type GetAPIClientByIDInput struct {
	ID uint
}

// GetAPIClientByAPIKeyInput represents the input for getting an API client by API key
type GetAPIClientByAPIKeyInput struct {
	APIKey string
}

// UpdateAPIClientInput represents the input for updating an API client
type UpdateAPIClientInput struct {
	ID          uint
	Name        string
	Description string
}

// RegenerateAPIKeyInput represents the input for regenerating an API key
type RegenerateAPIKeyInput struct {
	ID uint
}

// SetAPIClientActiveInput represents the input for setting an API client's active status
type SetAPIClientActiveInput struct {
	ID     uint
	Active bool
}

// DeleteAPIClientInput represents the input for deleting an API client
type DeleteAPIClientInput struct {
	ID uint
}

// ListAPIClientsInput represents the input for listing API clients
type ListAPIClientsInput struct {
	Page  int
	Limit int
}

// ListAPIClientsOutput represents the output for listing API clients
type ListAPIClientsOutput struct {
	APIClients []*entity.APIClient
	TotalCount int64
}
