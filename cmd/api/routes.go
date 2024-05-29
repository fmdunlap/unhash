package main

import "github.com/go-chi/chi/v5"

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)
		r.Route("/hashjob", func(r chi.Router) {
			r.Post("/", app.createHashJobHandler)
			r.Get("/{id}", app.getHashJobHandler)
		})
	})

	return r
}
