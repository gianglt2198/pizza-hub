package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/gianglt2198/pizza-app/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrderRepository_Create(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx        context.Context
		req        model.CreateOrderRequest
		totalPrice float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderRepository{
				db: tt.fields.db,
			}
			got, err := r.Create(tt.args.ctx, tt.args.req, tt.args.totalPrice)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_FindByID(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderRepository{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx    context.Context
		id     int
		status model.OrderStatus
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderRepository{
				db: tt.fields.db,
			}
			got, err := r.UpdateStatus(tt.args.ctx, tt.args.id, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.UpdateStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrderRepository(t *testing.T) {
	type args struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want OrderRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_FindAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *OrderRepository
		args    args
		want    []model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
