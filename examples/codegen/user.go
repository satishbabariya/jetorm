package main

//go:generate jetorm-gen -type=User -interface=UserRepository -input=user.go -output=user_repository_gen.go

import (
	"context"
	"time"
)

// User represents a user entity
type User struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	Email     string    `db:"email" jet:"unique,not_null"`
	Username  string    `db:"username" jet:"unique,not_null"`
	Age       int       `db:"age"`
	Status    string    `db:"status" jet:"default:'active'"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

// UserRepository defines the repository interface for User
type UserRepository interface {
	// Base repository methods are inherited from jetorm.Repository[User, int64]
	
	// Custom query methods (auto-implemented by code generation)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByAgeGreaterThan(ctx context.Context, age int) ([]*User, error)
	FindByStatusIn(ctx context.Context, statuses []string) ([]*User, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	DeleteByEmail(ctx context.Context, email string) error
	FindByStatusOrderByCreatedAtDesc(ctx context.Context, status string) ([]*User, error)
	
	// Custom implementations (manually written)
	FindActiveUsersWithOrders(ctx context.Context) ([]*User, error)
}

