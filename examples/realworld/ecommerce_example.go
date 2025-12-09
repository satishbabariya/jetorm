package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/satishbabariya/jetorm/core"
	"github.com/satishbabariya/jetorm/hooks"
)

// Real-world e-commerce example

// Product entity
type Product struct {
	ID          int64     `db:"id" jet:"primary_key,auto_increment"`
	Name        string    `db:"name" jet:"not_null" validate:"required,min:3"`
	SKU         string    `db:"sku" jet:"unique,not_null" validate:"required"`
	Price       float64   `db:"price" jet:"not_null" validate:"required,positive"`
	Stock       int       `db:"stock" jet:"default:0" validate:"min:0"`
	CategoryID  int64     `db:"category_id" jet:"foreign_key:categories.id"`
	CreatedAt   time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt   time.Time `db:"updated_at" jet:"auto_now"`
}

// Category entity
type Category struct {
	ID        int64     `db:"id" jet:"primary_key,auto_increment"`
	Name      string    `db:"name" jet:"unique,not_null" validate:"required"`
	Slug      string    `db:"slug" jet:"unique,not_null" validate:"required"`
	CreatedAt time.Time `db:"created_at" jet:"auto_now_add"`
}

// Order entity
type Order struct {
	ID         int64     `db:"id" jet:"primary_key,auto_increment"`
	UserID     int64     `db:"user_id" jet:"foreign_key:users.id"`
	Status     string    `db:"status" jet:"default:'pending'" validate:"in:pending,processing,shipped,delivered,cancelled"`
	Total      float64   `db:"total" jet:"not_null" validate:"required,positive"`
	CreatedAt  time.Time `db:"created_at" jet:"auto_now_add"`
	UpdatedAt  time.Time `db:"updated_at" jet:"auto_now"`
}

// OrderItem entity
type OrderItem struct {
	ID        int64   `db:"id" jet:"primary_key,auto_increment"`
	OrderID   int64   `db:"order_id" jet:"foreign_key:orders.id,on_delete:cascade"`
	ProductID int64   `db:"product_id" jet:"foreign_key:products.id"`
	Quantity  int     `db:"quantity" jet:"not_null" validate:"required,min:1"`
	Price     float64 `db:"price" jet:"not_null" validate:"required,positive"`
}

// ECommerceService provides e-commerce operations
type ECommerceService struct {
	productRepo core.Repository[Product, int64]
	categoryRepo core.Repository[Category, int64]
	orderRepo    core.Repository[Order, int64]
	orderItemRepo core.Repository[OrderItem, int64]
}

// NewECommerceService creates a new e-commerce service
func NewECommerceService(
	productRepo core.Repository[Product, int64],
	categoryRepo core.Repository[Category, int64],
	orderRepo core.Repository[Order, int64],
	orderItemRepo core.Repository[OrderItem, int64],
) *ECommerceService {
	return &ECommerceService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		orderRepo:    orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

// CreateProduct creates a new product
func (s *ECommerceService) CreateProduct(ctx context.Context, product *Product) (*Product, error) {
	// Validate
	validator := core.NewValidator()
	validator.RegisterRule("Name", core.All(core.Required(), core.MinLength(3)))
	validator.RegisterRule("Price", core.All(core.Required(), core.Positive()))
	
	if err := validator.Validate(product); err != nil {
		return nil, err
	}

	return s.productRepo.Save(ctx, product)
}

// GetProductsByCategory gets products by category
func (s *ECommerceService) GetProductsByCategory(ctx context.Context, categoryID int64) ([]*Product, error) {
	spec := core.Equal[Product]("category_id", categoryID)
	return s.productRepo.FindAllWithSpec(ctx, spec)
}

// CreateOrder creates a new order with items
func (s *ECommerceService) CreateOrder(ctx context.Context, userID int64, items []OrderItem) (*Order, error) {
	// Calculate total
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}

	// Create order
	order := &Order{
		UserID: userID,
		Status: "pending",
		Total:  total,
	}

	// Use transaction
	var savedOrder *Order
	err := s.orderRepo.(*core.BaseRepository[Order, int64]).db.Transaction(ctx, func(tx *core.Tx) error {
		// Save order
		var err error
		savedOrder, err = s.orderRepo.WithTx(tx).Save(ctx, order)
		if err != nil {
			return err
		}

		// Save order items
		txRepo := s.orderItemRepo.WithTx(tx)
		for _, item := range items {
			item.OrderID = savedOrder.ID
			if _, err := txRepo.Save(ctx, &item); err != nil {
				return err
			}
		}

		return nil
	})

	return savedOrder, err
}

// UpdateOrderStatus updates order status
func (s *ECommerceService) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Validate status
	validator := core.NewValidator()
	validator.RegisterRule("Status", core.InList("pending", "processing", "shipped", "delivered", "cancelled"))
	
	order.Status = status
	if err := validator.Validate(order); err != nil {
		return err
	}

	_, err = s.orderRepo.Update(ctx, order)
	return err
}

// GetLowStockProducts gets products with low stock
func (s *ECommerceService) GetLowStockProducts(ctx context.Context, threshold int) ([]*Product, error) {
	spec := core.LessThan[Product]("stock", threshold)
	return s.productRepo.FindAllWithSpec(ctx, spec)
}

// GetOrderHistory gets order history for a user
func (s *ECommerceService) GetOrderHistory(ctx context.Context, userID int64, page, size int) (*core.Page[Order], error) {
	spec := core.Equal[Order]("user_id", userID)
	pageable := core.PageRequest(page, size, core.Order{
		Field:     "created_at",
		Direction: core.Desc,
	})
	return s.orderRepo.FindAllPagedWithSpec(ctx, spec, pageable)
}

func exampleECommerce() {
	fmt.Println("E-Commerce Example")
	fmt.Println("==================")

	// Setup (would connect to database in real scenario)
	// db := core.Connect(config)
	// productRepo := core.NewBaseRepository[Product, int64](db)
	// categoryRepo := core.NewBaseRepository[Category, int64](db)
	// orderRepo := core.NewBaseRepository[Order, int64](db)
	// orderItemRepo := core.NewBaseRepository[OrderItem, int64](db)

	// service := NewECommerceService(productRepo, categoryRepo, orderRepo, orderItemRepo)

	// ctx := context.Background()

	// // Create product
	// product := &Product{
	// 	Name:       "Laptop",
	// 	SKU:        "LAP-001",
	// 	Price:      999.99,
	// 	Stock:      10,
	// 	CategoryID: 1,
	// }
	// savedProduct, err := service.CreateProduct(ctx, product)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Created product: %s (ID: %d)\n", savedProduct.Name, savedProduct.ID)

	// // Get low stock products
	// lowStock, err := service.GetLowStockProducts(ctx, 5)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Low stock products: %d\n", len(lowStock))
}

func main() {
	exampleECommerce()
}

