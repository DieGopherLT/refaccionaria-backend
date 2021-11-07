package models

import (
	"time"
)

type Product struct {
	ProductID   int      `json:"product_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Price       float32  `json:"price,omitempty"`
	Amount      int      `json:"amount,omitempty"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Provider    Provider `json:"provider,omitempty"`
}

type Category struct {
	CategoryID int    `json:"category_id,omitempty"`
	Name       string `json:"name,omitempty"`
}

type Provider struct {
	ProviderID int    `json:"provider_id,omitempty"`
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Enterprise string `json:"enterprise,omitempty"`
}

type ProductProvider struct {
	DeliveryDate time.Time `json:"delivery_date,omitempty"`
	ProductID    int       `json:"product_id,omitempty"`
	ProviderID   int       `json:"provider_id,omitempty"`
	Product      Product   `json:"product,omitempty"`
	Provider     Provider  `json:"provider,omitempty"`
}

type Sale struct {
	SaleID  int       `json:"sale_id,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	Amount  int       `json:"amount,omitempty"`
	Total   float32   `json:"total,omitempty"`
	Product Product   `json:"product,omitempty"`
}
