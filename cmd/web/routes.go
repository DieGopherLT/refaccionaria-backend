package main

import (
	"net/http"
	"time"

	"github.com/DieGopherLT/refaccionaria-backend/internal/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func Routes() http.Handler {
	mux := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(5 * time.Second))

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		message := []byte("Server is working")
		w.Write(message)
	})

	mux.Route("/api", func(r chi.Router) {

		r.Route("/v1", func(r chi.Router) {
			r.Route("/product", func(r chi.Router) {
				r.Get("/", controller.Repo.GetProducts)
				r.Post("/", controller.Repo.PostProduct)
				r.Put("/", controller.Repo.PutProduct)
				r.Delete("/", controller.Repo.DeleteProduct)
			})

			r.Route("/provider", func(r chi.Router) {
				r.Get("/", controller.Repo.GetProviders)
				r.Post("/", controller.Repo.PostProvider)
				r.Put("/", controller.Repo.PutProvider)
				r.Delete("/", controller.Repo.DeleteProvider)
			})

			r.Route("/sale", func(r chi.Router) {
				r.Get("/", controller.Repo.GetSales)
				r.Post("/", controller.Repo.PostSale)
				r.Put("/", controller.Repo.PutSale)
				r.Delete("/", controller.Repo.DeleteSale)
			})

			r.Route("/delivery", func(r chi.Router) {
				r.Get("/", controller.Repo.GetDeliveries)
				r.Post("/", controller.Repo.PostDelivery)
				r.Delete("/", controller.Repo.DeleteDelivery)
			})

			r.Route("/client", func(r chi.Router) {
				r.Get("/", controller.Repo.GetClients)
				r.Post("/", controller.Repo.PostClient)
				r.Put("/{id}", controller.Repo.PutClient)
				r.Delete("/{id}", controller.Repo.DeleteClient)
			})

			r.Route("/brand", func(r chi.Router) {
				r.Get("/", controller.Repo.GetBrands)
			})

			r.Route("/category", func(r chi.Router) {
				r.Get("/", controller.Repo.GetCategories)
			})
		})

	})

	return c.Handler(mux)
}
