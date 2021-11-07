package repository

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
)

type DatabaseRepo interface {
	/*
		These are the functions for the product CRUD
	*/
	InsertProduct(product models.ProductDTO) error
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(productID int, product models.ProductDTO) (int64, error)
	DeleteProduct(productID int) (int64, error)

	/*
		These are the functions for the provider CRUD
	*/
	GetAllProviders() ([]models.Provider, error)
	InsertProvider(provider models.ProviderDTO) error
	UpdateProvider(providerId int, provider models.ProviderDTO) (int64, error)
	DeleteProvider(providerID int) (int64, error)

	/*
		These are the functions for the sale CRUD
	*/
	GetAllSales() ([]models.Sale, error)
	InsertSale(sale models.SaleDTO) error
	UpdateSale(saleId int, sale models.SaleDTO) (int64, error)
	DeleteSale(saleId int) (int64, error)
}
