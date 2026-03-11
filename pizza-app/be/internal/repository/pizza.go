package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gianglt2198/pizza-app/internal/model"
)

type PizzaRepository struct {
	db *pgxpool.Pool
}

func NewPizzaRepository(db *pgxpool.Pool) *PizzaRepository {
	return &PizzaRepository{db: db}
}

func (r *PizzaRepository) FindAll(ctx context.Context) ([]model.Pizza, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, name, description, price, category, available, created_at
        FROM pizzas
        ORDER BY category, name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pizzas []model.Pizza
	for rows.Next() {
		var p model.Pizza
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price,
			&p.Category, &p.Available, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		pizzas = append(pizzas, p)
	}
	return pizzas, nil
}

func (r *PizzaRepository) FindByID(ctx context.Context, id int) (*model.Pizza, error) {
	var p model.Pizza
	err := r.db.QueryRow(ctx, `
        SELECT id, name, description, price, category, available, created_at
        FROM pizzas WHERE id = $1
    `, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price,
		&p.Category, &p.Available, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PizzaRepository) FindAvailable(ctx context.Context) ([]model.Pizza, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, name, description, price, category, available, created_at
        FROM pizzas WHERE available = true
        ORDER BY category, name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pizzas []model.Pizza
	for rows.Next() {
		var p model.Pizza
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price,
			&p.Category, &p.Available, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		pizzas = append(pizzas, p)
	}
	return pizzas, nil
}
