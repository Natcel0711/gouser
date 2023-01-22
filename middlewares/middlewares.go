package middlewares

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func UseMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
}
