package models

import (
	"time"
)

type Product struct {
	ProductID   int      `json:"product_id"`
	Name        string   `json:"name"`
	Brand       string   `json:"brand"`
	Price       float32  `json:"price"`
	Amount      int      `json:"amount"`
	Description string   `json:"description"`
	Category    Category `json:"category"`
}

type Category struct {
	CategoryID int    `json:"category_id"`
	Name       string `json:"name"`
}

type Provider struct {
	ProviderID int    `json:"provider_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Enterprise string `json:"enterprise"`
}

type ProductProvider struct {
	DeliveryDate time.Time `json:"delivery_date"`
	ProductID    int       `json:"product_id"`
	ProviderID   int       `json:"provider_id"`
	Product      Product   `json:"product"`
	Provider     Provider  `json:"provider"`
}

type Sale struct {
	SaleID    int       `json:"sale_id"`
	ProductID int       `json:"product_id"`
	Date      time.Time `json:"date"`
	Amount    int       `json:"amount"`
	Total     float32   `json:"total"`
	Product   Product   `json:"product"`
}
