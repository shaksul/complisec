package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"risknexus/backend/internal/repo"
)

// RoleRepository интерфейс для работы с ролями
type RoleRepository interface {
	Create(ctx context.Context, tenantID, name, description string) (*repo.Role, error)
	GetByID(ctx context.Context, id string) (*repo.Role, error)
	GetByName(ctx context.Context, tenantID, name string) (*repo.Role, error)
	List(ctx context.Context, tenantID string) ([]repo.Role, error)
	Update(ctx context.Context, id, name, description string) error
	Delete(ctx context.Context, id string) error
	SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]string, error)
	GetRoleWithPermissions(ctx context.Context, roleID string) (*repo.RoleWithPermissions, error)
	GetUsersByRole(ctx context.Context, roleID string) ([]repo.User, error)
	GetPermissions(ctx context.Context, tenantID string) ([]repo.Permission, error)
}

type RoleService struct {
	roleRepo  RoleRepository
	userRepo  *repo.UserRepo
	auditRepo *repo.AuditRepo
	cache     Cache
	cacheKey  CacheKey
	config    CacheConfig
}

func NewRoleService(roleRepo RoleRepository, userRepo *repo.UserRepo, auditRepo *repo.AuditRepo) *RoleService {
	return &RoleService{
		roleRepo:  roleRepo,
		userRepo:  userRepo,
		auditRepo: auditRepo,
		cacheKey:  CacheKey{},
		config:    DefaultCacheConfig(),
	}
}

// NewRoleServiceWithCache создает RoleService с кэшем
func NewRoleServiceWithCache(roleRepo RoleRepository, userRepo *repo.UserRepo, auditRepo *repo.AuditRepo, cache Cache) *RoleService {
	return &RoleService{
		roleRepo:  roleRepo,
		userRepo:  userRepo,
		auditRepo: auditRepo,
		cache:     cache,
		cacheKey:  CacheKey{},
		config:    DefaultCacheConfig(),
	}
}

// validateRoleName проверяет корректность имени роли
func (s *RoleService) validateRoleName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return NewValidationError("name", "role name cannot be empty")
	}
	if len(trimmed) > 100 {
		return NewValidationError("name", "role name cannot exceed 100 characters")
	}
	if len(trimmed) < 1 {
		return NewValidationError("name", "role name must be at least 1 character")
	}
	return nil
}

// validateDescription проверяет корректность описания роли
func (s *RoleService) validateDescription(description string) error {
	if len(description) > 500 {
		return NewValidationError("description", "role description cannot exceed 500 characters")
	}
	return nil
}

// validatePermissionIDs проверяет корректность ID прав
func (s *RoleService) validatePermissionIDs(ctx context.Context, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return nil // Пустой список разрешен
	}

	// Получаем все существующие права
	allPermissions, err := s.GetPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get permissions: %w", err)
	}

	// Создаем мапу для быстрого поиска
	permissionMap := make(map[string]bool)
	for _, perm := range allPermissions {
		permissionMap[perm.ID] = true
	}

	// Проверяем каждый ID
	for _, permID := range permissionIDs {
		if permID == "" {
			return NewValidationError("permission_ids", "permission ID cannot be empty")
		}
		if !permissionMap[permID] {
			return NewBusinessError("PERMISSION_NOT_FOUND",
				fmt.Sprintf("permission with ID %s does not exist", permID),
				map[string]interface{}{"permission_id": permID})
		}
	}

	return nil
}

