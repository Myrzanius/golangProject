package models

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockLevel  int     `json:"stock_level"`
	CategoryID  string  `json:"category_id"`
}
