package repository

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
)

type DatabaseRepo interface {
	InsertProduct(product models.Product) error
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(productID int, product models.Product) (int64, error)
	DeleteProduct(productID int) (int64, error)
}
