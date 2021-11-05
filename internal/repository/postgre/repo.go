package postgre

import (
	"context"
	"database/sql"
	"time"

	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *sql.DB
}

// InsertProduct inserts a product into database
func (r *Repository) InsertProduct(product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO producto (nombre, categoria_id, marca, precio, cantidad, descripcion)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Category.CategoryID,
		product.Brand,
		product.Price,
		product.Amount,
		product.Description,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetAllProducts fetch all products from databases
func (r *Repository) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	products := []models.Product{}
	query := `
		SELECT
			p.producto_id,
			p.nombre,
			p.marca,
			c.categoria_id,
			c.nombre as categoria,
			p.precio,
			p.cantidad,
			p.descripcion
		FROM
			producto p
			INNER JOIN categoria c ON p.categoria_id = c.categoria_id
			ORDER BY p.producto_id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		product := models.Product{}
		err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.Brand,
			&product.Category.CategoryID,
			&product.Category.Name,
			&product.Price,
			&product.Amount,
			&product.Description,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct updates a product in database
func (r *Repository) UpdateProduct(productID int, product models.Product) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE
			producto
		SET
			nombre = $1,
			marca = $2,
			categoria_id = $3,
			precio = $4,
			cantidad = $5,
			descripcion = $6
		WHERE
			producto_id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Brand,
		product.Category.CategoryID,
		product.Price,
		product.Amount,
		product.Description,
		productID,
	)

	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// DeleteProduct deletes a product from the database
func (r *Repository) DeleteProduct(productID int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM producto WHERE producto_id = $1;`

	result, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// GetAllProviders fetch all providers in database
func (r *Repository) GetAllProviders() ([]models.Provider, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	providers := []models.Provider{}
	query := `SELECT proveedor_id, nombre, correo, telefono, empresa FROM proveedor;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		provider := models.Provider{}
		err := rows.Scan(&provider.ProviderID, &provider.Name, &provider.Email, &provider.Phone, &provider.Enterprise)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
}

// InsertProvider inserts a provider in database
func (r *Repository) InsertProvider(provider models.Provider) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO proveedor (nombre, correo, telefono, empresa)
		VALUES ($1, $2, $3, $4);
	`

	_, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.Email,
		provider.Phone,
		provider.Enterprise,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateProvider(providerID int, provider models.Provider) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE
			proveedor
		SET
			nombre = $1,
			correo = $2,
			telefono = $3,
			empresa = $4
		WHERE
			proveedor_id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.Email,
		provider.Phone,
		provider.Enterprise,
		providerID,
	)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (r *Repository) DeleteProvider(providerID int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM proveedor WHERE proveedor_id = $1`
	result, err := r.db.ExecContext(ctx, query, providerID)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
