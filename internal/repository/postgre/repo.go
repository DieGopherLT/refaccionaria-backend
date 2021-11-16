package postgre

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
)

func NewRepository(pool *sql.DB) *Repository {
	return &Repository{
		db: pool,
	}
}

type Repository struct {
	db *sql.DB
}

// InsertProduct inserts a product into database
func (r *Repository) InsertProduct(product models.ProductDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO producto (nombre, categoria_id, marca, precio, cantidad, descripcion)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	var newID int
	_, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.CategoryID,
		product.Brand,
		product.Price,
		product.Amount,
		product.Description,
	)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO producto_proveedor (producto_id, proveedor_id, fecha_entrega)
		VALUES ($1, $2, NULL);
	`
	_, err = r.db.ExecContext(ctx, query, newID, product.ProviderID)
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
			p.descripcion,
			p.precio,
			p.cantidad,
			c.categoria_id,
			c.nombre as categoria,
			pr.proveedor_id,
		    pr.nombre,
			pr.correo as correo_proveedor,
			pr.telefono as tel_proveedor
		FROM producto p
		INNER JOIN categoria c
			ON c.categoria_id = p.categoria_id
		INNER JOIN producto_proveedor pp
			ON pp.producto_id = p.producto_id
		INNER JOIN proveedor pr
			ON pp.proveedor_id = pr.proveedor_id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := models.Product{}
		err := rows.Scan(
			&p.ProductID, &p.Name, &p.Brand, &p.Description, &p.Price, &p.Amount,
			&p.Category.CategoryID, &p.Category.Name,
			&p.Provider.ProviderID, &p.Provider.Name, &p.Provider.Email, &p.Provider.Phone,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct updates a product in database
func (r *Repository) UpdateProduct(productID int, product models.ProductDTO) (int64, error) {
	wg := sync.WaitGroup{}
	errChan := make(chan error)
	resultChan := make(chan sql.Result)
	defer close(errChan)
	defer close(resultChan)

	wg.Add(1)

	go func() {
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
				producto_id = $7;
		`
		result, err := r.db.ExecContext(ctx, query,
			product.Name,
			product.Brand,
			product.CategoryID,
			product.Price,
			product.Amount,
			product.Description,
			productID,
		)
		wg.Wait()
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		query := `UPDATE producto_proveedor SET proveedor_id = $1 WHERE producto_id = $2`
		_, err := r.db.ExecContext(ctx, query, product.ProviderID, productID)
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()

	var numRows int64
	select {
	case err := <-errChan:
		return 0, err
	case result := <-resultChan:
		rows, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		numRows = rows
		break
	}

	if numRows == 0 {
		return 0, nil
	}

	return numRows, nil
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
func (r *Repository) InsertProvider(provider models.ProviderDTO) error {
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

// UpdateProvider updates a provider in database
func (r *Repository) UpdateProvider(providerID int, provider models.ProviderDTO) (int64, error) {
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

// DeleteProvider deletes a provider in database
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

// GetAllSales fetches all sales stored in database
func (r *Repository) GetAllSales() ([]models.Sale, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sales := []models.Sale{}
	query := `
		SELECT 
			v.venta_id,
			v.fecha,
			v.total,
			v.cantidad,
			p.producto_id,
			p.nombre,
			p.marca,
		    p.precio
		FROM venta v
		INNER JOIN producto p 
			ON v.producto_id = p.producto_id;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := models.Sale{}
		err := rows.Scan(
			&s.SaleID, &s.Date, &s.Total, &s.Amount,
			&s.Product.ProductID, &s.Product.Name, &s.Product.Brand, &s.Product.Price,
		)
		if err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sales, nil
}

// InsertSale inserts a sale in database
func (r *Repository) InsertSale(sale models.SaleDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO venta (producto_id, fecha, total, cantidad)
		VALUES ($1, CURRENT_DATE, $2, $3);
	`

	_, err := r.db.ExecContext(ctx, query, sale.ProductID, sale.Total, sale.Amount)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSale updates a sale in database
func (r *Repository) UpdateSale(saleId int, sale models.SaleDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE venta
		SET producto_id = $1, total = $2, cantidad = $3
		WHERE venta_id = $4;
	`

	result, err := r.db.ExecContext(ctx, query, sale.ProductID, sale.Total, sale.Amount, saleId)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// DeleteSale deletes a sale in database
func (r *Repository) DeleteSale(saleId int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM venta WHERE venta_id = $1`

	result, err := r.db.ExecContext(ctx, query, saleId)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// GetAllDeliveries brings all the deliveries from database
func (r *Repository) GetAllDeliveries() ([]models.Delivery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	deliveries := []models.Delivery{}
	query := `
		SELECT 
			po.nombre,
			po.marca,
			pr.nombre,
			pr.correo,
			pp.fecha_entrega,
			pp.cantidad
		FROM 
			producto_proveedor pp
		INNER JOIN producto po
			ON pp.producto_id = po.producto_id
		INNER JOIN proveedor pr
			ON pp.proveedor_id = pr.proveedor_id
		WHERE pp.fecha_entrega IS NOT NULL;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		d := models.Delivery{}
		err := rows.Scan(
			&d.Product.Name, &d.Product.Brand,
			&d.Provider.Name, &d.Provider.Email,
			&d.DeliveryDate, &d.Amount,
		)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}

	return deliveries, nil
}

// InsertDelivery inserts a delivery in database
func (r *Repository) InsertDelivery(delivery models.DeliveryDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE producto_proveedor
		SET fecha_entrega = $1, cantidad = $2
		WHERE producto_id = $3 AND proveedor_id = $4;
	`

	result, err := r.db.ExecContext(ctx, query,
		delivery.DeliveryDate,
		delivery.Amount,
		delivery.ProductID,
		delivery.ProviderID,
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

// DeleteDelivery "deletes" a delivery in frontend perspective, it just updates some fields to NULL
func (r *Repository) DeleteDelivery(productID, providerID int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE producto_proveedor
		SET fecha_entrega = NULL, cantidad = NULL
		WHERE producto_id = $1 AND proveedor_id = $2; 
	`

	result, err := r.db.ExecContext(ctx, query, productID, providerID)
	if err != nil {
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

// GetAllBrands brings all the brands from providers from database without duplicates
func (r *Repository) GetAllBrands() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	brands := []string{}
	query := `SELECT DISTINCT empresa FROM proveedor;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		brand := ""
		err := rows.Scan(&brand)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return brands, nil
}

// GetAllCategories fetches all categories from database
func (r *Repository) GetAllCategories() ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	categories := []models.Category{}
	query := `SELECT categoria_id, nombre FROM categoria;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		category := models.Category{}
		err := rows.Scan(&category.CategoryID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
