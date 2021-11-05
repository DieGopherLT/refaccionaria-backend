package repository

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
)

type DatabaseRepo interface {
	/*
		These are the functions for the product CRUD
	*/
	InsertProduct(product models.Product) error
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(productID int, product models.Product) (int64, error)
	DeleteProduct(productID int) (int64, error)

	/*
		These are the functions for the provider CRUD
	*/
	GetAllProviders() ([]models.Provider, error)
	InsertProvider(provider models.Provider) error
	UpdateProvider(providerId int, provider models.Provider) (int64, error)
	DeleteProvider(providerID int) (int64, error)
}
