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
)

var Repo *Repository

type Repository struct {
	db repository.DatabaseRepo
}

func NewRepo(db repository.DatabaseRepo) *Repository {
	return &Repository{
		db: db,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// GetProducts handler for get request over product resource
func (m *Repository) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := m.db.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
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
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertProduct(product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	resp := helpers.Response{Message: "Producto creado"}
	helpers.WriteJsonMessage(w, http.StatusCreated, resp)
}

// PutProduct handler for put request over product resource
func (m *Repository) PutProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	product := models.ProductDTO{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateProduct(productId, product)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Error al actualizar el producto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado"}
		helpers.WriteJsonMessage(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Producto actualizado"}
	helpers.WriteJsonMessage(w, http.StatusOK, resp)
}

// DeleteProduct handler for delete request over product resource
func (m *Repository) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteProduct(productId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Error al eliminar el producto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado"}
		helpers.WriteJsonMessage(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro eliminado"}
	helpers.WriteJsonMessage(w, http.StatusOK, resp)
}

// GetProviders handler for get request over provider resource
func (m *Repository) GetProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := m.db.GetAllProviders()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
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
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	emptyField := validator.HasEmptyStringField(newProvider)
	if emptyField {
		resp := helpers.Response{Message: "Todos los campos son obligatorios", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	isValid, resp := validator.IsValidProvider(newProvider)
	if !isValid {
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertProvider(newProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal con la inserción del registro", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	resp = helpers.Response{Message: "Proveedor registrado correctamente"}
	helpers.WriteJsonMessage(w, http.StatusCreated, resp)
}

// PutProvider handler for put request over provider resource
func (m *Repository) PutProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	var updatedProvider models.ProviderDTO
	err = json.NewDecoder(r.Body).Decode(&updatedProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	emptyField := validator.HasEmptyStringField(updatedProvider)
	if emptyField {
		resp := helpers.Response{Message: "Todos los campos son obligatorios", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	isValid, resp := validator.IsValidProvider(updatedProvider)
	if !isValid {
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateProvider(providerId, updatedProvider)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Teléfono no válido", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonMessage(w, http.StatusNotFound, resp)
		return
	}

	resp = helpers.Response{Message: "Registro actualizado exitosamente"}
	helpers.WriteJsonMessage(w, http.StatusOK, resp)
}

// DeleteProvider handler for delete request over provider resource
func (m *Repository) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	providerId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteProvider(providerId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Error al eliminar el proveedor", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonMessage(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro eliminado"}
	helpers.WriteJsonMessage(w, http.StatusOK, resp)
}

// GetSales handler for get request over sale resource
func (m *Repository) GetSales(w http.ResponseWriter, r *http.Request) {
	sales, err := m.db.GetAllSales()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
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
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	hasEmptyField := validator.HasEmptyStringField(newSale)
	if hasEmptyField {
		resp := helpers.Response{Message: "Todos los campos son obligatorios", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	err = m.db.InsertSale(newSale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	resp := helpers.Response{Message: "Venta agregada exitosamente", Error: false}
	helpers.WriteJsonMessage(w, http.StatusCreated, resp)
}

// PutSale handler for put request over sale resource
func (m *Repository) PutSale(w http.ResponseWriter, r *http.Request) {
	saleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	sale := models.SaleDTO{}
	err = json.NewDecoder(r.Body).Decode(&sale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateSale(saleId, sale)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		fmt.Println(err)
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonMessage(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro actualizado extitosamente", Error: false}
	helpers.WriteJsonMessage(w, http.StatusOK, resp)
}

// DeleteSale handler for delete over sale resource
func (m *Repository) DeleteSale(w http.ResponseWriter, r *http.Request) {
	saleId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "La información se envió en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteSale(saleId)
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salió mal...", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	if rows == 0 {
		fmt.Println(err)
		resp := helpers.Response{Message: "Registro no encontrado", Error: true}
		helpers.WriteJsonMessage(w, http.StatusNotFound, resp)
		return
	}

	resp := helpers.Response{Message: "Registro eliminado exitosamente", Error: false}
	helpers.WriteJsonMessage(w, http.StatusOK, resp)
}

// GetBrands handler for get request over brand resource
func (m *Repository) GetBrands(w http.ResponseWriter, r *http.Request) {
	brands, err := m.db.GetAllBrands()
	if err != nil {
		fmt.Println(err)
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["brands"] = brands
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}
