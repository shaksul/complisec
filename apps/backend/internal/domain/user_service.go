package domain

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"risknexus/backend/internal/repo"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo  *repo.UserRepo
	roleRepo  *repo.RoleRepo
	assetRepo *repo.AssetRepo
}

func NewUserService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, assetRepo *repo.AssetRepo) *UserService {
	return &UserService{
		userRepo:  userRepo,
		roleRepo:  roleRepo,
		assetRepo: assetRepo,
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

func (s *UserService) GetUserByTenant(ctx context.Context, id, tenantID string) (*repo.User, error) {
	return s.userRepo.GetByIDAndTenant(ctx, id, tenantID)
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

func (s *UserService) UpdateUserByTenant(ctx context.Context, id, tenantID string, firstName, lastName *string, password *string, isActive *bool, roleIDs []string) error {
	log.Printf("DEBUG: user_service.UpdateUserByTenant id=%s tenant=%s", id, tenantID)
	user, err := s.userRepo.GetByIDAndTenant(ctx, id, tenantID)
	if err != nil {
		log.Printf("ERROR: user_service.UpdateUserByTenant GetByIDAndTenant: %v", err)
		return err
	}
	if user == nil {
		log.Printf("WARN: user_service.UpdateUserByTenant not found id=%s tenant=%s", id, tenantID)
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
	
	// Update password if provided (will be hashed in repo.Update)
	if password != nil && *password != "" {
		log.Printf("DEBUG: user_service.UpdateUserByTenant updating password for user=%s", id)
		user.PasswordHash = *password // Repo will hash it
	}

	err = s.userRepo.Update(ctx, *user)
	if err != nil {
		log.Printf("ERROR: user_service.UpdateUserByTenant repo.Update: %v", err)
		return err
	}

	// Update roles if provided
	if roleIDs != nil {
		log.Printf("DEBUG: user_service.UpdateUserByTenant roleIDs=%v", roleIDs)
		// Convert mixed role data (names and IDs) to role IDs
		var actualRoleIDs []string
		roleIDSet := make(map[string]bool) // For deduplication

		for _, roleData := range roleIDs {
			var roleID string

			// Check if it's a UUID (role ID) or a name
			if len(roleData) == 36 && strings.Contains(roleData, "-") {
				// It's a UUID, use as is
				roleID = roleData
			} else {
				// It's a role name, convert to ID
				role, err := s.roleRepo.GetByName(ctx, tenantID, roleData)
				if err != nil {
					log.Printf("ERROR: user_service.UpdateUserByTenant GetByName: %v", err)
					return err
				}
				if role != nil {
					roleID = role.ID
				}
			}

			// Add to set if not already present (deduplication)
			if roleID != "" && !roleIDSet[roleID] {
				actualRoleIDs = append(actualRoleIDs, roleID)
				roleIDSet[roleID] = true
			}
		}

		err = s.userRepo.SetUserRoles(ctx, id, actualRoleIDs)
		if err != nil {
			log.Printf("ERROR: user_service.UpdateUserByTenant SetUserRoles: %v", err)
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

func (s *UserService) GetUserWithRolesByTenant(ctx context.Context, userID, tenantID string) (*repo.UserWithRoles, error) {
	return s.userRepo.GetUserWithRolesByTenant(ctx, userID, tenantID)
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

	var updatedDescription *string
	if description != nil {
		updatedDescription = description
	} else {
		// Сохраняем существующее описание (может быть NULL)
		updatedDescription = role.Description
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

func (s *UserService) GetUserDetailByTenant(ctx context.Context, userID, tenantID string) (*repo.User, []string, map[string]int, error) {
	user, err := s.userRepo.GetByIDAndTenant(ctx, userID, tenantID)
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
			"sessions_count":  0,
			"login_count":     0,
			"activity_score":  0,
		}
	}

	return user, roles, stats, nil
}

// GetUserActivity retrieves user activity with pagination
func (s *UserService) GetUserActivity(ctx context.Context, userID string, page, pageSize int) ([]repo.UserActivity, int64, error) {
	return s.userRepo.GetUserActivity(ctx, userID, page, pageSize)
}

// GetUserActivityStats retrieves user activity statistics
func (s *UserService) GetUserActivityStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	return s.userRepo.GetUserActivityStats(ctx, userID)
}

// LogUserActivity logs a user activity
func (s *UserService) LogUserActivity(ctx context.Context, userID, action, description, ipAddress, userAgent string, metadata map[string]interface{}) error {
	return s.userRepo.LogUserActivity(ctx, userID, action, description, ipAddress, userAgent, metadata)
}

// LogLoginAttempt logs a login attempt
func (s *UserService) LogLoginAttempt(ctx context.Context, userID, tenantID, email, ipAddress, userAgent string, success bool, failureReason string) error {
	return s.userRepo.LogLoginAttempt(ctx, userID, tenantID, email, ipAddress, userAgent, success, failureReason)
}

// GetUserResponsibleAssets retrieves assets for which the user is responsible
func (s *UserService) GetUserResponsibleAssets(ctx context.Context, userID, tenantID string) ([]repo.Asset, error) {
	return s.assetRepo.GetUserResponsibleAssets(ctx, tenantID, userID)
}
