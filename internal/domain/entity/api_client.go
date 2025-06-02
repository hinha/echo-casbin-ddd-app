package entity

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// APIClient represents an API client in the system
type APIClient struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	APIKey      string     `json:"api_key"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// NewAPIClient creates a new API client
func NewAPIClient(name, description string) (*APIClient, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	return &APIClient{
		Name:        name,
		Description: description,
		APIKey:      apiKey,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// RegenerateAPIKey generates a new API key for the client
func (c *APIClient) RegenerateAPIKey() error {
	apiKey, err := generateAPIKey()
	if err != nil {
		return err
	}

	c.APIKey = apiKey
	c.UpdatedAt = time.Now()
	return nil
}

// SetActive sets the client's active status
func (c *APIClient) SetActive(active bool) {
	c.Active = active
	c.UpdatedAt = time.Now()
}

// UpdateInfo updates the client's information
func (c *APIClient) UpdateInfo(name, description string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	c.Name = name
	c.Description = description
	c.UpdatedAt = time.Now()
	return nil
}

// generateAPIKey generates a random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
