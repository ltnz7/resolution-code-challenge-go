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

func (h *BooksHandler) GetBooks(ctx *gin.Context) {
	books := h.booksProvider.GetBooks(ctx.Request.Context())
	ctx.JSON(http.StatusOK, books)
}

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