// @title Bookshop API
// @description resolución de challenge tecnico en go
// @version 1.0
// @host localhost:3000
// @BasePath /
package main

import (
	"fmt"

	"educabot.com/bookshop/handlers"
	"educabot.com/bookshop/pkg/bootstrap"
	"educabot.com/bookshop/providers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "educabot.com/bookshop/docs"
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
	
	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	fmt.Println("Starting server on :3000")
	fmt.Println("Swagger documentation available at: http://localhost:3000/swagger/index.html")
	router.Run(":3000")
}
