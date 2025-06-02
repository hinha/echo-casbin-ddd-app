package models

import (
	"time"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"gorm.io/gorm"
)

// User is the GORM model for users
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Username  string         `gorm:"uniqueIndex;size:255;not null"`
	Email     string         `gorm:"uniqueIndex;size:255;not null"`
	Password  string         `gorm:"size:255;not null"`
	Role      string         `gorm:"size:50;not null"`
	Active    bool           `gorm:"default:true"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for User
func (*User) TableName() string {
	return "public.users"
}

// ToEntity converts the model to a domain entity
func (u *User) ToEntity() *entity.User {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &entity.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

// FromEntity updates the model from a domain entity
func (u *User) FromEntity(user *entity.User) {
	u.Username = user.Username
	u.Email = user.Email
	u.Password = user.Password
	u.Role = user.Role
	u.Active = user.Active
	u.UpdatedAt = user.UpdatedAt

	if user.DeletedAt != nil {
		u.DeletedAt = gorm.DeletedAt{Time: *user.DeletedAt, Valid: true}
	} else {
		u.DeletedAt = gorm.DeletedAt{Valid: false}
	}
}
