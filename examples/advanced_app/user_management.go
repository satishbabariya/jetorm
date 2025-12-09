package main

import (
	"context"
	"fmt"
	"time"

	"github.com/satishbabariya/jetorm/core"
	"github.com/satishbabariya/jetorm/hooks"
)

// Comprehensive user management system example

// User entity
type User struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	Email     string    `db:"email" jet:"unique,not_null" validate:"required,email"`
	Username  string    `db:"username" jet:"unique,not_null" validate:"required,min:3"`
	Password  string    `db:"password" jet:"not_null" validate:"required,min:8"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Role      string    `db:"role" jet:"default:'user'" validate:"in:user,admin,moderator"`
	Active    bool      `db:"active" jet:"default:true"`
	LastLogin *time.Time `db:"last_login"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

// UserProfile entity
type UserProfile struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	UserID    int64     `db:"user_id" jet:"foreign_key:users.id,on_delete:cascade,unique"`
	Bio       string    `db:"bio" jet:"type:text"`
	Avatar    string    `db:"avatar"`
	Website   string    `db:"website" validate:"url"`
	Location  string    `db:"location"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

// UserSession entity
type UserSession struct {
	ID        string    `db:"id" jet:"primary_key"`
	UserID    int64     `db:"user_id" jet:"foreign_key:users.id,on_delete:cascade"`
	Token     string    `db:"token" jet:"unique,not_null"`
	ExpiresAt time.Time `db:"expires_at" jet:"not_null"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
}

// UserManagementService provides user management operations
type UserManagementService struct {
	userRepo    core.Repository[User, int64]
	profileRepo core.Repository[UserProfile, int64]
	sessionRepo core.Repository[UserSession, string]
	cache       core.Cache
	validator   *core.Validator
}

// NewUserManagementService creates a new user management service
func NewUserManagementService(
	userRepo core.Repository[User, int64],
	profileRepo core.Repository[UserProfile, int64],
	sessionRepo core.Repository[UserSession, string],
	cache core.Cache,
) *UserManagementService {
	validator := core.NewValidator()
	validator.RegisterRule("Email", core.All(core.Required(), core.Email()))
	validator.RegisterRule("Password", core.All(
		core.Required(),
		core.MinLength(8),
		core.HasLetter(),
		core.HasDigit(),
	))
	validator.RegisterRule("Role", core.InList("user", "admin", "moderator"))

	return &UserManagementService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		sessionRepo: sessionRepo,
		cache:       cache,
		validator:   validator,
	}
}

// RegisterUser registers a new user
func (s *UserManagementService) RegisterUser(ctx context.Context, user *User) (*User, error) {
	// Validate
	if err := s.validator.Validate(user); err != nil {
		return nil, core.WrapError(err, "validation failed")
	}

	// Check if email exists
	emailSpec := core.Equal[User]("email", user.Email)
	exists, err := s.userRepo.ExistsWithSpec(ctx, emailSpec)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, core.NewCodedError(core.ErrorCodeDuplicate, "email already exists", core.ErrEntityDuplicate)
	}

	// Create user
	saved, err := s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	// Create default profile
	profile := &UserProfile{
		UserID: saved.ID,
	}
	_, err = s.profileRepo.Save(ctx, profile)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

// GetUserByEmail gets a user by email
func (s *UserManagementService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("user:email:%s", email)
	if cached, ok := s.cache.Get(ctx, cacheKey); ok {
		if user, ok := cached.(*User); ok {
			return user, nil
		}
	}

	// Query database
	spec := core.Equal[User]("email", email)
	user, err := s.userRepo.FindOne(ctx, spec)
	if err != nil {
		return nil, err
	}

	// Cache result
	if user != nil {
		s.cache.Set(ctx, cacheKey, user, 5*time.Minute)
	}

	return user, nil
}

