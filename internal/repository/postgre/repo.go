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
		INSERT INTO producto (nombre_producto, id_categoria, marca, precio_publico, stock, descripcion)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_producto;
	`

	var newID int
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.CategoryID,
		product.Brand,
		product.Price,
		product.Amount,
		product.Description,
	).Scan(&newID)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO producto_proveedor (id_producto, id_proveedor, fecha_entrega, cantidad_surtir)
		VALUES ($1, $2, NULL, NULL);
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
			p.id_producto,
			p.nombre_producto,
			p.marca,
			p.descripcion,
			p.precio_publico,
			p.stock,
			c.id_categoria,
			c.nombre_categoria as categoria,
			pr.codigo,
		    pr.nombre_proveedor,
			pr.correo as correo_proveedor,
			pr.telefono_proveedor as tel_proveedor
		FROM producto p
		INNER JOIN categoria c
			ON c.id_categoria = p.id_categoria
		INNER JOIN producto_proveedor pp
			ON pp.id_producto = p.id_producto
		INNER JOIN proveedor pr
			ON pp.id_proveedor = pr.codigo;
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

		/*
			From GUI product stock is not directly modified, but due to a database trigger, stock column needs to receive
			a value not equal to NULL not to activate the trigger.
		*/
		query := `
			UPDATE
				producto
			SET
				nombre_producto = $1,
				marca = $2,
				id_categoria = $3,
				precio_publico = $4,
			    stock = $5,
				descripcion = $6
			WHERE
				id_producto = $7;
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

		query := `UPDATE producto_proveedor SET id_proveedor = $1 WHERE id_producto = $2`
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

	query := `DELETE FROM producto WHERE id_producto = $1;`

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
	query := `
		SELECT 
		   codigo, 
		   nombre_proveedor, 
		   correo, 
		   telefono_proveedor, 
		   empresa,
		   direccion_proveedor,
		   precio_proveedor
		FROM 
		     proveedor;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		provider := models.Provider{}
		err := rows.Scan(
			&provider.ProviderID,
			&provider.Name,
			&provider.Email,
			&provider.Phone,
			&provider.Enterprise,
			&provider.Address,
			&provider.ProviderPrice,
		)
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
		INSERT INTO proveedor 
		    (nombre_proveedor, correo, telefono_proveedor, empresa, direccion_proveedor, precio_proveedor)
		VALUES 
		       ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.Email,
		provider.Phone,
		provider.Enterprise,
		provider.Address,
		provider.ProviderPrice,
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
			nombre_proveedor = $1,
			correo = $2,
			telefono_proveedor = $3,
			empresa = $4,
		    direccion_proveedor = $5,
		    precio_proveedor = $6
		WHERE
			codigo = $7;
	`

	result, err := r.db.ExecContext(ctx, query,
		provider.Name,
		provider.Email,
		provider.Phone,
		provider.Enterprise,
		provider.Address,
		provider.ProviderPrice,
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

	query := `DELETE FROM proveedor WHERE codigo = $1`
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
			v.id_venta,
			v.fecha,
			v.total,
			v.cantidad_vendida,
			p.id_producto,
			p.nombre_producto,
			p.marca,
		    p.precio_publico
		FROM venta v
		INNER JOIN producto p 
			ON v.id_producto = p.id_producto;
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
		INSERT INTO venta (id_producto, fecha, total, cantidad_vendida)
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
		SET id_producto = $1, total = $2, cantidad_vendida = $3
		WHERE id_venta = $4;
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

	query := `DELETE FROM venta WHERE id_venta = $1`

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
		    po.id_producto,
			po.nombre_producto,
			po.marca,
		    pr.codigo,
			pr.nombre_proveedor,
			pr.correo,
			pp.fecha_entrega,
			pp.cantidad_surtir
		FROM 
			producto_proveedor pp
		INNER JOIN producto po
			ON pp.id_producto = po.id_producto
		INNER JOIN proveedor pr
			ON pp.id_proveedor = pr.codigo
		WHERE pp.fecha_entrega IS NOT NULL;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		d := models.Delivery{}
		err := rows.Scan(
			&d.Product.ProductID, &d.Product.Name, &d.Product.Brand,
			&d.Provider.ProviderID, &d.Provider.Name, &d.Provider.Email,
			&d.DeliveryDate, &d.Amount,
		)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, d)
	}

	if err := rows.Err(); err != nil {
		return deliveries, err
	}

	return deliveries, nil
}

// InsertDelivery inserts a delivery in database
func (r *Repository) InsertDelivery(delivery models.DeliveryDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		UPDATE producto_proveedor
		SET fecha_entrega = $1, cantidad_surtir = $2
		WHERE id_producto = $3 AND id_proveedor = $4;
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
		SET fecha_entrega = NULL, cantidad_surtir = NULL
		WHERE id_producto = $1 AND id_proveedor = $2; 
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

// GetAllClients fetches all clients from database
func (r *Repository) GetAllClients() ([]models.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	clients := []models.Client{}
	query := `
		 SELECT 
			id_cliente,
		    nombre_cliente,
		    direccion_cliente,
		    telefono_cliente
		FROM
			cliente;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := models.Client{}
		err := rows.Scan(&c.ClientID, &c.Name, &c.Address, &c.Phone)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

// InsertClient inserts a client into database
func (r *Repository) InsertClient(client models.ClientDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
		INSERT INTO cliente (nombre_cliente, telefono_cliente, direccion_cliente)
		VALUES ($1, $2, $3);
`
	_, err := r.db.ExecContext(ctx, query, client.Name, client.Phone, client.Address)
	if err != nil {
		return err
	}

	return nil
}

// DeleteClient deletes a client in database
func (r *Repository) DeleteClient(clientId int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `DELETE FROM cliente WHERE id_cliente = $1;`
	result, err := r.db.ExecContext(ctx, query, clientId)
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
	query := `SELECT id_categoria, nombre_categoria FROM categoria ORDER BY id_categoria;`

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
