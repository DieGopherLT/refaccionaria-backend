package models

import (
	"time"
)

type ProductDTO struct {
	Name        string  `json:"name"`
	Brand       string  `json:"brand"`
	Price       float32 `json:"price"`
	Amount      int     `json:"amount,omitempty"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	ProviderID  int     `json:"provider_id"`
}

type ProviderDTO struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Enterprise    string `json:"enterprise"`
	Address       string `json:"address"`
	ProviderPrice int    `json:"provider_price"`
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
