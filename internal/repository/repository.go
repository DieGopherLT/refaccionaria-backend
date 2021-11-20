package repository

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
)

type DatabaseRepo interface {
	InsertProduct(product models.ProductDTO) error
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(productID int, product models.ProductDTO) (int64, error)
	DeleteProduct(productID int) (int64, error)

	GetAllProviders() ([]models.Provider, error)
	InsertProvider(provider models.ProviderDTO) error
	UpdateProvider(providerId int, provider models.ProviderDTO) (int64, error)
	DeleteProvider(providerID int) (int64, error)

	GetAllSales() ([]models.Sale, error)
	InsertSale(sale models.SaleDTO) error
	UpdateSale(saleId int, sale models.SaleDTO) (int64, error)
	DeleteSale(saleId int) (int64, error)

	GetAllDeliveries() ([]models.Delivery, error)
	InsertDelivery(delivery models.DeliveryDTO) (int64, error)
	DeleteDelivery(productID, providerID int) (int64, error)

	GetAllClients() ([]models.Client, error)
	InsertClient(client models.ClientDTO) error
	DeleteClient(clientId int) (int64, error)

	GetAllBrands() ([]string, error)

	GetAllCategories() ([]models.Category, error)
}
