package domain

import (
	"context"
	"errors"
	"log"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repo.UserRepo
	roleRepo *repo.RoleRepo
}

func NewUserService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, tenantID, email, password, firstName, lastName string, roleIDs []string) (*repo.User, error) {
	log.Printf("DEBUG: user_service.CreateUser tenant=%s email=%s", tenantID, email)
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, tenantID, email)
	if err != nil {
		log.Printf("ERROR: user_service.CreateUser GetByEmail: %v", err)
		return nil, err
	}
	if existingUser != nil {
		log.Printf("WARN: user_service.CreateUser email already exists tenant=%s email=%s", tenantID, email)
		return nil, errors.New("user already exists")
	}

	// Create user
	user := repo.User{
		ID:           uuid.New().String(),
		TenantID:     tenantID,
		Email:        email,
		PasswordHash: password, // Will be hashed in repo
		FirstName:    &firstName,
		LastName:     &lastName,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("ERROR: user_service.CreateUser userRepo.Create: %v", err)
		return nil, err
	}

	// Assign roles
	if len(roleIDs) > 0 {
		log.Printf("DEBUG: user_service.CreateUser assigning roles=%v", roleIDs)
		err = s.userRepo.SetUserRoles(ctx, user.ID, roleIDs)
		if err != nil {
			log.Printf("ERROR: user_service.CreateUser SetUserRoles: %v", err)
			return nil, err
		}
	}

	log.Printf("DEBUG: user_service.CreateUser success userID=%s", user.ID)

	return &user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*repo.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, tenantID string) ([]repo.User, error) {
	return s.userRepo.List(ctx, tenantID)
}

func (s *UserService) ListUsersPaginated(ctx context.Context, tenantID string, page, pageSize int) ([]repo.User, int64, error) {
	return s.userRepo.ListPaginated(ctx, tenantID, page, pageSize)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, firstName, lastName *string, isActive *bool, roleIDs []string) error {
	log.Printf("DEBUG: user_service.UpdateUser id=%s", id)
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		log.Printf("ERROR: user_service.UpdateUser GetByID: %v", err)
		return err
	}
	if user == nil {
		log.Printf("WARN: user_service.UpdateUser not found id=%s", id)
		return errors.New("user not found")
	}

	// Update user fields
	if firstName != nil {
		user.FirstName = firstName
	}
	if lastName != nil {
		user.LastName = lastName
	}
	if isActive != nil {
		user.IsActive = *isActive
	}

	err = s.userRepo.Update(ctx, *user)
	if err != nil {
		log.Printf("ERROR: user_service.UpdateUser repo.Update: %v", err)
		return err
	}

	// Update roles if provided
	if roleIDs != nil {
		err = s.userRepo.SetUserRoles(ctx, id, roleIDs)
		if err != nil {
			log.Printf("ERROR: user_service.UpdateUser SetUserRoles: %v", err)
			return err
		}
	}

	return nil
}

func (s *UserService) CreateRole(ctx context.Context, tenantID, name, description string, permissionIDs []string) (*repo.Role, error) {
	role, err := s.roleRepo.Create(ctx, tenantID, name, description)
	if err != nil {
		return nil, err
	}

	// Assign permissions
	if len(permissionIDs) > 0 {
		err = s.roleRepo.SetRolePermissions(ctx, role.ID, permissionIDs)
		if err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (s *UserService) ListRoles(ctx context.Context, tenantID string) ([]repo.Role, error) {
	return s.roleRepo.List(ctx, tenantID)
}

func (s *UserService) GetPermissions(ctx context.Context) ([]repo.Permission, error) {
	return s.roleRepo.GetPermissions(ctx, "")
}

func (s *UserService) GetUserWithRoles(ctx context.Context, userID string) (*repo.UserWithRoles, error) {
	return s.userRepo.GetUserWithRoles(ctx, userID)
}

func (s *UserService) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	return s.userRepo.GetUserPermissions(ctx, userID)
}

func (s *UserService) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	permissions, err := s.userRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm == permission {
			return true, nil
		}
	}
	return false, nil
}

func (s *UserService) UpdateRole(ctx context.Context, id, name, description *string, permissionIDs []string) error {
	if id == nil {
		return errors.New("role id is required")
	}
	role, err := s.roleRepo.GetByID(ctx, *id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Compute updated fields safely
	updatedName := role.Name
	if name != nil {
		updatedName = *name
	}
	var updatedDescription string
	if role.Description != nil {
		updatedDescription = *role.Description
	}
	if description != nil {
		updatedDescription = *description
	}

	err = s.roleRepo.Update(ctx, *id, updatedName, updatedDescription)
	if err != nil {
		return err
	}

	// Update permissions if provided (non-nil slice means caller intends update)
	if permissionIDs != nil {
		err = s.roleRepo.SetRolePermissions(ctx, *id, permissionIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) SearchUsers(ctx context.Context, tenantID string, search, role string, isActive *bool, sortBy, sortDir string, page, pageSize int) ([]repo.User, int64, error) {
	return s.userRepo.SearchUsers(ctx, tenantID, search, role, isActive, sortBy, sortDir, page, pageSize)
}

func (s *UserService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return s.userRepo.GetUserRoles(ctx, userID)
}

func (s *UserService) GetUserDetail(ctx context.Context, userID string) (*repo.User, []string, map[string]int, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	if user == nil {
		return nil, nil, nil, errors.New("user not found")
	}

	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		// Если не удается получить роли, возвращаем пустой массив
		roles = []string{}
	}

	stats, err := s.userRepo.GetUserStats(ctx, userID)
	if err != nil {
		// Если не удается получить статистику, возвращаем нули
		stats = map[string]int{
			"documents_count": 0,
			"risks_count":     0,
			"incidents_count": 0,
			"assets_count":    0,
		}
	}

	return user, roles, stats, nil
}
