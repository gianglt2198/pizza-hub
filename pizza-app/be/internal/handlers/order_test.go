package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gianglt2198/pizza-app/internal/mocks"
	"github.com/gianglt2198/pizza-app/internal/model"
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
		setupMock      func(mOrder *mocks.MockOrderRepo, mPizza *mocks.MockPizzaRepo, mCache *mocks.MockCache)
		expectedStatus int
	}{
		{
			name: "success",
			reqBody: createBody("John Doe", 1, 2),
			setupMock: func(mOrder *mocks.MockOrderRepo, mPizza *mocks.MockPizzaRepo, mCache *mocks.MockCache) {
				mPizza.On("FindByID", mock.Anything, 1).Return(&model.Pizza{ID: 1, Price: 10.0}, nil)
				// 2 quantity * 10 = 20 total price
				mOrder.On("Create", mock.Anything, mock.AnythingOfType("model.CreateOrderRequest"), 20.0).
					Return(&model.Order{ID: 1, TotalPrice: 20.0, CustomerName: "John Doe"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid pizza id",
			reqBody: createBody("John Doe", 99, 2),
			setupMock: func(mOrder *mocks.MockOrderRepo, mPizza *mocks.MockPizzaRepo, mCache *mocks.MockCache) {
				mPizza.On("FindByID", mock.Anything, 99).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mOrder := new(mocks.MockOrderRepo)
			mPizza := new(mocks.MockPizzaRepo)
			mCache := new(mocks.MockCache)
			tt.setupMock(mOrder, mPizza, mCache)

			handler := NewOrderHandler(mOrder, mPizza, mCache)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewBuffer(tt.reqBody))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Create(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mOrder.AssertExpectations(t)
			mPizza.AssertExpectations(t)
			mCache.AssertExpectations(t)
		})
	}
}
