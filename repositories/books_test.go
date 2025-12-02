package repositories

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

func TestHTTPBooksRepository_GetBooks_OK(t *testing.T) {
	// Mock server
	mockBooks := []models.Book{
		{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 100, Price: 20},
		{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 200, Price: 30},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockBooks)
	}))
	defer server.Close()

	// Set environment variable
	os.Setenv("BOOKS_API_URL", server.URL)
	defer os.Unsetenv("BOOKS_API_URL")

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(context.Background())

	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Equal(t, "Book 1", books[0].Name)
	assert.Equal(t, "Author 1", books[0].Author)
}

func TestHTTPBooksRepository_GetBooks_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	os.Setenv("BOOKS_API_URL", server.URL)
	defer os.Unsetenv("BOOKS_API_URL")

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(context.Background())

	// The current implementation doesn't check status codes, so it will try to decode the response
	assert.Error(t, err)
	assert.Nil(t, books)
	assert.Equal(t, "failed to decode response", err.Error())
}

func TestHTTPBooksRepository_GetBooks_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	os.Setenv("BOOKS_API_URL", server.URL)
	defer os.Unsetenv("BOOKS_API_URL")

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(context.Background())

	assert.Error(t, err)
	assert.Nil(t, books)
	assert.Equal(t, "failed to decode response", err.Error())
}

func TestHTTPBooksRepository_GetBooks_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.Book{})
	}))
	defer server.Close()

	os.Setenv("BOOKS_API_URL", server.URL)
	defer os.Unsetenv("BOOKS_API_URL")

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(context.Background())

	assert.NoError(t, err)
	assert.Len(t, books, 0)
}

func TestHTTPBooksRepository_GetBooks_NetworkError(t *testing.T) {
	// Use localhost with a port that's definitely not listening
	os.Setenv("BOOKS_API_URL", "http://localhost:99999")
	defer os.Unsetenv("BOOKS_API_URL")

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(context.Background())

	assert.Error(t, err)
	assert.Nil(t, books)
	assert.Equal(t, "failed to make HTTP request", err.Error())
}

func TestHTTPBooksRepository_GetBooks_NoURL(t *testing.T) {
	// Ensure no URL is set
	os.Unsetenv("BOOKS_API_URL")

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(context.Background())

	assert.Error(t, err)
	assert.Nil(t, books)
	assert.Equal(t, "API URL not configured", err.Error())
}

func TestHTTPBooksRepository_GetBooks_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		select {
		case <-r.Context().Done():
			return
		}
	}))
	defer server.Close()

	os.Setenv("BOOKS_API_URL", server.URL)
	defer os.Unsetenv("BOOKS_API_URL")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	repo := NewHTTPBooksRepository(log.New(os.Stdout, "", log.LstdFlags))
	books, err := repo.GetBooks(ctx)

	assert.Error(t, err)
	assert.Nil(t, books)
	assert.Equal(t, "failed to make HTTP request", err.Error())
}
