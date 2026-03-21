package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gianglt2198/pizza-app/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createBody := func(name string, pID, qty int) []byte {
		body := map[string]any{
			"customer_name": name,
			"items": []map[string]any{
				{"pizza_id": pID, "quantity": qty},
			},
		}
		b, _ := json.Marshal(body)
		return b
	}

	tests := []struct {
		name           string
		reqBody        []byte
		setupMock      func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo)
		expectedStatus int
	}{
		{
			name:    "success",
			reqBody: createBody("John Doe", 1, 2),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mPizza.On("FindByID", mock.Anything, 1).Return(&model.Pizza{ID: 1, Price: 10.0}, nil)
				// 2 quantity * 10 = 20 total price
				mOrder.On("Create", mock.Anything, mock.AnythingOfType("model.CreateOrderRequest"), 20.0).
					Return(&model.Order{ID: 1, TotalPrice: 20.0, CustomerName: "John Doe"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:    "invalid pizza id",
			reqBody: createBody("John Doe", 99, 2),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mPizza.On("FindByID", mock.Anything, 99).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mOrder := new(MockOrderRepo)
			mPizza := new(MockPizzaRepo)
			tt.setupMock(mOrder, mPizza)

			handler := NewOrderHandler(mOrder, mPizza)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(tt.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Create(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mOrder.AssertExpectations(t)
			mPizza.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		idParam        string
		setupMock      func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo)
		expectedStatus int
	}{
		{
			name:    "order not found",
			idParam: "999",
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mOrder.On("FindByID", context.Background(), mock.Anything).Return(nil, pgx.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "cache hit",
			idParam: "1",
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mOrder.On("FindByID", context.Background(), mock.Anything).Return(&model.Order{ID: 1, Status: model.StatusPending}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mOrder := new(MockOrderRepo)
			mPizza := new(MockPizzaRepo)
			tt.setupMock(mOrder, mPizza)

			handler := NewOrderHandler(mOrder, mPizza)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/orders/"+tt.idParam, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.idParam}}

			handler.GetByID(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mOrder.AssertExpectations(t)
			mPizza.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_UpdateStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	makeBody := func(status string) []byte {
		b, _ := json.Marshal(map[string]any{"status": status})
		return b
	}

	tests := []struct {
		name              string
		idParam           string
		reqBody           []byte
		setupMock         func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo)
		expectedStatus    int
		assertNoRepoCall  bool
	}{
		{
			name:    "invalid id",
			idParam: "abc",
			reqBody: makeBody("completed"),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				// no repo expectation; handler should return 400 before repo call
			},
			expectedStatus:   http.StatusBadRequest,
			assertNoRepoCall: true,
		},
		{
			name:    "invalid status",
			idParam: "1",
			reqBody: makeBody("wrong-status"),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				// handler calls repo with provided status; mock it explicitly
			},
			expectedStatus:   http.StatusBadRequest,
			assertNoRepoCall: true,
		},
		{
			name:    "success",
			idParam: "1",
			reqBody: makeBody("done"),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mOrder.On("UpdateStatus", context.Background(), mock.Anything, model.OrderStatus("done")).
					Return(&model.Order{ID: 1, Status: "done"}, nil)
			},
			expectedStatus:   http.StatusOK,
			assertNoRepoCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mOrder := new(MockOrderRepo)
			mPizza := new(MockPizzaRepo)
			tt.setupMock(mOrder, mPizza)

			handler := NewOrderHandler(mOrder, mPizza)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPatch, "/api/v1/orders/"+tt.idParam+"/status", bytes.NewBuffer(tt.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.idParam}}

			handler.UpdateStatus(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.assertNoRepoCall {
				mOrder.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything, mock.Anything)
			}
			mOrder.AssertExpectations(t)
			mPizza.AssertExpectations(t)
		})
	}
}

func TestOrderHandler_calculateTotal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	createBody := func(name string, items []map[string]any) []byte {
		body := map[string]any{
			"customer_name": name,
			"items":         items,
		}
		b, _ := json.Marshal(body)
		return b
	}

	tests := []struct {
		name           string
		reqBody        []byte
		setupMock      func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo)
		expectedStatus int
	}{
		{
			name: "calculate total success with multiple items",
			reqBody: createBody("Jane", []map[string]any{
				{"pizza_id": 1, "quantity": 2}, // 2 * 10
				{"pizza_id": 2, "quantity": 1}, // 1 * 15
			}),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mPizza.On("FindByID", context.Background(), 1).Return(&model.Pizza{ID: 1, Price: 10.0}, nil)
				mPizza.On("FindByID", context.Background(), 2).Return(&model.Pizza{ID: 2, Price: 15.0}, nil)
				mOrder.On("Create", context.Background(), mock.AnythingOfType("model.CreateOrderRequest"), 35.0).
					Return(&model.Order{ID: 10, TotalPrice: 35.0, CustomerName: "Jane"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "calculate total fails when pizza lookup returns error",
			reqBody: createBody("Jane", []map[string]any{
				{"pizza_id": 1, "quantity": 2},
				{"pizza_id": 999, "quantity": 1},
			}),
			setupMock: func(mOrder *MockOrderRepo, mPizza *MockPizzaRepo) {
				mPizza.On("FindByID", context.Background(), 1).Return(&model.Pizza{ID: 1, Price: 10.0}, nil)
				mPizza.On("FindByID", context.Background(), 999).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mOrder := new(MockOrderRepo)
			mPizza := new(MockPizzaRepo)
			tt.setupMock(mOrder, mPizza)

			handler := NewOrderHandler(mOrder, mPizza)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(tt.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Create(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus != http.StatusCreated {
				mOrder.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
			}
			mOrder.AssertExpectations(t)
			mPizza.AssertExpectations(t)
		})
	}
}