// CreateRole создает новую роль с правами
func (s *RoleService) CreateRole(ctx context.Context, tenantID, name, description string, permissionIDs []string) (*repo.Role, error) {
	// Валидация входных данных
	if err := s.validateRoleName(name); err != nil {
		return nil, err
	}
	if err := s.validateDescription(description); err != nil {
		return nil, err
	}
	if err := s.validatePermissionIDs(ctx, permissionIDs); err != nil {
		return nil, err
	}

	// Проверяем, что роль с таким именем не существует в тенанте
	existingRole, err := s.roleRepo.GetByName(ctx, tenantID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole != nil {
		return nil, ErrRoleAlreadyExists
	}

	// Создаем роль в транзакции
	role, err := s.roleRepo.Create(ctx, tenantID, name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Назначаем права роли (уже в транзакции в SetRolePermissions)
	if len(permissionIDs) > 0 {
		err = s.roleRepo.SetRolePermissions(ctx, role.ID, permissionIDs)
		if err != nil {
			// В случае ошибки пытаемся удалить созданную роль
			if deleteErr := s.roleRepo.Delete(ctx, role.ID); deleteErr != nil {
				return nil, fmt.Errorf("failed to set permissions and cleanup role: %w (cleanup error: %v)", err, deleteErr)
			}
			return nil, fmt.Errorf("failed to set role permissions: %w", err)
		}
	}

	// Инвалидируем кэш
	s.invalidateRoleCache(ctx, tenantID, role.ID)

	// Логируем создание роли
	if s.auditRepo != nil {
		auditData := map[string]interface{}{
			"role_name":        name,
			"description":      description,
			"permission_count": len(permissionIDs),
		}
		s.auditRepo.LogAction(ctx, tenantID, "system", "role.create", "role", &role.ID, auditData)
	}

	return role, nil
}

// GetRole получает роль по ID
func (s *RoleService) GetRole(ctx context.Context, roleID string) (*repo.Role, error) {
	// Если кэш не настроен, используем прямой запрос
	if s.cache == nil {
		return s.roleRepo.GetByID(ctx, roleID)
	}

	// Пытаемся получить из кэша
	key := s.cacheKey.Role("", roleID) // tenantID будет получен из роли
	if cached, err := s.cache.Get(ctx, key); err == nil {
		if role, ok := cached.(*repo.Role); ok {
			return role, nil
		}
	}

	// Получаем из базы данных
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Сохраняем в кэш
	cacheKey := s.cacheKey.Role(role.TenantID, roleID)
	s.cache.Set(ctx, cacheKey, role, s.config.RoleTTL)

	return role, nil
}

// GetRoleWithPermissions получает роль с правами
func (s *RoleService) GetRoleWithPermissions(ctx context.Context, roleID string) (*repo.RoleWithPermissions, error) {
	return s.roleRepo.GetRoleWithPermissions(ctx, roleID)
}

// ListRoles получает список ролей тенанта
func (s *RoleService) ListRoles(ctx context.Context, tenantID string) ([]repo.Role, error) {
	return s.roleRepo.List(ctx, tenantID)
}

// UpdateRole обновляет роль
func (s *RoleService) UpdateRole(ctx context.Context, roleID string, name, description *string, permissionIDs []string) error {
	// Получаем существующую роль
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// Валидация входных данных
	if name != nil {
		if err := s.validateRoleName(*name); err != nil {
			return err
		}
		// Проверяем уникальность имени в рамках тенанта
		existingRole, err := s.roleRepo.GetByName(ctx, role.TenantID, *name)
		if err != nil {
			return fmt.Errorf("failed to check existing role name: %w", err)
		}
		if existingRole != nil && existingRole.ID != roleID {
			return errors.New("role with this name already exists")
		}
	}

	if description != nil {
		if err := s.validateDescription(*description); err != nil {
			return err
		}
	}

	if permissionIDs != nil {
		if err := s.validatePermissionIDs(ctx, permissionIDs); err != nil {
			return err
		}
	}

	// Обновляем поля роли
	updateName := role.Name
	updateDescription := ""
	if role.Description != nil {
		updateDescription = *role.Description
	}

	if name != nil {
		updateName = *name
	}
	if description != nil {
		updateDescription = *description
	}

	err = s.roleRepo.Update(ctx, roleID, updateName, updateDescription)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	// Обновляем права если предоставлены
	if permissionIDs != nil {
		err = s.roleRepo.SetRolePermissions(ctx, roleID, permissionIDs)
		if err != nil {
			return fmt.Errorf("failed to update role permissions: %w", err)
		}
	}

	// Инвалидируем кэш
	s.invalidateRoleCache(ctx, role.TenantID, roleID)

	// Логируем обновление роли
	if s.auditRepo != nil {
		auditData := map[string]interface{}{
			"role_name":        updateName,
			"description":      updateDescription,
			"permission_count": len(permissionIDs),
		}
		s.auditRepo.LogAction(ctx, role.TenantID, "system", "role.update", "role", &roleID, auditData)
	}

	return nil
}

// DeleteRole удаляет роль
func (s *RoleService) DeleteRole(ctx context.Context, roleID string) error {
	// Проверяем, что роль не используется пользователями
	users, err := s.roleRepo.GetUsersByRole(ctx, roleID)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return ErrRoleInUse
	}

	// Получаем роль для аудита
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	err = s.roleRepo.Delete(ctx, roleID)
	if err != nil {
		return err
	}

	// Инвалидируем кэш
	if role != nil {
		s.invalidateRoleCache(ctx, role.TenantID, roleID)
	}

	// Логируем удаление роли
	if s.auditRepo != nil && role != nil {
		auditData := map[string]interface{}{
			"role_name":   role.Name,
			"description": role.Description,
		}
		s.auditRepo.LogAction(ctx, role.TenantID, "system", "role.delete", "role", &roleID, auditData)
	}

	return nil
}

// GetPermissions получает все права в системе
func (s *RoleService) GetPermissions(ctx context.Context) ([]repo.Permission, error) {
	return s.roleRepo.GetPermissions(ctx, "")
}

// GetRolePermissions получает права роли
func (s *RoleService) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	return s.roleRepo.GetRolePermissions(ctx, roleID)
}

