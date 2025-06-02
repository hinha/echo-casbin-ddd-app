package auth

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/hinha/echo-casbin-ddd-app/internal/config"
	"gorm.io/gorm"
)

// CasbinService handles authorization using Casbin
type CasbinService struct {
	enforcer *casbin.SyncedEnforcer
}

// NewCasbinService creates a new CasbinService
func NewCasbinService(db *gorm.DB, config *config.Config) (*CasbinService, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedEnforcer("casbin/model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// Load policies from the database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	//enforcer.StartAutoLoadPolicy(2 * time.Second)

	return &CasbinService{
		enforcer: enforcer,
	}, nil
}

// Enforce checks if a subject can access an object with the given action in the specified domain
func (s *CasbinService) Enforce(sub, dom, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, dom, obj, act)
}

// AddPolicy adds a policy rule to the enforcer
func (s *CasbinService) AddPolicy(sub, dom, obj, act string) (bool, error) {
	return s.enforcer.AddPolicy(sub, dom, obj, act)
}

// RemovePolicy removes a policy rule from the enforcer
func (s *CasbinService) RemovePolicy(sub, dom, obj, act string) (bool, error) {
	return s.enforcer.RemovePolicy(sub, dom, obj, act)
}

// AddRoleForUser adds a role for a user in a domain
func (s *CasbinService) AddRoleForUser(user, role, domain string) (bool, error) {
	return s.enforcer.AddGroupingPolicy(user, role, domain)
}

// DeleteRoleForUser removes a role for a user in a domain
func (s *CasbinService) DeleteRoleForUser(user, role, domain string) (bool, error) {
	return s.enforcer.RemoveGroupingPolicy(user, role, domain)
}

// GetRolesForUser gets roles for a user in a domain
func (s *CasbinService) GetRolesForUser(user, domain string) ([]string, error) {
	return s.enforcer.GetRolesForUserInDomain(user, domain), nil
}

// GetUsersForRole gets users for a role in a domain
func (s *CasbinService) GetUsersForRole(role, domain string) ([]string, error) {
	return s.enforcer.GetUsersForRoleInDomain(role, domain), nil
}
