package handlers

import (
	"context"
	"time"

	"github.com/gianglt2198/pizza-app/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *MockCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, keys ...string) error {
	callArgs := []any{ctx}
	for _, k := range keys {
		callArgs = append(callArgs, k)
	}
	args := m.Called(callArgs...)
	return args.Error(0)
}

type MockPizzaRepo struct {
	mock.Mock
}

func (m *MockPizzaRepo) FindAll(ctx context.Context) ([]model.Pizza, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Pizza), args.Error(1)
}

func (m *MockPizzaRepo) FindAvailable(ctx context.Context) ([]model.Pizza, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Pizza), args.Error(1)
}

func (m *MockPizzaRepo) FindByID(ctx context.Context, id int) (*model.Pizza, error) {
	args := m.Called(ctx, id)
	var p *model.Pizza
	if args.Get(0) != nil {
		p = args.Get(0).(*model.Pizza)
	}
	return p, args.Error(1)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, req model.CreateOrderRequest, totalPrice float64) (*model.Order, error) {
	args := m.Called(ctx, req, totalPrice)
	var o *model.Order
	if args.Get(0) != nil {
		o = args.Get(0).(*model.Order)
	}
	return o, args.Error(1)
}

func (m *MockOrderRepo) FindAll(ctx context.Context) ([]model.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Order), args.Error(1)
}

func (m *MockOrderRepo) FindByID(ctx context.Context, id int) (*model.Order, error) {
	args := m.Called(ctx, id)
	var o *model.Order
	if args.Get(0) != nil {
		o = args.Get(0).(*model.Order)
	}
	return o, args.Error(1)
}

func (m *MockOrderRepo) UpdateStatus(ctx context.Context, id int, status model.OrderStatus) (*model.Order, error) {
	args := m.Called(ctx, id, status)
	var o *model.Order
	if args.Get(0) != nil {
		o = args.Get(0).(*model.Order)
	}
	return o, args.Error(1)
}