// SetRolePermissions устанавливает права роли
func (s *RoleService) SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return s.roleRepo.SetRolePermissions(ctx, roleID, permissionIDs)
}

// GetUsersByRole получает пользователей с определенной ролью
func (s *RoleService) GetUsersByRole(ctx context.Context, roleID string) ([]repo.User, error) {
	return s.roleRepo.GetUsersByRole(ctx, roleID)
}

// AssignRoleToUser назначает роль пользователю
func (s *RoleService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// Проверяем, что пользователь существует
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Проверяем, что роль существует
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrRoleNotFound
	}

	// Получаем текущие роли пользователя
	currentRoles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return err
	}

	// Проверяем, что роль еще не назначена
	for _, r := range currentRoles {
		if r == role.Name {
			return errors.New("user already has this role")
		}
	}

	// Получаем ID ролей для обновления
	allRoles, err := s.roleRepo.List(ctx, user.TenantID)
	if err != nil {
		return err
	}

	var roleIDs []string
	for _, r := range allRoles {
		for _, currentRole := range currentRoles {
			if r.Name == currentRole {
				roleIDs = append(roleIDs, r.ID)
				break
			}
		}
	}
	roleIDs = append(roleIDs, roleID)

	return s.userRepo.SetUserRoles(ctx, userID, roleIDs)
}

// RemoveRoleFromUser убирает роль у пользователя
func (s *RoleService) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	// Получаем текущие роли пользователя
	currentRoles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return err
	}

	// Получаем пользователя для получения tenantID
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	allRoles, err := s.roleRepo.List(ctx, user.TenantID)
	if err != nil {
		return err
	}

	// Находим роль для удаления
	var roleToRemove *repo.Role
	for _, role := range allRoles {
		if role.ID == roleID {
			roleToRemove = &role
			break
		}
	}

	if roleToRemove == nil {
		return errors.New("role not found")
	}

	// Проверяем, что роль назначена пользователю
	hasRole := false
	for _, currentRole := range currentRoles {
		if currentRole == roleToRemove.Name {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return errors.New("user does not have this role")
	}

	// Формируем новый список ролей без удаляемой
	var newRoleIDs []string
	for _, role := range allRoles {
		shouldKeep := false
		for _, currentRole := range currentRoles {
			if role.Name == currentRole && role.ID != roleID {
				shouldKeep = true
				break
			}
		}
		if shouldKeep {
			newRoleIDs = append(newRoleIDs, role.ID)
		}
	}

	return s.userRepo.SetUserRoles(ctx, userID, newRoleIDs)
}

// invalidateRoleCache инвалидирует кэш роли
func (s *RoleService) invalidateRoleCache(ctx context.Context, tenantID, roleID string) {
	if s.cache == nil {
		return
	}

	// Удаляем роль из кэша
	s.cache.Delete(ctx, s.cacheKey.Role(tenantID, roleID))

	// Удаляем список ролей тенанта
	s.cache.Delete(ctx, s.cacheKey.RoleList(tenantID))

	// Удаляем права роли
	s.cache.Delete(ctx, s.cacheKey.RolePermissions(roleID))
}

// invalidateUserRoleCache инвалидирует кэш ролей пользователя
func (s *RoleService) invalidateUserRoleCache(ctx context.Context, userID string) {
	if s.cache == nil {
		return
	}

	s.cache.Delete(ctx, s.cacheKey.UserRoles(userID))
}
