package repo

import (
	"context"
	"log"
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
	log.Printf("DEBUG: cached_role_repo.Create invalidate cache tenant=%s", tenantID)
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
		log.Printf("DEBUG: cached_role_repo.Update invalidate cache tenant=%s role=%s", role.TenantID, role.ID)
		r.invalidateRoleCache(ctx, role.TenantID)

		roleKey := cache.GenerateKey("role", role.ID)
		roleWithPermsKey := cache.GenerateKey("role_with_permissions", role.ID)
		log.Printf("DEBUG: cached_role_repo.Update delete keys role=%s role_with_permissions=%s", roleKey, roleWithPermsKey)
		r.cache.Delete(ctx, roleKey)
		r.cache.Delete(ctx, roleWithPermsKey)
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
		log.Printf("DEBUG: cached_role_repo.Delete invalidate cache tenant=%s role=%s", role.TenantID, role.ID)
		r.invalidateRoleCache(ctx, role.TenantID)

		roleKey := cache.GenerateKey("role", id)
		log.Printf("DEBUG: cached_role_repo.Delete delete key role=%s", roleKey)
		r.cache.Delete(ctx, roleKey)
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

	role, getErr := r.roleRepo.GetByID(ctx, roleID)
	if getErr != nil {
		log.Printf("ERROR: cached_role_repo.SetRolePermissions failed to fetch role for cache invalidation role=%s err=%v", roleID, getErr)
	}

	// Инвалидируем кэш прав роли
	rolePermsKey := cache.GenerateKey("role_permissions", roleID)
	roleWithPermsKey := cache.GenerateKey("role_with_permissions", roleID)
	log.Printf("DEBUG: cached_role_repo.SetRolePermissions delete keys role_permissions=%s role_with_permissions=%s", rolePermsKey, roleWithPermsKey)
	r.cache.Delete(ctx, rolePermsKey)
	// Инвалидируем кэш роли с правами
	r.cache.Delete(ctx, roleWithPermsKey)

	if role != nil {
		log.Printf("DEBUG: cached_role_repo.SetRolePermissions invalidate tenant caches tenant=%s role=%s", role.TenantID, role.ID)
		r.invalidateRoleCache(ctx, role.TenantID)

		roleKey := cache.GenerateKey("role", role.ID)
		log.Printf("DEBUG: cached_role_repo.SetRolePermissions delete key role=%s", roleKey)
		r.cache.Delete(ctx, roleKey)
	}

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
	rolesKey := cache.GenerateKey("roles", tenantID)
	permissionsKey := cache.GenerateKey("permissions", tenantID)
	log.Printf("DEBUG: cached_role_repo.invalidateRoleCache delete keys roles=%s permissions=%s", rolesKey, permissionsKey)
	r.cache.Delete(ctx, rolesKey)
	r.cache.Delete(ctx, permissionsKey)
}
