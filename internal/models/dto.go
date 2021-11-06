package models

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
