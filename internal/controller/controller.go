package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DieGopherLT/refaccionaria-backend/internal/helpers"
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
	"github.com/DieGopherLT/refaccionaria-backend/internal/repository"
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

func (m *Repository) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = m.db.InsertProduct(product)
	if err != nil {
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}

	resp := helpers.Response{Message: "Producto creado"}
	helpers.WriteJsonMessage(w, http.StatusCreated, resp)
}

func (m *Repository) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := m.db.GetAllProducts()
	if err != nil {
		resp := helpers.Response{Message: "Algo salio mal", Error: true}
		helpers.WriteJsonMessage(w, http.StatusInternalServerError, resp)
		return
	}
	data := make(map[string]interface{})
	data["products"] = products
	data["error"] = false
	helpers.WriteJsonResponse(w, http.StatusOK, data)
}

func (m *Repository) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	product := models.Product{}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.UpdateProduct(productId, product)
	if err != nil {
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

func (m *Repository) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		resp := helpers.Response{Message: "La información se envío en un formato incorrecto", Error: true}
		helpers.WriteJsonMessage(w, http.StatusBadRequest, resp)
		return
	}

	rows, err := m.db.DeleteProduct(productId)
	if err != nil {
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
