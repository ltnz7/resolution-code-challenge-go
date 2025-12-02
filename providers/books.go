package providers

import (
	"context"
	"log"
	"slices"

	"educabot.com/bookshop/models"
	"educabot.com/bookshop/repositories"
)

// BooksMetrics represents statistical metrics about books
type BooksMetrics struct {
	MeanUnitsSold        uint   `json:"mean_units_sold" example:"10000"`
	CheapestBook         string `json:"cheapest_book" example:"The Go Programming Language"`
	BooksWrittenByAuthor uint   `json:"books_written_by_author" example:"2"`
}

type BooksProvider interface {
	GetBooks(ctx context.Context) []models.Book
	GetMetrics(ctx context.Context, author string) (*BooksMetrics, error)
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

func (p *booksProvider) GetMetrics(ctx context.Context, author string) (*BooksMetrics, error) {
	books := p.GetBooks(ctx)

	if len(books) == 0 {
		return &BooksMetrics{}, nil
	}

	meanUnitsSold := p.meanUnitsSold(books)
	cheapestBook := p.cheapestBook(books)
	booksWrittenByAuthor := p.booksWrittenByAuthor(books, author)

	return &BooksMetrics{
		MeanUnitsSold:        meanUnitsSold,
		CheapestBook:         cheapestBook.Name,
		BooksWrittenByAuthor: booksWrittenByAuthor,
	}, nil
}

func (p *booksProvider) meanUnitsSold(books []models.Book) uint {
	var sum uint
	for _, book := range books {
		sum += book.UnitsSold
	}
	return sum / uint(len(books))
}

func (p *booksProvider) cheapestBook(books []models.Book) models.Book {
	return slices.MinFunc(books, func(a, b models.Book) int {
		return int(a.Price - b.Price)
	})
}

func (p *booksProvider) booksWrittenByAuthor(books []models.Book, author string) uint {
	var count uint
	for _, book := range books {
		if book.Author == author {
			count++
		}
	}
	return count
}
