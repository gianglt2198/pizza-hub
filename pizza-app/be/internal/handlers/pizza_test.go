package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gianglt2198/pizza-app/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPizzaHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(mCache *MockCache, mRepo *MockPizzaRepo)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "cache hit",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
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
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
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
			mCache := new(MockCache)
			mRepo := new(MockPizzaRepo)
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

func TestPizzaHandler_GetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		idParam          string
		setupMock        func(mCache *MockCache, mRepo *MockPizzaRepo)
		expectedStatus   int
		expectedBody     string
		expectedContains string
	}{
		{
			name:    "cache hit",
			idParam: "1",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
				mCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Pizza")).
					Run(func(args mock.Arguments) {
						dest := args.Get(2).(*model.Pizza)
						*dest = model.Pizza{ID: 1, Name: "Margherita"}
					}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Margherita","description":"","price":0,"category":"","available":false,"created_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:    "cache miss, db success",
			idParam: "1",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
				mCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Pizza")).
					Return(errors.New("redis.Nil"))
				mRepo.On("FindByID", mock.Anything, mock.Anything).
					Return(&model.Pizza{ID: 1, Name: "Pepperoni"}, nil)
				mCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, 5*time.Minute).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"Pepperoni","description":"","price":0,"category":"","available":false,"created_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:    "cache miss, db not found",
			idParam: "999",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
				mCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Pizza")).
					Return(errors.New("redis.Nil"))
				mRepo.On("FindByID", mock.Anything, mock.Anything).
					Return(nil, pgx.ErrNoRows)
			},
			expectedStatus:   http.StatusNotFound,
			expectedContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mCache := new(MockCache)
			mRepo := new(MockPizzaRepo)
			tt.setupMock(mCache, mRepo)

			handler := NewPizzaHandler(mRepo, mCache)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/pizzas/"+tt.idParam, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.idParam}}

			handler.GetByID(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
			if tt.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tt.expectedContains)
			}
			mCache.AssertExpectations(t)
			mRepo.AssertExpectations(t)
		})
	}
}

func TestPizzaHandler_ListAvailable(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(mCache *MockCache, mRepo *MockPizzaRepo)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "cache hit",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
				mCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*[]model.Pizza")).
					Run(func(args mock.Arguments) {
						dest := args.Get(2).(*[]model.Pizza)
						*dest = []model.Pizza{{ID: 3, Name: "Veggie", Available: true}}
					}).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":3,"name":"Veggie","description":"","price":0,"category":"","available":true,"created_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "cache miss, db success",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
				mCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*[]model.Pizza")).
					Return(errors.New("redis.Nil"))
				mRepo.On("FindAvailable", mock.Anything).
					Return([]model.Pizza{{ID: 4, Name: "Hawaiian", Available: true}}, nil)
				mCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, 5*time.Minute).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":4,"name":"Hawaiian","description":"","price":0,"category":"","available":true,"created_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "cache miss, db empty",
			setupMock: func(mCache *MockCache, mRepo *MockPizzaRepo) {
				mCache.On("Get", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*[]model.Pizza")).
					Return(errors.New("redis.Nil"))
				mRepo.On("FindAvailable", mock.Anything).
					Return([]model.Pizza{}, nil)
				mCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, 5*time.Minute).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mCache := new(MockCache)
			mRepo := new(MockPizzaRepo)
			tt.setupMock(mCache, mRepo)

			handler := NewPizzaHandler(mRepo, mCache)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/pizzas/available", nil)

			handler.ListAvailable(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			mCache.AssertExpectations(t)
			mRepo.AssertExpectations(t)
		})
	}
}
