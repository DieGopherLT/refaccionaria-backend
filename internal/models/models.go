package models

import (
	"time"
)

type Product struct {
	ProductID   int      `json:"product_id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Price       float32  `json:"price,omitempty"`
	Amount      int      `json:"amount"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	Provider    Provider `json:"provider,omitempty"`
}

type Category struct {
	CategoryID int    `json:"category_id,omitempty"`
	Name       string `json:"name,omitempty"`
}

type Provider struct {
	ProviderID    int     `json:"provider_id,omitempty"`
	Email         string  `json:"email,omitempty"`
	Name          string  `json:"name,omitempty"`
	Phone         string  `json:"phone,omitempty"`
	Enterprise    string  `json:"enterprise,omitempty"`
	Address       string  `json:"address,omitempty"`
	ProviderPrice float32 `json:"provider_price,omitempty"`
}

type Delivery struct {
	DeliveryDate time.Time `json:"delivery_date,omitempty"`
	Product      Product   `json:"product,omitempty"`
	Provider     Provider  `json:"provider,omitempty"`
	Amount       int       `json:"amount,omitempty"`
}

type Sale struct {
	SaleID   int       `json:"sale_id,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Amount   int       `json:"amount"`
	Total    float32   `json:"total,omitempty"`
	SubTotal float32   `json:"sub_total,omitempty"`
	Product  Product   `json:"product,omitempty"`
}

type Client struct {
	ClientID int `json:"client_id,omitempty"`
	ClientDTO
}