// UpdateUserProfile updates user profile
func (s *UserManagementService) UpdateUserProfile(ctx context.Context, userID int64, profile *UserProfile) (*UserProfile, error) {
	// Validate
	if err := s.validator.Validate(profile); err != nil {
		return nil, core.WrapError(err, "validation failed")
	}

	// Find existing profile
	spec := core.Equal[UserProfile]("user_id", userID)
	existing, err := s.profileRepo.FindOne(ctx, spec)
	if err != nil && !core.IsNotFound(err) {
		return nil, err
	}

	if existing != nil {
		profile.ID = existing.ID
		return s.profileRepo.Update(ctx, profile)
	}

	profile.UserID = userID
	return s.profileRepo.Save(ctx, profile)
}

// CreateSession creates a user session
func (s *UserManagementService) CreateSession(ctx context.Context, userID int64, token string, expiresAt time.Time) (*UserSession, error) {
	session := &UserSession{
		ID:        fmt.Sprintf("%d-%d", userID, time.Now().Unix()),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return s.sessionRepo.Save(ctx, session)
}

// GetActiveUsers gets active users with pagination
func (s *UserManagementService) GetActiveUsers(ctx context.Context, page, size int) (*core.Page[User], error) {
	spec := core.Equal[User]("active", true)
	pageable := core.PageRequest(page, size, core.Order{
		Field:     "created_at",
		Direction: core.Desc,
	})
	return s.userRepo.FindAllPagedWithSpec(ctx, spec, pageable)
}

// SearchUsers searches users by query
func (s *UserManagementService) SearchUsers(ctx context.Context, query string) ([]*User, error) {
	spec := core.Or(
		core.Like[User]("email", "%"+query+"%"),
		core.Like[User]("username", "%"+query+"%"),
		core.Like[User]("first_name", "%"+query+"%"),
		core.Like[User]("last_name", "%"+query+"%"),
	)
	return s.userRepo.FindAllWithSpec(ctx, spec)
}

// GetUsersByRole gets users by role
func (s *UserManagementService) GetUsersByRole(ctx context.Context, role string) ([]*User, error) {
	spec := core.And(
		core.Equal[User]("role", role),
		core.Equal[User]("active", true),
	)
	return s.userRepo.FindAllWithSpec(ctx, spec)
}

// DeactivateUser deactivates a user
func (s *UserManagementService) DeactivateUser(ctx context.Context, userID int64) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Active = false
	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:email:%s", user.Email)
	s.cache.Delete(ctx, cacheKey)

	return nil
}

// GetUserStats gets user statistics
func (s *UserManagementService) GetUserStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Total users
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Active users
	activeSpec := core.Equal[User]("active", true)
	active, err := s.userRepo.CountWithSpec(ctx, activeSpec)
	if err != nil {
		return nil, err
	}
	stats["active"] = active

	// Users by role
	roles := []string{"user", "admin", "moderator"}
	for _, role := range roles {
		roleSpec := core.Equal[User]("role", role)
		count, err := s.userRepo.CountWithSpec(ctx, roleSpec)
		if err != nil {
			return nil, err
		}
		stats[role] = count
	}

	return stats, nil
}

func exampleUserManagement() {
	fmt.Println("User Management Example")
	fmt.Println("=======================")

	// Setup (would connect to database in real scenario)
	// db := core.Connect(config)
	// userRepo := core.NewBaseRepository[User, int64](db)
	// profileRepo := core.NewBaseRepository[UserProfile, int64](db)
	// sessionRepo := core.NewBaseRepository[UserSession, string](db)
	// cache := core.NewInMemoryCache()

	// service := NewUserManagementService(userRepo, profileRepo, sessionRepo, cache)
	// ctx := context.Background()

	// // Register user
	// user := &User{
	// 	Email:    "john@example.com",
	// 	Username: "johndoe",
	// 	Password: "securepass123",
	// 	Role:     "user",
	// }
	// savedUser, err := service.RegisterUser(ctx, user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Get user stats
	// stats, err := service.GetUserStats(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("User stats: %+v\n", stats)
}

func main() {
	exampleUserManagement()
}

