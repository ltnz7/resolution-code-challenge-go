package bootstrap

import (
	"log"
	"os"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func GetBooksAPIURL() string {
	return os.Getenv("BOOKS_API_URL")
}