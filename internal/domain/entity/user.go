package entity

import (
	"errors"
	"time"

	"github.com/hinha/echo-casbin-ddd-app/pkg/argon2"
)

// User represents a user in the system
type User struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // Password is not exposed in JSON
	Role      string     `json:"role"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// NewUser creates a new user
func NewUser(username, email, password, role string) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}

	hashedPassword, err := argon2.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:  username,
		Email:     email,
		Password:  hashedPassword,
		Role:      role,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// ValidatePassword checks if the provided password matches the stored hash
func (u *User) ValidatePassword(password string) bool {
	isValid, _ := argon2.VerifyHash(password, u.Password)
	return isValid
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := argon2.GenerateHash(password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(username, email string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if email == "" {
		return errors.New("email cannot be empty")
	}

	u.Username = username
	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

// SetRole sets the user's role
func (u *User) SetRole(role string) {
	u.Role = role
	u.UpdatedAt = time.Now()
}

// SetActive sets the user's active status
func (u *User) SetActive(active bool) {
	u.Active = active
	u.UpdatedAt = time.Now()
}
