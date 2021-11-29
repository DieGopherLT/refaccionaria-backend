package models

import (
	"time"
)

type Product struct {
	ProductID      int      `json:"product_id,omitempty"`
	Classification string   `json:"classification"`
	Brand          string   `json:"brand,omitempty"`
	PublicPrice    float32  `json:"public_price"`
	ProviderPrice  float32  `json:"provider_price"`
	Amount         int      `json:"amount"`
	Category       Category `json:"category,omitempty"`
	Provider       Provider `json:"provider,omitempty"`
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
	Address    string `json:"address,omitempty"`
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
