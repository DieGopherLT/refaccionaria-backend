package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DieGopherLT/refaccionaria-backend/internal/helpers"
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
	"github.com/DieGopherLT/refaccionaria-backend/internal/repository"
	"github.com/DieGopherLT/refaccionaria-backend/internal/validator"
	"github.com/go-chi/chi/v5"
)

var Repo *Repository

// Repository is a repository that will store all handlers for incoming http requests
type Repository struct {
	db repository.DatabaseRepo
}

// NewHandlersRepo creates a new repository for handlers with a database pool connection
func NewHandlersRepo(db repository.DatabaseRepo) *Repository {
	return &Repository{
		db: db,
	}
}

// SetHandlersRepo sets the repository for handlers that will be used on routes
func SetHandlersRepo(r *Repository) {
	Repo = r
}

// GetProducts handler for get request over product resource
func (m *Repository) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := m.db.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["products"] = products
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

// PostProduct handler for post request over product resource
func (m *Repository) PostProduct(w http.ResponseWriter, r *http.Request) {
	var product models.ProductDTO

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envio en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertProduct(product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	resp := helpers.Response{Message: "Producto creado"}
	helpers.WriteJsonResponse(w, http.StatusCreated, resp)
}

// PutProduct handler for put request over product resource
func (m *Repository) PutProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	product := models.ProductDTO{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateProduct(productId, product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Error al actualizar el producto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado"}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Producto actualizado"}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// DeleteProduct handler for delete request over product resource
func (m *Repository) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteProduct(productId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Error al eliminar el producto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado"}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro eliminado"}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// GetProviders handler for get request over provider resource
func (m *Repository) GetProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := m.db.GetAllProviders()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["providers"] = providers
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

// PostProvider handler for post request over provider resource
func (m *Repository) PostProvider(w http.ResponseWriter, r *http.Request) {
	var newProvider models.ProviderDTO

	err := json.NewDecoder(r.Body).Decode(&newProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	emptyField := validator.HasEmptyStringField(newProvider)
	if emptyField {
		resp := helpers.Response{Message: "Todos los campos son obligatorios", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	isValid, resp := validator.IsValidProvider(newProvider)
	if !isValid {
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertProvider(newProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal con la inserción del registro", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	resp = helpers.Response{Message: "Proveedor registrado correctamente"}
	helpers.WriteJsonResponse(w, http.StatusCreated, resp)
}

// PutProvider handler for put request over provider resource
func (m *Repository) PutProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	var updatedProvider models.ProviderDTO
	err = json.NewDecoder(r.Body).Decode(&updatedProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	emptyField := validator.HasEmptyStringField(updatedProvider)
	if emptyField {
		resp := helpers.Response{Message: "Todos los campos son obligatorios", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	isValid, resp := validator.IsValidProvider(updatedProvider)
	if !isValid {
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateProvider(providerId, updatedProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Teléfono no válido", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp = helpers.Response{Message: "Registro actualizado exitosamente"}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// DeleteProvider handler for delete request over provider resource
func (m *Repository) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteProvider(providerId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Error al eliminar el proveedor", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro eliminado"}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// GetSales handler for get request over sale resource
func (m *Repository) GetSales(w http.ResponseWriter, r *http.Request) {
	sales, err := m.db.GetAllSales()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["sales"] = sales
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

// PostSale handler for post request over sale resource
func (m *Repository) PostSale(w http.ResponseWriter, r *http.Request) {
	var newSale models.SaleDTO

	err := json.NewDecoder(r.Body).Decode(&newSale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	hasEmptyField := validator.HasEmptyStringField(newSale)
	if hasEmptyField {
		resp := helpers.Response{Message: "Todos los campos son obligatorios", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertSale(newSale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	resp := helpers.Response{Message: "Venta agregada exitosamente", Error: false}
	helpers.WriteJsonResponse(w, http.StatusCreated, resp)
}

// PutSale handler for put request over sale resource
func (m *Repository) PutSale(w http.ResponseWriter, r *http.Request) {
	saleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	sale := models.SaleDTO{}
	err = json.NewDecoder(r.Body).Decode(&sale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateSale(saleId, sale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		fmt.Println(err)
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro actualizado extitosamente", Error: false}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// DeleteSale handler for delete request over sale resource
func (m *Repository) DeleteSale(w http.ResponseWriter, r *http.Request) {
	saleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteSale(saleId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		fmt.Println(err)
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro eliminado exitosamente", Error: false}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// GetDeliveries handler for get request over delivery resource
func (m *Repository) GetDeliveries(w http.ResponseWriter, r *http.Request) {
	deliveries, err := m.db.GetAllDeliveries()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusOK, resp)
		return
	}
	data := make(map[string]interface{})
	data["deliveries"] = deliveries
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

// PostDelivery handler for post request over delivery resource
func (m *Repository) PostDelivery(w http.ResponseWriter, r *http.Request) {
	var deliveryDTO models.DeliveryDTO

	err := json.NewDecoder(r.Body).Decode(&deliveryDTO)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.InsertDelivery(deliveryDTO)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "No se encontró el producto o proveedor", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Entrega registrada correctamente", Error: false}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// DeleteDelivery handler for delete request over delivery resource
func (m *Repository) DeleteDelivery(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("productId"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	providerId, err := strconv.Atoi(r.URL.Query().Get("providerId"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteDelivery(productId, providerId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "No se encontró la entrega", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Entrega dada de alta", Error: true}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// GetClients handler for get request over client resource
func (m *Repository) GetClients(w http.ResponseWriter, r *http.Request) {
	clients, err := m.db.GetAllClients()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["clients"] = clients
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

// PostClient handler for post request over client resource
func (m *Repository) PostClient(w http.ResponseWriter, r *http.Request) {
	var client models.ClientDTO

	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información fue enviada en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	isValid, resp := validator.IsValidClient(client)
	if !isValid {
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertClient(client)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	resp = helpers.Response{Message: "Cliente registrado exitosamente", Error: false}
	helpers.WriteJsonResponse(w, http.StatusCreated, resp)
}

// PutClient handler for put request over client resource
func (m *Repository) PutClient(w http.ResponseWriter, r *http.Request) {
	var client models.ClientDTO

	clientId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	isValid, resp := validator.IsValidClient(client)
	if !isValid {
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateClient(clientId, client)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Cliente no encontrado", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp = helpers.Response{Message: "Cliente actualizado exitosamente", Error: false}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// DeleteClient handler for delete request over client resource
func (m *Repository) DeleteClient(w http.ResponseWriter, r *http.Request) {
	clientId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonResponse(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteClient(clientId)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		fmt.Println(err)
		resp := helpers.Response{Message: "Cliente no encontrado", Error: true}
		helpers.WriteJsonResponse(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Cliente dado de baja", Error: false}
	helpers.WriteJsonResponse(w, http.StatusOK, resp)
}

// GetBrands handler for get request over brand resource
func (m *Repository) GetBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := m.db.GetAllBrands()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["brands"] = brands
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

// GetCategories handler for get request over category resource
func (m *Repository) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := m.db.GetAllCategories()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonResponse(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["categories"] = categories
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}
