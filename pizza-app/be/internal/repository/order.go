package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gianglt2198/pizza-app/internal/model"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, req model.CreateOrderRequest, totalPrice float64) (*model.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var order model.Order
	err = tx.QueryRow(ctx, `
        INSERT INTO orders (customer_name, status, total_price)
        VALUES ($1, $2, $3)
        RETURNING id, customer_name, status, total_price, created_at, updated_at
    `, req.CustomerName, model.StatusPending, totalPrice).
		Scan(&order.ID, &order.CustomerName, &order.Status,
			&order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		var oi model.OrderItem
		err = tx.QueryRow(ctx, `
            INSERT INTO order_items (order_id, pizza_id, quantity, price)
            SELECT $1, $2, $3, price FROM pizzas WHERE id = $2
            RETURNING id, order_id, pizza_id, quantity, price
        `, order.ID, item.PizzaID, item.Quantity).
			Scan(&oi.ID, &oi.OrderID, &oi.PizzaID, &oi.Quantity, &oi.Price)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, oi)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id int) (*model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(ctx, `
        SELECT id, customer_name, status, total_price, created_at, updated_at
        FROM orders WHERE id = $1
    `, id).Scan(&order.ID, &order.CustomerName, &order.Status,
		&order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
        SELECT id, order_id, pizza_id, quantity, price
        FROM order_items WHERE order_id = $1
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.PizzaID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return &order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (*model.Order, error) {
	var order model.Order
	err := r.db.QueryRow(ctx, `
        UPDATE orders SET status = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING id, customer_name, status, total_price, created_at, updated_at
    `, status, id).Scan(&order.ID, &order.CustomerName, &order.Status,
		&order.TotalPrice, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, customer_name, status, total_price, created_at, updated_at
        FROM orders ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.ID, &o.CustomerName, &o.Status,
			&o.TotalPrice, &o.CreatedAt, &o.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
