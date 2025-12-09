package main

import (
	"context"
	"fmt"

	"github.com/satishbabariya/jetorm/core"
	"github.com/satishbabariya/jetorm/query"
)

// Example demonstrating query building features

type User struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Username string `db:"username"`
	Age      int    `db:"age"`
	Status   string `db:"status"`
}

func exampleBasicQuery() {
	// Basic query builder
	qb := query.NewQueryBuilder("users")
	qb.WhereEqual("status", "active")
	qb.OrderBy("created_at", "DESC")
	qb.Limit(10)
	
	query, args := qb.Build()
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func exampleComposableQuery() {
	// Composable query with specification
	spec := core.And(
		core.Equal[User]("age", 18),
		core.Equal[User]("status", "active"),
	)
	
	cq := query.NewComposableQuery[User]("users")
	cq.WithSpecification(spec)
	cq.OrderBy("email", "ASC")
	cq.Limit(20)
	
	queryStr, args := cq.Build()
	fmt.Printf("Query: %s\nArgs: %v\n", queryStr, args)
}

func exampleJoinQuery() {
	// Query with joins
	jq := query.NewJoinQuery[User]("users")
	jq.InnerJoin("profiles", "users.id = profiles.user_id")
	jq.WhereEqual("users.status", "active")
	jq.Select("users.id", "users.email", "profiles.bio")
	
	query, args := jq.Build()
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func exampleConditionBuilder() {
	// Using condition builder
	cb := query.NewConditionBuilder()
	cb.Equal("status", "active")
	cb.GreaterThan("age", 18)
	cb.Like("email", "%@example.com")
	
	whereClause, args := cb.Build()
	fmt.Printf("WHERE: %s\nArgs: %v\n", whereClause, args)
}

func exampleDynamicQuery() {
	// Dynamic query building
	status := "active"
	minAge := 18
	
	dq := query.NewDynamicQuery[User]("users")
	dq.When(status != "", func(q *query.ComposableQuery[User]) *query.ComposableQuery[User] {
		return q.Where("status = $1", status)
	})
	dq.When(minAge > 0, func(q *query.ComposableQuery[User]) *query.ComposableQuery[User] {
		return q.Where("age >= $1", minAge)
	})
	
	query, args := dq.Build()
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func exampleRepositoryQuery(ctx context.Context, repo core.Repository[User, int64]) {
	// Repository-integrated query
	rq := query.NewRepositoryQuery[User, int64](repo, "users")
	rq.WhereEqual("status", "active")
	rq.OrderBy("created_at", "DESC")
	rq.Limit(10)
	
	users, err := rq.Find(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Found %d users\n", len(users))
}

func examplePagination(ctx context.Context, repo core.Repository[User, int64]) {
	// Pagination example
	rq := query.NewRepositoryQuery[User, int64](repo, "users")
	rq.WhereEqual("status", "active")
	
	pageable := core.PageRequest(0, 20, core.Order{
		Field:     "created_at",
		Direction: core.Desc,
	})
	
	page, err := rq.Paginate(ctx, pageable)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Page %d of %d, Total: %d\n", page.Number, page.TotalPages, page.TotalElements)
}

func exampleComplexQuery() {
	// Complex query with multiple conditions
	cb := query.NewConditionBuilder()
	cb.Equal("status", "active")
	cb.GreaterThan("age", 18)
	
	cb2 := query.NewConditionBuilder()
	cb2.Equal("status", "pending")
	cb2.LessThan("age", 65)
	
	// Combine with OR
	combined := cb.Or(cb2)
	
	whereClause, args := combined.Build()
	fmt.Printf("WHERE: %s\nArgs: %v\n", whereClause, args)
}

func exampleSubquery() {
	// Query with subquery
	sq := query.NewSubqueryQuery[User]("users")
	sq.Select("id", "email")
	sq.WithSubquery("SELECT COUNT(*) FROM orders WHERE orders.user_id = users.id", nil, "order_count")
	sq.WhereEqual("status", "active")
	
	query, args := sq.Build()
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func exampleHelperFunctions() {
	// Using helper functions
	helper := query.NewQueryBuilderHelper()
	
	query, args := helper.BuildSelectQuery("users",
		query.WithSelect("id", "email", "name"),
		query.WithWhere("status = $1", "active"),
		query.WithOrderBy("created_at", "DESC"),
		query.WithLimit(10),
	)
	
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func examplePostgreSQLFeatures() {
	// PostgreSQL-specific features
	cb := query.NewConditionBuilder()
	
	// Full-text search
	cb = query.TextSearch("description", "search term")
	
	// Array operations
	cb = query.ArrayContains("tags", "golang")
	cb = query.ArrayOverlaps("categories", []interface{}{"tech", "programming"})
	
	whereClause, args := cb.Build()
	fmt.Printf("WHERE: %s\nArgs: %v\n", whereClause, args)
}

func main() {
	fmt.Println("Query Building Examples")
	fmt.Println("======================")
	
	exampleBasicQuery()
	fmt.Println()
	
	exampleComposableQuery()
	fmt.Println()
	
	exampleJoinQuery()
	fmt.Println()
	
	exampleConditionBuilder()
	fmt.Println()
	
	exampleDynamicQuery()
	fmt.Println()
	
	exampleComplexQuery()
	fmt.Println()
	
	exampleSubquery()
	fmt.Println()
	
	exampleHelperFunctions()
	fmt.Println()
	
	examplePostgreSQLFeatures()
}

