package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/gianglt2198/pizza-app/internal/model"
)

func TestPizzaRepository_FindAll(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *PizzaRepository
		args    args
		want    []model.Pizza
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PizzaRepository.FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PizzaRepository.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPizzaRepository_FindByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		r       *PizzaRepository
		args    args
		want    *model.Pizza
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PizzaRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PizzaRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPizzaRepository_FindAvailable(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *PizzaRepository
		args    args
		want    []model.Pizza
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindAvailable(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("PizzaRepository.FindAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PizzaRepository.FindAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
