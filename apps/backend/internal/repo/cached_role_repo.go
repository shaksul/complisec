package repo

import (
	"context"
	"time"

	"risknexus/backend/internal/cache"
)

// CachedRoleRepo репозиторий ролей с кэшированием
type CachedRoleRepo struct {
	roleRepo *RoleRepo
	cache    cache.Cache
	ttl      time.Duration
}

// NewCachedRoleRepo создает новый кэшированный репозиторий ролей
func NewCachedRoleRepo(roleRepo *RoleRepo, cache cache.Cache, ttl time.Duration) *CachedRoleRepo {
	return &CachedRoleRepo{
		roleRepo: roleRepo,
		cache:    cache,
		ttl:      ttl,
	}
}

// GetByID получает роль по ID с кэшированием
func (r *CachedRoleRepo) GetByID(ctx context.Context, id string) (*Role, error) {
	key := cache.GenerateKey("role", id)

	// Пытаемся получить из кэша
	if cached, found := r.cache.Get(ctx, key); found {
		if role, ok := cached.(*Role); ok {
			return role, nil
		}
	}

	// Получаем из базы данных
	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if role != nil {
		r.cache.Set(ctx, key, role, r.ttl)
	}

	return role, nil
}

// List получает список ролей с кэшированием
func (r *CachedRoleRepo) List(ctx context.Context, tenantID string) ([]Role, error) {
	key := cache.GenerateKey("roles", tenantID)

	// Пытаемся получить из кэша
	if cached, found := r.cache.Get(ctx, key); found {
		if roles, ok := cached.([]Role); ok {
			return roles, nil
		}
	}

	// Получаем из базы данных
	roles, err := r.roleRepo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.Set(ctx, key, roles, r.ttl)

	return roles, nil
}

// Create создает роль и инвалидирует кэш
func (r *CachedRoleRepo) Create(ctx context.Context, tenantID, name, description string) (*Role, error) {
	role, err := r.roleRepo.Create(ctx, tenantID, name, description)
	if err != nil {
		return nil, err
	}

	// Инвалидируем кэш
	r.invalidateRoleCache(ctx, tenantID)

	return role, nil
}

// GetByName получает роль по имени
func (r *CachedRoleRepo) GetByName(ctx context.Context, tenantID, name string) (*Role, error) {
	return r.roleRepo.GetByName(ctx, tenantID, name)
}

// Update обновляет роль и инвалидирует кэш
func (r *CachedRoleRepo) Update(ctx context.Context, id, name, description string) error {
	err := r.roleRepo.Update(ctx, id, name, description)
	if err != nil {
		return err
	}

	// Получаем роль для инвалидации кэша
	role, err := r.roleRepo.GetByID(ctx, id)
	if err == nil && role != nil {
		r.invalidateRoleCache(ctx, role.TenantID)
		r.cache.Delete(ctx, cache.GenerateKey("role", role.ID))
		r.cache.Delete(ctx, cache.GenerateKey("role_with_permissions", role.ID))
	}

	return nil
}

// Delete удаляет роль и инвалидирует кэш
func (r *CachedRoleRepo) Delete(ctx context.Context, id string) error {
	// Получаем роль для получения tenantID
	role, err := r.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.roleRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Инвалидируем кэш
	if role != nil {
		r.invalidateRoleCache(ctx, role.TenantID)
		r.cache.Delete(ctx, cache.GenerateKey("role", id))
	}

	return nil
}

// GetPermissions получает права с кэшированием
func (r *CachedRoleRepo) GetPermissions(ctx context.Context, tenantID string) ([]Permission, error) {
	key := cache.GenerateKey("permissions", tenantID)

	// Пытаемся получить из кэша
	if cached, found := r.cache.Get(ctx, key); found {
		if permissions, ok := cached.([]Permission); ok {
			return permissions, nil
		}
	}

	// Получаем из базы данных
	permissions, err := r.roleRepo.GetPermissions(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.Set(ctx, key, permissions, r.ttl)

	return permissions, nil
}

// GetRolePermissions получает права роли с кэшированием
func (r *CachedRoleRepo) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	key := cache.GenerateKey("role_permissions", roleID)

	// Пытаемся получить из кэша
	if cached, found := r.cache.Get(ctx, key); found {
		if permissions, ok := cached.([]string); ok {
			return permissions, nil
		}
	}

	// Получаем из базы данных
	permissions, err := r.roleRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.Set(ctx, key, permissions, r.ttl)

	return permissions, nil
}

// SetRolePermissions устанавливает права роли и инвалидирует кэш
func (r *CachedRoleRepo) SetRolePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	err := r.roleRepo.SetRolePermissions(ctx, roleID, permissionIDs)
	if err != nil {
		return err
	}

	// Инвалидируем кэш прав роли
	r.cache.Delete(ctx, cache.GenerateKey("role_permissions", roleID))
	// Инвалидируем кэш роли с правами
	r.cache.Delete(ctx, cache.GenerateKey("role_with_permissions", roleID))

	return nil
}

// GetRoleWithPermissions получает роль с правами с кэшированием
func (r *CachedRoleRepo) GetRoleWithPermissions(ctx context.Context, roleID string) (*RoleWithPermissions, error) {
	key := cache.GenerateKey("role_with_permissions", roleID)

	// Пытаемся получить из кэша
	if cached, found := r.cache.Get(ctx, key); found {
		if roleWithPerms, ok := cached.(*RoleWithPermissions); ok {
			return roleWithPerms, nil
		}
	}

	// Получаем из базы данных
	roleWithPerms, err := r.roleRepo.GetRoleWithPermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	if roleWithPerms != nil {
		r.cache.Set(ctx, key, roleWithPerms, r.ttl)
	}

	return roleWithPerms, nil
}

// GetUsersByRole получает пользователей роли с кэшированием
func (r *CachedRoleRepo) GetUsersByRole(ctx context.Context, roleID string) ([]User, error) {
	key := cache.GenerateKey("role_users", roleID)

	// Пытаемся получить из кэша
	if cached, found := r.cache.Get(ctx, key); found {
		if users, ok := cached.([]User); ok {
			return users, nil
		}
	}

	// Получаем из базы данных
	users, err := r.roleRepo.GetUsersByRole(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.Set(ctx, key, users, r.ttl)

	return users, nil
}

// invalidateRoleCache инвалидирует кэш ролей для тенанта
func (r *CachedRoleRepo) invalidateRoleCache(ctx context.Context, tenantID string) {
	r.cache.Delete(ctx, cache.GenerateKey("roles", tenantID))
	r.cache.Delete(ctx, cache.GenerateKey("permissions", tenantID))
}
