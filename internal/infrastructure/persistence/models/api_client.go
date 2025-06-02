package models

import (
	"time"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"gorm.io/gorm"
)

// APIClient is the GORM model for API clients
type APIClient struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"size:255;not null"`
	Description string         `gorm:"size:1000"`
	APIKey      string         `gorm:"uniqueIndex;size:255;not null"`
	Active      bool           `gorm:"default:true"`
	CreatedAt   time.Time      `gorm:"not null"`
	UpdatedAt   time.Time      `gorm:"not null"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for APIClient
func (*APIClient) TableName() string {
	return "public.api_clients"
}

// ToEntity converts the model to a domain entity
func (c *APIClient) ToEntity() *entity.APIClient {
	var deletedAt *time.Time
	if c.DeletedAt.Valid {
		deletedAt = &c.DeletedAt.Time
	}

	return &entity.APIClient{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		APIKey:      c.APIKey,
		Active:      c.Active,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// FromEntity updates the model from a domain entity
func (c *APIClient) FromEntity(client *entity.APIClient) {
	c.Name = client.Name
	c.Description = client.Description
	c.APIKey = client.APIKey
	c.Active = client.Active
	c.UpdatedAt = client.UpdatedAt

	if client.DeletedAt != nil {
		c.DeletedAt = gorm.DeletedAt{Time: *client.DeletedAt, Valid: true}
	} else {
		c.DeletedAt = gorm.DeletedAt{Valid: false}
	}
}
