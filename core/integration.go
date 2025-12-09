package core

import (
	"context"
	"fmt"
	"time"

	"github.com/satishbabariya/jetorm/hooks"
)

// Integration utilities for combining features

// CachedRepositoryWithHooks combines caching and hooks
type CachedRepositoryWithHooks[T any, ID comparable] struct {
	*CachedRepository[T, ID]
	hooks *hooks.Hooks[T]
}

// NewCachedRepositoryWithHooks creates a repository with caching and hooks
func NewCachedRepositoryWithHooks[T any, ID comparable](
	repo Repository[T, ID],
	cache Cache,
	entityType string,
	ttl time.Duration,
	hooks *hooks.Hooks[T],
) *CachedRepositoryWithHooks[T, ID] {
	return &CachedRepositoryWithHooks[T, ID]{
		CachedRepository: NewCachedRepository(repo, cache, entityType, ttl),
		hooks:            hooks,
	}
}

// Save implements Repository.Save with caching and hooks
func (cr *CachedRepositoryWithHooks[T, ID]) Save(ctx context.Context, entity *T) (*T, error) {
	// Execute before hooks
	if cr.hooks != nil {
		// Check if entity is new (no ID) to determine create vs update
		id, err := ExtractID[T, ID](entity)
		if err != nil || IsZero(id) {
			if err := cr.hooks.ExecuteBeforeCreate(ctx, entity); err != nil {
				return nil, err
			}
		} else {
			if err := cr.hooks.ExecuteBeforeUpdate(ctx, entity); err != nil {
				return nil, err
			}
		}
	}

	// Save using cached repository
	saved, err := cr.CachedRepository.Save(ctx, entity)
	if err != nil {
		return nil, err
	}

	// Execute after hooks
	if cr.hooks != nil {
		id, _ := ExtractID[T, ID](saved)
		if IsZero(id) {
			if err := cr.hooks.ExecuteAfterCreate(ctx, saved); err != nil {
				return nil, err
			}
		} else {
			if err := cr.hooks.ExecuteAfterUpdate(ctx, saved); err != nil {
				return nil, err
			}
		}
	}

	return saved, nil
}

// RepositoryWithValidation wraps a repository with validation
type RepositoryWithValidation[T any, ID comparable] struct {
	repo      Repository[T, ID]
	validator *Validator
}

// NewRepositoryWithValidation creates a repository with validation
func NewRepositoryWithValidation[T any, ID comparable](
	repo Repository[T, ID],
	validator *Validator,
) *RepositoryWithValidation[T, ID] {
	return &RepositoryWithValidation[T, ID]{
		repo:      repo,
		validator: validator,
	}
}

// Save implements Repository.Save with validation
func (rv *RepositoryWithValidation[T, ID]) Save(ctx context.Context, entity *T) (*T, error) {
	// Validate entity
	if rv.validator != nil {
		if err := rv.validator.Validate(entity); err != nil {
			return nil, WrapError(err, "validation failed")
		}
	}

	return rv.repo.Save(ctx, entity)
}

// RepositoryWithMetrics wraps a repository with performance metrics
type RepositoryWithMetrics[T any, ID comparable] struct {
	repo    Repository[T, ID]
	profiler *QueryProfiler
}

// NewRepositoryWithMetrics creates a repository with metrics
func NewRepositoryWithMetrics[T any, ID comparable](
	repo Repository[T, ID],
	profiler *QueryProfiler,
) *RepositoryWithMetrics[T, ID] {
	return &RepositoryWithMetrics[T, ID]{
		repo:     repo,
		profiler: profiler,
	}
}

