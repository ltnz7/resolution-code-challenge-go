package models

// Book represents a book in the bookshop
type Book struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"The Go Programming Language"`
	Author    string `json:"author" example:"Alan Donovan"`
	UnitsSold uint   `json:"units_sold" example:"5000"`
	Price     uint   `json:"price" example:"45"`
}
