package model

import "time"

type Pizza struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Category    string    `json:"category" db:"category"` // classic, premium, veggie
	Available   bool      `json:"available" db:"available"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusPreparing  OrderStatus = "preparing"
	StatusReady      OrderStatus = "ready"
	StatusDelivering OrderStatus = "delivering"
	StatusDone       OrderStatus = "done"
	StatusCancelled  OrderStatus = "cancelled"
)

type OrderItem struct {
	ID       int     `json:"id" db:"id"`
	OrderID  int     `json:"order_id" db:"order_id"`
	PizzaID  int     `json:"pizza_id" db:"pizza_id"`
	Quantity int     `json:"quantity" db:"quantity"`
	Price    float64 `json:"price" db:"price"` // snapshot price at order time
}

type Order struct {
	ID           int         `json:"id" db:"id"`
	CustomerName string      `json:"customer_name" db:"customer_name"`
	Status       OrderStatus `json:"status" db:"status"`
	TotalPrice   float64     `json:"total_price" db:"total_price"`
	Items        []OrderItem `json:"items"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
}

// Request DTOs
type CreateOrderRequest struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Items        []struct {
		PizzaID  int `json:"pizza_id" binding:"required"`
		Quantity int `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending cancelled preparing ready delivering done"`
}