// FindByID implements Repository.FindByID with metrics
func (rm *RepositoryWithMetrics[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	query := fmt.Sprintf("FindByID(%v)", id)
	var result *T
	var err error

	if rm.profiler != nil {
		err = rm.profiler.Profile(ctx, query, func(ctx context.Context) error {
			result, err = rm.repo.FindByID(ctx, id)
			return err
		})
	} else {
		result, err = rm.repo.FindByID(ctx, id)
	}

	return result, err
}

// FullFeaturedRepository combines all features
type FullFeaturedRepository[T any, ID comparable] struct {
	repo          Repository[T, ID]
	cache         Cache
	hooks         *hooks.Hooks[T]
	validator     *Validator
	profiler      *QueryProfiler
	healthChecker *HealthChecker
	keyGen        *CacheKeyGenerator[T, ID]
	ttl           time.Duration
	entityType    string
}

// NewFullFeaturedRepository creates a repository with all features
func NewFullFeaturedRepository[T any, ID comparable](
	repo Repository[T, ID],
	cache Cache,
	entityType string,
	ttl time.Duration,
	hooks *hooks.Hooks[T],
	validator *Validator,
	profiler *QueryProfiler,
	db *Database,
) *FullFeaturedRepository[T, ID] {
	return &FullFeaturedRepository[T, ID]{
		repo:          repo,
		cache:         cache,
		validator:     validator,
		hooks:         hooks,
		profiler:      profiler,
		healthChecker: NewHealthChecker(db),
		keyGen:        NewCacheKeyGenerator[T, ID](entityType),
		ttl:           ttl,
		entityType:    entityType,
	}
}

// FindByID implements Repository.FindByID with all features
func (fr *FullFeaturedRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	// Try cache first
	if fr.cache != nil {
		key := fr.keyGen.KeyForID(id)
		if cached, ok := fr.cache.Get(ctx, key); ok {
			if entity, ok := cached.(*T); ok {
				return entity, nil
			}
		}
	}

	// Profile query if profiler available
	var result *T
	var err error
	if fr.profiler != nil {
		query := fmt.Sprintf("FindByID(%v)", id)
		err = fr.profiler.Profile(ctx, query, func(ctx context.Context) error {
			result, err = fr.repo.FindByID(ctx, id)
			return err
		})
	} else {
		result, err = fr.repo.FindByID(ctx, id)
	}

	if err != nil {
		return nil, err
	}

	// Cache result
	if fr.cache != nil && result != nil {
		key := fr.keyGen.KeyForID(id)
		fr.cache.Set(ctx, key, result, fr.ttl)
	}

	return result, nil
}

// Save implements Repository.Save with all features
func (fr *FullFeaturedRepository[T, ID]) Save(ctx context.Context, entity *T) (*T, error) {
	// Validate
	if fr.validator != nil {
		if err := fr.validator.Validate(entity); err != nil {
			return nil, WrapError(err, "validation failed")
		}
	}

	// Execute before hooks
	if fr.hooks != nil {
		id, _ := ExtractID[T, ID](entity)
		if IsZero(id) {
			if err := fr.hooks.ExecuteBeforeCreate(ctx, entity); err != nil {
				return nil, err
			}
		} else {
			if err := fr.hooks.ExecuteBeforeUpdate(ctx, entity); err != nil {
				return nil, err
			}
		}
	}

	// Save
	saved, err := fr.repo.Save(ctx, entity)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	if fr.cache != nil {
		fr.cache.Clear(ctx)
	}

	// Execute after hooks
	if fr.hooks != nil {
		id, _ := ExtractID[T, ID](saved)
		if IsZero(id) {
			if err := fr.hooks.ExecuteAfterCreate(ctx, saved); err != nil {
				return nil, err
			}
		} else {
			if err := fr.hooks.ExecuteAfterUpdate(ctx, saved); err != nil {
				return nil, err
			}
		}
	}

	return saved, nil
}

// Update implements Repository.Update
func (fr *FullFeaturedRepository[T, ID]) Update(ctx context.Context, entity *T) (*T, error) {
	return fr.Save(ctx, entity) // Save handles update
}

// DeleteByID implements Repository.DeleteByID
func (fr *FullFeaturedRepository[T, ID]) DeleteByID(ctx context.Context, id ID) error {
	err := fr.repo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	if fr.cache != nil {
		key := fr.keyGen.KeyForID(id)
		fr.cache.Delete(ctx, key)
	}

	return nil
}

// HealthCheck performs a health check
func (fr *FullFeaturedRepository[T, ID]) HealthCheck(ctx context.Context) HealthCheck {
	return fr.healthChecker.Check(ctx)
}

