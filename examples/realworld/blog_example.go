package main

import (
	"context"
	"fmt"
	"time"

	"github.com/satishbabariya/jetorm/core"
)

// Blog example with posts, comments, and tags

// Post entity
type Post struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	Title     string    `db:"title" jet:"not_null" validate:"required,min:5"`
	Slug      string    `db:"slug" jet:"unique,not_null" validate:"required"`
	Content   string    `db:"content" jet:"type:text"`
	AuthorID  int64     `db:"author_id" jet:"foreign_key:users.id"`
	Published bool      `db:"published" jet:"default:false"`
	Views     int64     `db:"views" jet:"default:0"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt time.Time `db:"updated_at" jet:"auto_now"`
}

// Comment entity
type Comment struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	PostID    int64     `db:"post_id" jet:"foreign_key:posts.id,on_delete:cascade"`
	AuthorID  int64     `db:"author_id" jet:"foreign_key:users.id"`
	Content   string    `db:"content" jet:"type:text,not_null" validate:"required,min:10"`
	Approved  bool      `db:"approved" jet:"default:false"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
}

// Tag entity
type Tag struct {
	ID   int64  `db:"id" jet:"primary_key,auto_increment"`
	Name string `db:"name" jet:"unique,not_null" validate:"required"`
	Slug string `db:"slug" jet:"unique,not_null" validate:"required"`
}

// BlogService provides blog operations
type BlogService struct {
	postRepo    core.Repository[Post, int64]
	commentRepo core.Repository[Comment, int64]
	tagRepo     core.Repository[Tag, int64]
}

// NewBlogService creates a new blog service
func NewBlogService(
	postRepo core.Repository[Post, int64],
	commentRepo core.Repository[Comment, int64],
	tagRepo core.Repository[Tag, int64],
) *BlogService {
	return &BlogService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		tagRepo:     tagRepo,
	}
}

// CreatePost creates a new blog post
func (s *BlogService) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	// Validate
	validator := core.NewValidator()
	validator.RegisterRule("Title", core.All(core.Required(), core.MinLength(5)))
	validator.RegisterRule("Slug", core.Required())
	
	if err := validator.Validate(post); err != nil {
		return nil, err
	}

	return s.postRepo.Save(ctx, post)
}

// GetPublishedPosts gets published posts with pagination
func (s *BlogService) GetPublishedPosts(ctx context.Context, page, size int) (*core.Page[Post], error) {
	spec := core.Equal[Post]("published", true)
	pageable := core.PageRequest(page, size, core.Order{
		Field:     "created_at",
		Direction: core.Desc,
	})
	return s.postRepo.FindAllPagedWithSpec(ctx, spec, pageable)
}

// GetPostBySlug gets a post by slug
func (s *BlogService) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	spec := core.Equal[Post]("slug", slug)
	return s.postRepo.FindOne(ctx, spec)
}

// AddComment adds a comment to a post
func (s *BlogService) AddComment(ctx context.Context, comment *Comment) (*Comment, error) {
	// Validate
	validator := core.NewValidator()
	validator.RegisterRule("Content", core.All(core.Required(), core.MinLength(10)))
	
	if err := validator.Validate(comment); err != nil {
		return nil, err
	}

	return s.commentRepo.Save(ctx, comment)
}

// GetPostComments gets comments for a post
func (s *BlogService) GetPostComments(ctx context.Context, postID int64) ([]*Comment, error) {
	spec := core.And(
		core.Equal[Comment]("post_id", postID),
		core.Equal[Comment]("approved", true),
	)
	return s.commentRepo.FindAllWithSpec(ctx, spec)
}

// IncrementViews increments post view count
func (s *BlogService) IncrementViews(ctx context.Context, postID int64) error {
	post, err := s.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}

	post.Views++
	_, err = s.postRepo.Update(ctx, post)
	return err
}

// SearchPosts searches posts by title or content
func (s *BlogService) SearchPosts(ctx context.Context, query string) ([]*Post, error) {
	spec := core.Or(
		core.Like[Post]("title", "%"+query+"%"),
		core.Like[Post]("content", "%"+query+"%"),
	)
	return s.postRepo.FindAllWithSpec(ctx, spec)
}

func exampleBlog() {
	fmt.Println("Blog Example")
	fmt.Println("============")

	// Setup (would connect to database in real scenario)
	// db := core.Connect(config)
	// postRepo := core.NewBaseRepository[Post, int64](db)
	// commentRepo := core.NewBaseRepository[Comment, int64](db)
	// tagRepo := core.NewBaseRepository[Tag, int64](db)

	// service := NewBlogService(postRepo, commentRepo, tagRepo)
	// ctx := context.Background()

	// // Create post
	// post := &Post{
	// 	Title:     "Getting Started with JetORM",
	// 	Slug:      "getting-started-with-jetorm",
	// 	Content:   "This is a comprehensive guide...",
	// 	AuthorID:  1,
	// 	Published: true,
	// }
	// savedPost, err := service.CreatePost(ctx, post)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Get published posts
	// posts, err := service.GetPublishedPosts(ctx, 0, 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Published posts: %d\n", len(posts.Content))
}

func main() {
	exampleBlog()
}

