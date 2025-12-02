package providers

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"educabot.com/bookshop/models"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of BooksRepository
type mockBooksRepository struct {
	books       []models.Book
	shouldError bool
}

func (m *mockBooksRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	if m.shouldError {
		return nil, errors.New("repository error")
	}
	return m.books, nil
}

func TestBooksProvider_GetBooks_OK(t *testing.T) {
	mockRepo := &mockBooksRepository{
		books: []models.Book{
			{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 100, Price: 20},
			{ID: 2, Name: "Book 2", Author: "Author 2", UnitsSold: 200, Price: 30},
		},
	}

	provider := &booksProvider{
		repo:   mockRepo,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	books := provider.GetBooks(context.Background())

	assert.Len(t, books, 2)
	assert.Equal(t, "Book 1", books[0].Name)
	assert.Equal(t, "Author 1", books[0].Author)
}

func TestBooksProvider_GetBooks_Error(t *testing.T) {
	mockRepo := &mockBooksRepository{
		shouldError: true,
	}

	provider := &booksProvider{
		repo:   mockRepo,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	books := provider.GetBooks(context.Background())

	assert.Len(t, books, 0)
}

func TestBooksProvider_GetMetrics_OK(t *testing.T) {
	mockRepo := &mockBooksRepository{
		books: []models.Book{
			{ID: 1, Name: "The Go Programming Language", Author: "Alan Donovan", UnitsSold: 5000, Price: 40},
			{ID: 2, Name: "Clean Code", Author: "Robert C. Martin", UnitsSold: 15000, Price: 50},
			{ID: 3, Name: "The Pragmatic Programmer", Author: "Andrew Hunt", UnitsSold: 13000, Price: 45},
		},
	}

	provider := &booksProvider{
		repo:   mockRepo,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	metrics, err := provider.GetMetrics(context.Background(), "Alan Donovan")

	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, uint(11000), metrics.MeanUnitsSold)
	assert.Equal(t, "The Go Programming Language", metrics.CheapestBook)
	assert.Equal(t, uint(1), metrics.BooksWrittenByAuthor)
}

func TestBooksProvider_GetMetrics_EmptyBooks(t *testing.T) {
	mockRepo := &mockBooksRepository{
		books: []models.Book{},
	}

	provider := &booksProvider{
		repo:   mockRepo,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	metrics, err := provider.GetMetrics(context.Background(), "Any Author")

	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, uint(0), metrics.MeanUnitsSold)
	assert.Equal(t, "", metrics.CheapestBook)
	assert.Equal(t, uint(0), metrics.BooksWrittenByAuthor)
}

func TestBooksProvider_GetMetrics_NoAuthorMatch(t *testing.T) {
	mockRepo := &mockBooksRepository{
		books: []models.Book{
			{ID: 1, Name: "Book 1", Author: "Author 1", UnitsSold: 100, Price: 20},
		},
	}

	provider := &booksProvider{
		repo:   mockRepo,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	metrics, err := provider.GetMetrics(context.Background(), "Nonexistent Author")

	assert.NoError(t, err)
	assert.NotNil(t, metrics)
	assert.Equal(t, uint(100), metrics.MeanUnitsSold)
	assert.Equal(t, "Book 1", metrics.CheapestBook)
	assert.Equal(t, uint(0), metrics.BooksWrittenByAuthor)
}

func TestBooksProvider_CalculateMeanUnitsSold(t *testing.T) {
	provider := &booksProvider{}
	books := []models.Book{
		{UnitsSold: 100},
		{UnitsSold: 200},
		{UnitsSold: 300},
	}

	mean := provider.meanUnitsSold(books)
	assert.Equal(t, uint(200), mean)
}

func TestBooksProvider_cheapestBook(t *testing.T) {
	provider := &booksProvider{}
	books := []models.Book{
		{Name: "Expensive", Price: 100},
		{Name: "Cheap", Price: 20},
		{Name: "Medium", Price: 50},
	}

	cheapest := provider.cheapestBook(books)
	assert.Equal(t, "Cheap", cheapest.Name)
}

func TestBooksProvider_booksWrittenByAuthor(t *testing.T) {
	provider := &booksProvider{}
	books := []models.Book{
		{Author: "Author A"},
		{Author: "Author B"},
		{Author: "Author A"},
	}

	count := provider.booksWrittenByAuthor(books, "Author A")
	assert.Equal(t, uint(2), count)

	count = provider.booksWrittenByAuthor(books, "Author C")
	assert.Equal(t, uint(0), count)
}
