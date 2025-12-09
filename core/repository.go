package core

import "context"

// Repository is the generic repository interface providing CRUD operations
type Repository[T any, ID comparable] interface {
	// Basic CRUD
	Save(ctx context.Context, entity *T) (*T, error)
	SaveAll(ctx context.Context, entities []*T) ([]*T, error)
	FindByID(ctx context.Context, id ID) (*T, error)
	FindAll(ctx context.Context) ([]*T, error)
	FindAllByIDs(ctx context.Context, ids []ID) ([]*T, error)
	Delete(ctx context.Context, entity *T) error
	DeleteByID(ctx context.Context, id ID) error
	DeleteAll(ctx context.Context, entities []*T) error
	Count(ctx context.Context) (int64, error)
	ExistsById(ctx context.Context, id ID) (bool, error)

	// Pagination
	FindAllPaged(ctx context.Context, pageable Pageable) (*Page[T], error)

	// Transaction
	WithTx(tx *Tx) Repository[T, ID]
}

// Pageable represents pagination and sorting request
type Pageable struct {
	Page int  // Zero-based page number
	Size int  // Page size
	Sort Sort // Sort specification
}

// Sort represents sort specification
type Sort struct {
	Orders []Order
}

// Order represents a single sort order
type Order struct {
	Field     string
	Direction Direction
}

// Direction represents sort direction
type Direction int

const (
	Asc Direction = iota
	Desc
)

// Page represents a page of results
type Page[T any] struct {
	Content       []*T  // Page content
	Number        int   // Current page number
	Size          int   // Page size
	TotalElements int64 // Total elements
	TotalPages    int   // Total pages
	First         bool  // Is first page
	Last          bool  // Is last page
}

// PageRequest creates a Pageable with the given page, size and sort orders
func PageRequest(page, size int, orders ...Order) Pageable {
	return Pageable{
		Page: page,
		Size: size,
		Sort: Sort{Orders: orders},
	}
}

// Unpaged creates a Pageable that represents no pagination
func Unpaged() Pageable {
	return Pageable{
		Page: 0,
		Size: -1,
	}
}

