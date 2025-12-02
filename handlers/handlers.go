package handlers

import (
	"net/http"

	"educabot.com/bookshop/providers"
	"github.com/gin-gonic/gin"
)

type BooksHandler struct {
	booksProvider providers.BooksProvider
}

type GetMetricsRequest struct {
	Author string `form:"author"`
}

func NewBooksHandler(booksProvider providers.BooksProvider) *BooksHandler {
	return &BooksHandler{booksProvider: booksProvider}
}

// GetBooks godoc
// @Summary Get all books
// @Description Get a list of all available books
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func (h *BooksHandler) GetBooks(ctx *gin.Context) {
	books := h.booksProvider.GetBooks(ctx.Request.Context())
	ctx.JSON(http.StatusOK, books)
}

// GetMetrics godoc
// @Summary Get books metrics
// @Description Get statistical metrics about books, optionally filtered by author
// @Tags books
// @Accept json
// @Produce json
// @Param author query string false "Author name to filter metrics"
// @Success 200 {object} providers.BooksMetrics
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/metrics [get]
func (h *BooksHandler) GetMetrics(ctx *gin.Context) {
	var query GetMetricsRequest
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	metrics, err := h.booksProvider.GetMetrics(ctx.Request.Context(), query.Author)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics"})
		return
	}

	ctx.JSON(http.StatusOK, metrics)
}