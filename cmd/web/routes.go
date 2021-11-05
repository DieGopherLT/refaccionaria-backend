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
	mux.Use(middleware.Timeout(5 * time.Second))

	mux.Route("/api/v1", func(r chi.Router) {

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
	})

	return mux
}
