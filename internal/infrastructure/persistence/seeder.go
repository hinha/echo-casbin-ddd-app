package persistence

import (
	"context"
	"errors"
	"github.com/hinha/echo-casbin-ddd-app/internal/config"
	"log"

	"github.com/hinha/echo-casbin-ddd-app/internal/domain/entity"
	"github.com/hinha/echo-casbin-ddd-app/internal/domain/repository"
	"github.com/hinha/echo-casbin-ddd-app/pkg/argon2"
)

// Seeder handles database seeding
type Seeder struct {
	cfg      *config.Config
	userRepo repository.UserRepository
}

// NewSeeder creates a new seeder
func NewSeeder(cfg *config.Config, userRepo repository.UserRepository) *Seeder {
	return &Seeder{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

// SeedUsers seeds initial users into the database
func (s *Seeder) SeedUsers(ctx context.Context) error {
	// Check if admin user already exists
	existingAdmin, err := s.userRepo.GetByUsername(ctx, s.cfg.Server.InitialAdminUsername)
	if err != nil {
		return err
	}

	// If admin user doesn't exist, create it
	if existingAdmin == nil {
		// Create admin user with Argon2 hashed password
		password, err := argon2.GenerateHash(s.cfg.Server.InitialAdminPassword)
		if err != nil {
			return err
		}

		adminUser := &entity.User{
			Username: s.cfg.Server.InitialAdminUsername,
			Email:    "admin@example.com",
			Password: password,
			Role:     "superadmin",
			Active:   true,
		}

		if err := s.userRepo.Create(ctx, adminUser); err != nil {
			return err
		}

		// Verify that the user was created correctly by retrieving it from the database
		createdUser, err := s.userRepo.GetByUsername(ctx, s.cfg.Server.InitialAdminUsername)
		if err != nil {
			return err
		}

		if createdUser == nil {
			return errors.New("failed to retrieve created admin user")
		}

		// Verify that the password is correct
		if !createdUser.ValidatePassword(s.cfg.Server.InitialAdminPassword) {
			return errors.New("password verification failed for admin user")
		}

		log.Println("Admin user created and verified successfully")
	}

	return nil
}

// Seed seeds all initial data
func (s *Seeder) Seed(ctx context.Context) error {
	if err := s.SeedUsers(ctx); err != nil {
		return err
	}

	return nil
}
