package main

import (
	"net/http"
	"time"

	"github.com/DieGopherLT/refaccionaria-backend/internal/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(3 * time.Second))

	mux.Route("/api/v1", func(r chi.Router) {

		r.Route("/product", func(r chi.Router) {
			r.Get("/", controller.Repo.GetAllProducts)
			r.Post("/", controller.Repo.AddProduct)
			r.Put("/", controller.Repo.UpdateProduct)
			r.Delete("/", controller.Repo.DeleteProduct)
		})
	})

	return mux
}
