package testing

import (
	"context"
	"time"
)

// Fixture represents a test data fixture
type Fixture interface {
	Setup(ctx context.Context) error
	Teardown(ctx context.Context) error
}

// FixtureManager manages test fixtures
type FixtureManager struct {
	fixtures []Fixture
}

// NewFixtureManager creates a new fixture manager
func NewFixtureManager() *FixtureManager {
	return &FixtureManager{
		fixtures: make([]Fixture, 0),
	}
}

// Register registers a fixture
func (fm *FixtureManager) Register(fixture Fixture) {
	fm.fixtures = append(fm.fixtures, fixture)
}

// SetupAll sets up all registered fixtures
func (fm *FixtureManager) SetupAll(ctx context.Context) error {
	for _, fixture := range fm.fixtures {
		if err := fixture.Setup(ctx); err != nil {
			return err
		}
	}
	return nil
}

// TeardownAll tears down all registered fixtures
func (fm *FixtureManager) TeardownAll(ctx context.Context) error {
	// Teardown in reverse order
	for i := len(fm.fixtures) - 1; i >= 0; i-- {
		if err := fm.fixtures[i].Teardown(ctx); err != nil {
			return err
		}
	}
	return nil
}

// UserFixture is an example fixture for user entities
type UserFixture struct {
	Users []interface{}
}

// Setup sets up user test data
func (uf *UserFixture) Setup(ctx context.Context) error {
	// Example: Create test users
	uf.Users = []interface{}{
		map[string]interface{}{
			"email":    "test1@example.com",
			"username": "testuser1",
			"age":      25,
		},
		map[string]interface{}{
			"email":    "test2@example.com",
			"username": "testuser2",
			"age":      30,
		},
	}
	return nil
}

// Teardown cleans up user test data
func (uf *UserFixture) Teardown(ctx context.Context) error {
	uf.Users = nil
	return nil
}

// TimeFixture provides time utilities for testing
type TimeFixture struct{}

// Now returns the current time
func (tf *TimeFixture) Now() time.Time {
	return time.Now()
}

// FixedTime returns a fixed time for testing
func (tf *TimeFixture) FixedTime() time.Time {
	return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
}

