package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/providers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of BooksProvider
type mockBooksProvider struct {
	books       []models.Book
	shouldError bool
}

func (m *mockBooksProvider) GetBooks(ctx context.Context) []models.Book {
	return m.books
}

func (m *mockBooksProvider) GetMetrics(ctx context.Context, author string) (*providers.BooksMetrics, error) {
	if m.shouldError {
		return nil, errors.New("provider error")
	}
	return &providers.BooksMetrics{
		MeanUnitsSold:        10000,
		CheapestBook:         "The Go Programming Language",
		BooksWrittenByAuthor: 1,
	}, nil
}

func TestGetBooks_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{
		books: []models.Book{
			{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 100, Price: 20},
			{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 200, Price: 30},
		},
	}

	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books", handler.GetBooks)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var books []models.Book
	err := json.Unmarshal(res.Body.Bytes(), &books)
	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Equal(t, "Book 1", books[0].Name)
}

func TestGetBooks_EmptyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{
		books: []models.Book{},
	}

	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books", handler.GetBooks)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var books []models.Book
	err := json.Unmarshal(res.Body.Bytes(), &books)
	assert.NoError(t, err)
	assert.Len(t, books, 0)
}

func TestGetMetrics_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{
		shouldError: false,
	}

	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books/metrics", handler.GetMetrics)

	req := httptest.NewRequest(http.MethodGet, "/books/metrics?author=Alan+Donovan", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)

	assert.Equal(t, 10000, int(resBody["mean_units_sold"].(float64)))
	assert.Equal(t, "The Go Programming Language", resBody["cheapest_book"])
	assert.Equal(t, 1, int(resBody["books_written_by_author"].(float64)))
}

func TestGetMetrics_ValidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{}
	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books/metrics", handler.GetMetrics)

	// Test with valid query - the current struct doesn't have validation that would cause binding errors
	// So we test that it accepts various valid formats
	req := httptest.NewRequest(http.MethodGet, "/books/metrics?author=Valid+Author", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)
	assert.Equal(t, 10000, int(resBody["mean_units_sold"].(float64)))
}

func TestGetMetrics_ProviderError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{
		shouldError: true,
	}

	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books/metrics", handler.GetMetrics)

	req := httptest.NewRequest(http.MethodGet, "/books/metrics?author=Test+Author", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to get metrics", resBody["error"])
}

func TestGetMetrics_EmptyAuthor(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{
		shouldError: false,
	}

	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books/metrics", handler.GetMetrics)

	req := httptest.NewRequest(http.MethodGet, "/books/metrics?author=", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)

	// Should still return metrics even with empty author
	assert.Equal(t, 10000, int(resBody["mean_units_sold"].(float64)))
}

func TestGetMetrics_NoAuthorParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockProvider := &mockBooksProvider{
		shouldError: false,
	}

	handler := NewBooksHandler(mockProvider)
	r := gin.Default()
	r.GET("/books/metrics", handler.GetMetrics)

	// Test without author parameter at all
	req := httptest.NewRequest(http.MethodGet, "/books/metrics", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var resBody map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &resBody)
	assert.NoError(t, err)

	// Should return metrics with empty author (empty string)
	assert.Equal(t, 10000, int(resBody["mean_units_sold"].(float64)))
}
