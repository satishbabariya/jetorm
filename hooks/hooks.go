package hooks

import (
	"context"
	"time"
)

// HookType represents the type of lifecycle hook
type HookType int

const (
	HookBeforeCreate HookType = iota
	HookAfterCreate
	HookBeforeUpdate
	HookAfterUpdate
	HookBeforeDelete
	HookAfterDelete
	HookBeforeSave
	HookAfterSave
)

// HookFunc is a function that can be registered as a lifecycle hook
type HookFunc[T any] func(ctx context.Context, entity *T) error

// Hooks manages lifecycle hooks for an entity type
type Hooks[T any] struct {
	beforeCreate []HookFunc[T]
	afterCreate  []HookFunc[T]
	beforeUpdate []HookFunc[T]
	afterUpdate  []HookFunc[T]
	beforeDelete []HookFunc[T]
	afterDelete  []HookFunc[T]
	beforeSave   []HookFunc[T]
	afterSave    []HookFunc[T]
}

// NewHooks creates a new Hooks instance
func NewHooks[T any]() *Hooks[T] {
	return &Hooks[T]{
		beforeCreate: make([]HookFunc[T], 0),
		afterCreate:  make([]HookFunc[T], 0),
		beforeUpdate: make([]HookFunc[T], 0),
		afterUpdate:  make([]HookFunc[T], 0),
		beforeDelete: make([]HookFunc[T], 0),
		afterDelete:  make([]HookFunc[T], 0),
		beforeSave:   make([]HookFunc[T], 0),
		afterSave:    make([]HookFunc[T], 0),
	}
}

// RegisterBeforeCreate registers a hook to run before entity creation
func (h *Hooks[T]) RegisterBeforeCreate(fn HookFunc[T]) {
	h.beforeCreate = append(h.beforeCreate, fn)
}

// RegisterAfterCreate registers a hook to run after entity creation
func (h *Hooks[T]) RegisterAfterCreate(fn HookFunc[T]) {
	h.afterCreate = append(h.afterCreate, fn)
}

// RegisterBeforeUpdate registers a hook to run before entity update
func (h *Hooks[T]) RegisterBeforeUpdate(fn HookFunc[T]) {
	h.beforeUpdate = append(h.beforeUpdate, fn)
}

// RegisterAfterUpdate registers a hook to run after entity update
func (h *Hooks[T]) RegisterAfterUpdate(fn HookFunc[T]) {
	h.afterUpdate = append(h.afterUpdate, fn)
}

// RegisterBeforeDelete registers a hook to run before entity deletion
func (h *Hooks[T]) RegisterBeforeDelete(fn HookFunc[T]) {
	h.beforeDelete = append(h.beforeDelete, fn)
}

// RegisterAfterDelete registers a hook to run after entity deletion
func (h *Hooks[T]) RegisterAfterDelete(fn HookFunc[T]) {
	h.afterDelete = append(h.afterDelete, fn)
}

// RegisterBeforeSave registers a hook to run before save (create or update)
func (h *Hooks[T]) RegisterBeforeSave(fn HookFunc[T]) {
	h.beforeSave = append(h.beforeSave, fn)
}

// RegisterAfterSave registers a hook to run after save (create or update)
func (h *Hooks[T]) RegisterAfterSave(fn HookFunc[T]) {
	h.afterSave = append(h.afterSave, fn)
}

// ExecuteBeforeCreate executes all before-create hooks
func (h *Hooks[T]) ExecuteBeforeCreate(ctx context.Context, entity *T) error {
	for _, fn := range h.beforeCreate {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	for _, fn := range h.beforeSave {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteAfterCreate executes all after-create hooks
func (h *Hooks[T]) ExecuteAfterCreate(ctx context.Context, entity *T) error {
	for _, fn := range h.afterCreate {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	for _, fn := range h.afterSave {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteBeforeUpdate executes all before-update hooks
func (h *Hooks[T]) ExecuteBeforeUpdate(ctx context.Context, entity *T) error {
	for _, fn := range h.beforeUpdate {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	for _, fn := range h.beforeSave {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteAfterUpdate executes all after-update hooks
func (h *Hooks[T]) ExecuteAfterUpdate(ctx context.Context, entity *T) error {
	for _, fn := range h.afterUpdate {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	for _, fn := range h.afterSave {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteBeforeDelete executes all before-delete hooks
func (h *Hooks[T]) ExecuteBeforeDelete(ctx context.Context, entity *T) error {
	for _, fn := range h.beforeDelete {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteAfterDelete executes all after-delete hooks
func (h *Hooks[T]) ExecuteAfterDelete(ctx context.Context, entity *T) error {
	for _, fn := range h.afterDelete {
		if err := fn(ctx, entity); err != nil {
			return err
		}
	}
	return nil
}

// Auditable interface for entities that support auditing
type Auditable interface {
	SetCreatedAt(t time.Time)
	SetUpdatedAt(t time.Time)
	SetCreatedBy(userID string)
	SetUpdatedBy(userID string)
}

// AuditHook creates a hook that automatically sets audit fields
func AuditHook[T Auditable]() HookFunc[T] {
	return func(ctx context.Context, entity *T) error {
		now := time.Now()
		auditable := any(*entity).(Auditable)
		auditable.SetUpdatedAt(now)
		
		// Try to get user ID from context
		if userID, ok := ctx.Value("user_id").(string); ok {
			auditable.SetUpdatedBy(userID)
		}
		
		*entity = auditable.(T)
		return nil
	}
}

// CreateAuditHook creates a hook that sets created_at and created_by
func CreateAuditHook[T Auditable]() HookFunc[T] {
	return func(ctx context.Context, entity *T) error {
		now := time.Now()
		auditable := any(*entity).(Auditable)
		auditable.SetCreatedAt(now)
		auditable.SetUpdatedAt(now)
		
		if userID, ok := ctx.Value("user_id").(string); ok {
			auditable.SetCreatedBy(userID)
			auditable.SetUpdatedBy(userID)
		}
		
		*entity = auditable.(T)
		return nil
	}
}

