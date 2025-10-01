package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// CacheItem элемент кэша
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// MemoryCache простое in-memory кэширование
type MemoryCache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

// NewMemoryCache создает новый экземпляр кэша
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	// Запускаем горутину для очистки просроченных элементов
	go cache.cleanup()

	return cache
}

// Set сохраняет значение в кэш
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Get получает значение из кэша
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mutex.RLock()
	item, exists := c.items[key]
	c.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	// Проверяем, не истек ли срок действия
	if time.Now().After(item.ExpiresAt) {
		// Удаляем просроченный элемент
		c.mutex.Lock()
		delete(c.items, key)
		c.mutex.Unlock()
		return nil, false
	}

	return item.Value, true
}

// Delete удаляет значение из кэша
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
	return nil
}

// Clear очищает весь кэш
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items = make(map[string]CacheItem)
	return nil
}

// cleanup периодически очищает просроченные элементы
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// Cache интерфейс для кэширования
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, bool)
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// CachedService сервис с кэшированием
type CachedService struct {
	cache Cache
	ttl   time.Duration
}

// NewCachedService создает новый сервис с кэшированием
func NewCachedService(cache Cache, ttl time.Duration) *CachedService {
	return &CachedService{
		cache: cache,
		ttl:   ttl,
	}
}

// GetOrSet получает значение из кэша или вычисляет и сохраняет
func (s *CachedService) GetOrSet(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
	// Пытаемся получить из кэша
	if value, found := s.cache.Get(ctx, key); found {
		return value, nil
	}

	// Вычисляем значение
	value, err := fn()
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	s.cache.Set(ctx, key, value, s.ttl)

	return value, nil
}

// Invalidate удаляет ключ из кэша
func (s *CachedService) Invalidate(ctx context.Context, key string) error {
	return s.cache.Delete(ctx, key)
}

// InvalidatePattern удаляет ключи по паттерну
func (s *CachedService) InvalidatePattern(ctx context.Context, pattern string) error {
	// Для простоты очищаем весь кэш
	// В реальной системе можно использовать более сложную логику
	return s.cache.Clear(ctx)
}

// GenerateKey генерирует ключ кэша
func GenerateKey(prefix string, parts ...interface{}) string {
	key := prefix
	for _, part := range parts {
		key += fmt.Sprintf(":%v", part)
	}
	return key
}

// Serialize сериализует значение для кэша
func Serialize(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Deserialize десериализует значение из кэша
func Deserialize(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
