package testing

import (
	"context"
	"errors"

	"github.com/satishbabariya/jetorm/core"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository[T any, ID comparable] struct {
	SaveFunc           func(ctx context.Context, entity *T) (*T, error)
	SaveAllFunc        func(ctx context.Context, entities []*T) ([]*T, error)
	UpdateFunc         func(ctx context.Context, entity *T) (*T, error)
	UpdateAllFunc      func(ctx context.Context, entities []*T) ([]*T, error)
	FindByIDFunc       func(ctx context.Context, id ID) (*T, error)
	FindAllFunc        func(ctx context.Context) ([]*T, error)
	FindAllByIDsFunc   func(ctx context.Context, ids []ID) ([]*T, error)
	DeleteFunc         func(ctx context.Context, entity *T) error
	DeleteByIDFunc     func(ctx context.Context, id ID) error
	DeleteAllFunc      func(ctx context.Context) error
	DeleteAllByIDsFunc func(ctx context.Context, ids []ID) error
	CountFunc          func(ctx context.Context) (int64, error)
	ExistsByIdFunc     func(ctx context.Context, id ID) (bool, error)
	FindAllPagedFunc   func(ctx context.Context, pageable core.Pageable) (*core.Page[T], error)
	SaveBatchFunc      func(ctx context.Context, entities []*T, batchSize int) error
	QueryFunc          func(ctx context.Context, query string, args ...interface{}) ([]*T, error)
	QueryOneFunc       func(ctx context.Context, query string, args ...interface{}) (*T, error)
	ExecFunc           func(ctx context.Context, query string, args ...interface{}) (int64, error)
	WithTxFunc         func(ctx context.Context, fn func(ctx context.Context) error) error
}

// NewMockRepository creates a new mock repository
func NewMockRepository[T any, ID comparable]() *MockRepository[T, ID] {
	return &MockRepository[T, ID]{}
}

// Save implements Repository.Save
func (m *MockRepository[T, ID]) Save(ctx context.Context, entity *T) (*T, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, entity)
	}
	return nil, errors.New("Save not implemented")
}

// SaveAll implements Repository.SaveAll
func (m *MockRepository[T, ID]) SaveAll(ctx context.Context, entities []*T) ([]*T, error) {
	if m.SaveAllFunc != nil {
		return m.SaveAllFunc(ctx, entities)
	}
	return nil, errors.New("SaveAll not implemented")
}

// Update implements Repository.Update
func (m *MockRepository[T, ID]) Update(ctx context.Context, entity *T) (*T, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, entity)
	}
	return nil, errors.New("Update not implemented")
}

// UpdateAll implements Repository.UpdateAll
func (m *MockRepository[T, ID]) UpdateAll(ctx context.Context, entities []*T) ([]*T, error) {
	if m.UpdateAllFunc != nil {
		return m.UpdateAllFunc(ctx, entities)
	}
	return nil, errors.New("UpdateAll not implemented")
}

// FindByID implements Repository.FindByID
func (m *MockRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, core.ErrNotFound
}

// FindAll implements Repository.FindAll
func (m *MockRepository[T, ID]) FindAll(ctx context.Context) ([]*T, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return []*T{}, nil
}

// FindAllByIDs implements Repository.FindAllByIDs
func (m *MockRepository[T, ID]) FindAllByIDs(ctx context.Context, ids []ID) ([]*T, error) {
	if m.FindAllByIDsFunc != nil {
		return m.FindAllByIDsFunc(ctx, ids)
	}
	return []*T{}, nil
}

// Delete implements Repository.Delete
func (m *MockRepository[T, ID]) Delete(ctx context.Context, entity *T) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, entity)
	}
	return errors.New("Delete not implemented")
}

// DeleteByID implements Repository.DeleteByID
func (m *MockRepository[T, ID]) DeleteByID(ctx context.Context, id ID) error {
	if m.DeleteByIDFunc != nil {
		return m.DeleteByIDFunc(ctx, id)
	}
	return errors.New("DeleteByID not implemented")
}

// DeleteAll implements Repository.DeleteAll
func (m *MockRepository[T, ID]) DeleteAll(ctx context.Context) error {
	if m.DeleteAllFunc != nil {
		return m.DeleteAllFunc(ctx)
	}
	return errors.New("DeleteAll not implemented")
}

// DeleteAllByIDs implements Repository.DeleteAllByIDs
func (m *MockRepository[T, ID]) DeleteAllByIDs(ctx context.Context, ids []ID) error {
	if m.DeleteAllByIDsFunc != nil {
		return m.DeleteAllByIDsFunc(ctx, ids)
	}
	return errors.New("DeleteAllByIDs not implemented")
}

// Count implements Repository.Count
func (m *MockRepository[T, ID]) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

// ExistsById implements Repository.ExistsById
func (m *MockRepository[T, ID]) ExistsById(ctx context.Context, id ID) (bool, error) {
	if m.ExistsByIdFunc != nil {
		return m.ExistsByIdFunc(ctx, id)
	}
	return false, nil
}

// FindAllPaged implements Repository.FindAllPaged
func (m *MockRepository[T, ID]) FindAllPaged(ctx context.Context, pageable core.Pageable) (*core.Page[T], error) {
	if m.FindAllPagedFunc != nil {
		return m.FindAllPagedFunc(ctx, pageable)
	}
	return &core.Page[T]{}, nil
}

// SaveBatch implements Repository.SaveBatch
func (m *MockRepository[T, ID]) SaveBatch(ctx context.Context, entities []*T, batchSize int) error {
	if m.SaveBatchFunc != nil {
		return m.SaveBatchFunc(ctx, entities, batchSize)
	}
	return errors.New("SaveBatch not implemented")
}

// Query implements Repository.Query
func (m *MockRepository[T, ID]) Query(ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, query, args...)
	}
	return []*T{}, nil
}

// QueryOne implements Repository.QueryOne
func (m *MockRepository[T, ID]) QueryOne(ctx context.Context, query string, args ...interface{}) (*T, error) {
	if m.QueryOneFunc != nil {
		return m.QueryOneFunc(ctx, query, args...)
	}
	return nil, core.ErrNotFound
}

// Exec implements Repository.Exec
func (m *MockRepository[T, ID]) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, query, args...)
	}
	return 0, nil
}

// WithTx implements Repository.WithTx
func (m *MockRepository[T, ID]) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.WithTxFunc != nil {
		return m.WithTxFunc(ctx, fn)
	}
	return fn(ctx)
}

