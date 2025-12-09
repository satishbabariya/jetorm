package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/satishbabariya/jetorm/core"
)

// User entity with relationships
type User struct {
	ID       int64   `db:"id" jet:"primary_key,auto_increment"`
	Email    string  `db:"email" jet:"unique,not_null"`
	Profile  *Profile `db:"-" jet:"one_to_one:Profile,foreign_key:user_id"`
	Posts    []*Post  `db:"-" jet:"one_to_many:Post,mapped_by:user_id"`
	Roles    []*Role  `db:"-" jet:"many_to_many:Role,join_table:user_roles,join_column:user_id,inverse_join_column:role_id"`
}

// Profile entity
type Profile struct {
	ID     int64 `db:"id" jet:"primary_key,auto_increment"`
	UserID int64 `db:"user_id" jet:"foreign_key:users.id,on_delete:cascade"`
	Bio    string `db:"bio"`
	User   *User  `db:"-" jet:"many_to_one:User,foreign_key:user_id"`
}

// Post entity
type Post struct {
	ID     int64 `db:"id" jet:"primary_key,auto_increment"`
	UserID int64 `db:"user_id" jet:"foreign_key:users.id,on_delete:cascade"`
	Title  string `db:"title"`
	User   *User  `db:"-" jet:"many_to_one:User,foreign_key:user_id"`
}

// Role entity
type Role struct {
	ID   int64  `db:"id" jet:"primary_key,auto_increment"`
	Name string `db:"name" jet:"unique"`
	Users []*User `db:"-" jet:"many_to_many:User"`
}

func exampleRelationships() {
	fmt.Println("Relationship Examples")
	fmt.Println("=====================")
	
	// Load relationships from entity
	userType := reflect.TypeOf(User{})
	relationships := core.LoadRelationships(userType)
	
	fmt.Printf("Found %d relationships for User:\n", len(relationships))
	for _, rel := range relationships {
		fmt.Printf("  - %s: %s -> %s\n", rel.Field, rel.Type, rel.TargetEntity)
	}
}

func exampleEagerLoading(ctx context.Context, userRepo core.Repository[User, int64]) {
	// Load user with all relationships
	user, err := userRepo.FindByID(ctx, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Eager load relationships
	err = core.EagerLoad(userRepo, []*User{user}, "Profile", "Posts", "Roles")
	if err != nil {
		fmt.Printf("Error loading relationships: %v\n", err)
		return
	}
	
	fmt.Printf("User loaded with relationships:\n")
	fmt.Printf("  Profile: %v\n", user.Profile)
	fmt.Printf("  Posts: %d\n", len(user.Posts))
	fmt.Printf("  Roles: %d\n", len(user.Roles))
}

func exampleLazyLoading(ctx context.Context, userRepo core.Repository[User, int64]) {
	user, err := userRepo.FindByID(ctx, 1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	// Lazy load profile when accessed
	err = core.LazyLoad(userRepo, user, "Profile")
	if err != nil {
		fmt.Printf("Error loading profile: %v\n", err)
		return
	}
	
	fmt.Printf("Profile loaded: %v\n", user.Profile)
}

// Uncomment to run example:
// func main() {
// 	exampleRelationships()
// }

