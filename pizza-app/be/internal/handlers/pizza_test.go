package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gianglt2198/pizza-app/internal/mocks"
	"github.com/gianglt2198/pizza-app/internal/model"
)

func TestPizzaHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(mCache *mocks.MockCache, mRepo *mocks.MockPizzaRepo)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "cache hit",
			setupMock: func(mCache *mocks.MockCache, mRepo *mocks.MockPizzaRepo) {
				mCache.On("Get", mock.Anything, "pizzas:all", mock.AnythingOfType("*[]model.Pizza")).
					Run(func(args mock.Arguments) {
						dest := args.Get(2).(*[]model.Pizza)
						*dest = []model.Pizza{{ID: 1, Name: "Margherita"}}
					}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"Margherita","description":"","price":0,"category":"","available":false,"created_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "cache miss, db success",
			setupMock: func(mCache *mocks.MockCache, mRepo *mocks.MockPizzaRepo) {
				mCache.On("Get", mock.Anything, "pizzas:all", mock.Anything).Return(errors.New("redis.Nil"))
				mRepo.On("FindAll", mock.Anything).Return([]model.Pizza{{ID: 2, Name: "Pepperoni"}}, nil)
				mCache.On("Set", mock.Anything, "pizzas:all", mock.Anything, 5*time.Minute).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":2,"name":"Pepperoni","description":"","price":0,"category":"","available":false,"created_at":"0001-01-01T00:00:00Z"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mCache := new(mocks.MockCache)
			mRepo := new(mocks.MockPizzaRepo)
			tt.setupMock(mCache, mRepo)

			handler := NewPizzaHandler(mRepo, mCache)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/pizzas", nil)

			handler.List(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mCache.AssertExpectations(t)
			mRepo.AssertExpectations(t)
		})
	}
}
