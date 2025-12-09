package hooks

import (
	"context"
	"time"
)

// SoftDeletable interface for entities that support soft delete
type SoftDeletable interface {
	SetDeletedAt(t *time.Time)
	GetDeletedAt() *time.Time
	IsDeleted() bool
}

// SoftDeleteHook creates a hook that performs soft delete instead of hard delete
func SoftDeleteHook[T SoftDeletable]() HookFunc[T] {
	return func(ctx context.Context, entity *T) error {
		now := time.Now()
		deletable := any(*entity).(SoftDeletable)
		deletable.SetDeletedAt(&now)
		*entity = deletable.(T)
		return nil
	}
}

// RestoreHook restores a soft-deleted entity
func RestoreHook[T SoftDeletable]() HookFunc[T] {
	return func(ctx context.Context, entity *T) error {
		deletable := any(*entity).(SoftDeletable)
		deletable.SetDeletedAt(nil)
		*entity = deletable.(T)
		return nil
	}
}

// IsSoftDeleted checks if an entity is soft deleted
func IsSoftDeleted[T SoftDeletable](entity T) bool {
	return entity.IsDeleted()
}

