package cache

import "fmt"

const (
	KeyPizzasAll       = "pizzas:all"
	KeyPizzasAvailable = "pizzas:available"
)

func KeyPizzaByID(id int) string { return fmt.Sprintf("pizza:%d", id) }