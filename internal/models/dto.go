package models

import (
	"time"
)

type ProductDTO struct {
	Classification string  `json:"classification"`
	Brand          string  `json:"brand"`
	PublicPrice    float32 `json:"public_price"`
	ProviderPrice  float32 `json:"provider_price"`
	Amount         int     `json:"amount,omitempty"`
	CategoryID     int     `json:"category_id"`
	ProviderID     int     `json:"provider_id"`
}

type ProviderDTO struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Enterprise string `json:"enterprise"`
	Address    string `json:"address"`
}

type SaleDTO struct {
	ProductID int       `json:"product_id"`
	ClientID  int       `json:"client_id"`
	Date      time.Time `json:"date"`
	Total     float32   `json:"total"`
	Subtotal  float32   `json:"subtotal"`
	Amount    int       `json:"amount"`
}

type DeliveryDTO struct {
	ProductID    int    `json:"product_id"`
	ProviderID   int    `json:"provider_id"`
	DeliveryDate string `json:"delivery_date"`
	Amount       int    `json:"amount"`
}

type ClientDTO struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
}
