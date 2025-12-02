package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/pkg/bootstrap"
)

type BooksRepository interface {
	GetBooks(ctx context.Context) ([]models.Book, error)
}

type HTTPBooksRepository struct {
	client *http.Client
	logger *log.Logger
}

func NewHTTPBooksRepository(logger *log.Logger) BooksRepository {
	return &HTTPBooksRepository{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (r *HTTPBooksRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	url := bootstrap.GetBooksAPIURL()

	if url == "" {
		r.logger.Println("BOOKS_API_URL not configured")
		return nil, errors.New("API URL not configured")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		r.logger.Printf("Error creating request: %v", err)
		return nil, errors.New("failed to create request")
	}

	resp, err := r.client.Do(req)
	if err != nil {
		r.logger.Printf("Error making HTTP request: %v", err)
		return nil, errors.New("failed to make HTTP request")
	}
	defer resp.Body.Close()

	var books []models.Book
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		r.logger.Printf("Error decoding response: %v", err)
		return nil, errors.New("failed to decode response")
	}

	return books, nil
}
