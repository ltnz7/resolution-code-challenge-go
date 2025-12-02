package main

import (
	"fmt"

	"educabot.com/bookshop/handlers"
	"educabot.com/bookshop/pkg/bootstrap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	router := gin.New()
	router.SetTrustedProxies(nil)

	metricsHandler := handlers.NewGetMetrics(mockImpls.NewMockBooksProvider())
	router.GET("/", metricsHandler.Handle())
	router.Run(":3000")
	fmt.Println("Starting server on :3000")
}
