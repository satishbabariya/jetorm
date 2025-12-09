package core

import (
	"context"
	"fmt"
	"time"
)

// Cache interface for caching repository results
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) (interface{}, bool)
	
	// Set stores a value in cache
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	
	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error
	
	// Clear clears all cache entries
	Clear(ctx context.Context) error
}

// CacheKeyGenerator generates cache keys for entities
type CacheKeyGenerator[T any, ID comparable] struct {
	entityType string
}

// NewCacheKeyGenerator creates a new cache key generator
func NewCacheKeyGenerator[T any, ID comparable](entityType string) *CacheKeyGenerator[T, ID] {
	return &CacheKeyGenerator[T, ID]{
		entityType: entityType,
	}
}

// KeyForID generates a cache key for an entity ID
func (ckg *CacheKeyGenerator[T, ID]) KeyForID(id ID) string {
	return fmt.Sprintf("%s:id:%v", ckg.entityType, id)
}

// KeyForQuery generates a cache key for a query
func (ckg *CacheKeyGenerator[T, ID]) KeyForQuery(query string, args ...interface{}) string {
	key := fmt.Sprintf("%s:query:%s", ckg.entityType, query)
	for _, arg := range args {
		key += fmt.Sprintf(":%v", arg)
	}
	return key
}

// CachedRepository wraps a repository with caching
type CachedRepository[T any, ID comparable] struct {
	repo  Repository[T, ID]
	cache Cache
	ttl   time.Duration
	keyGen *CacheKeyGenerator[T, ID]
}

// NewCachedRepository creates a new cached repository
func NewCachedRepository[T any, ID comparable](
	repo Repository[T, ID],
	cache Cache,
	entityType string,
	ttl time.Duration,
) *CachedRepository[T, ID] {
	return &CachedRepository[T, ID]{
		repo:   repo,
		cache:  cache,
		ttl:    ttl,
		keyGen: NewCacheKeyGenerator[T, ID](entityType),
	}
}

// FindByID implements Repository.FindByID with caching
func (cr *CachedRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	key := cr.keyGen.KeyForID(id)
	
	// Try cache first
	if cached, ok := cr.cache.Get(ctx, key); ok {
		if entity, ok := cached.(*T); ok {
			return entity, nil
		}
	}
	
	// Cache miss - load from repository
	entity, err := cr.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	if entity != nil {
		cr.cache.Set(ctx, key, entity, cr.ttl)
	}
	
	return entity, nil
}

// Save implements Repository.Save with cache invalidation
func (cr *CachedRepository[T, ID]) Save(ctx context.Context, entity *T) (*T, error) {
	saved, err := cr.repo.Save(ctx, entity)
	if err != nil {
		return nil, err
	}
	
	// Invalidate cache for this entity
	// Note: Would need to extract ID from entity
	// This is a simplified version
	cr.cache.Clear(ctx) // Clear all for simplicity
	
	return saved, nil
}

// Delete implements Repository.Delete with cache invalidation
func (cr *CachedRepository[T, ID]) Delete(ctx context.Context, entity *T) error {
	err := cr.repo.Delete(ctx, entity)
	if err != nil {
		return err
	}
	
	// Invalidate cache
	cr.cache.Clear(ctx)
	
	return nil
}

// InMemoryCache is a simple in-memory cache implementation
type InMemoryCache struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheEntry),
	}
}

// Get retrieves a value from cache
func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}
	
	// Check expiration
	if time.Now().After(entry.expiresAt) {
		delete(c.data, key)
		return nil, false
	}
	
	return entry.value, true
}

// Set stores a value in cache
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

// Delete removes a value from cache
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

// Clear clears all cache entries
func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.data = make(map[string]cacheEntry)
	return nil
}

