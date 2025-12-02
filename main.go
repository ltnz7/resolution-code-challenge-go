package main

import (
	"fmt"

	"educabot.com/bookshop/handlers"
	"educabot.com/bookshop/pkg/bootstrap"
	"educabot.com/bookshop/providers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	router := gin.New()
	router.SetTrustedProxies(nil)

	booksProvider := providers.NewBooksProvider(l)
	booksHandler := handlers.NewBooksHandler(booksProvider)
	
	router.GET("/books", booksHandler.GetBooks)
	router.GET("/books/metrics", booksHandler.GetMetrics)
	router.Run(":3000")
	fmt.Println("Starting server on :3000")
}
