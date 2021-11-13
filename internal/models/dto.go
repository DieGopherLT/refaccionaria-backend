package models

import (
	"time"
)

type ProductDTO struct {
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Price       float32 `json:"price"`
	Amount      int     `json:"amount"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	ProviderID  int     `json:"provider_id"`
}

type ProviderDTO struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Enterprise string `json:"enterprise"`
}

type SaleDTO struct {
	ProductID int       `json:"product_id"`
	Date      time.Time `json:"date"`
	Total     float32   `json:"total"`
	Amount    int       `json:"amount"`
}

type DeliveryDTO struct {
	ProductID    int       `json:"product_id"`
	ProviderID   int       `json:"provider_id"`
	DeliveryDate time.Time `json:"delivery_date"`
	Amount       int       `json:"amount"`
}
