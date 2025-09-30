package domain

import (
	"context"
	"time"
)

// Cache интерфейс для кэширования
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// CacheKey генерирует ключи для кэша
type CacheKey struct{}

func (CacheKey) Role(tenantID, roleID string) string {
	return "role:" + tenantID + ":" + roleID
}

func (CacheKey) RoleList(tenantID string) string {
	return "roles:" + tenantID
}

func (CacheKey) Permissions() string {
	return "permissions:all"
}

func (CacheKey) RolePermissions(roleID string) string {
	return "role_permissions:" + roleID
}

func (CacheKey) UserRoles(userID string) string {
	return "user_roles:" + userID
}

// CacheConfig конфигурация кэша
type CacheConfig struct {
	DefaultTTL    time.Duration
	RoleTTL       time.Duration
	PermissionTTL time.Duration
	UserRoleTTL   time.Duration
}

// DefaultCacheConfig возвращает конфигурацию по умолчанию
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		DefaultTTL:    5 * time.Minute,
		RoleTTL:       10 * time.Minute,
		PermissionTTL: 30 * time.Minute,
		UserRoleTTL:   15 * time.Minute,
	}
}
