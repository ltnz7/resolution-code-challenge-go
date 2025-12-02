package providers

import (
	"context"
	"log"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/repositories"
)

type BooksProvider interface {
	GetBooks(ctx context.Context) []models.Book
}

type booksProvider struct {
	repo   repositories.BooksRepository
	logger *log.Logger
}

func NewBooksProvider(logger *log.Logger) BooksProvider {
	return &booksProvider{
		repo:   repositories.NewHTTPBooksRepository(logger),
		logger: logger,
	}
}

func (p *booksProvider) GetBooks(ctx context.Context) []models.Book {
	books, err := p.repo.GetBooks(ctx)
	if err != nil {
		p.logger.Printf("Error fetching books: %v", err)
		return []models.Book{}
	}
	return books
}
