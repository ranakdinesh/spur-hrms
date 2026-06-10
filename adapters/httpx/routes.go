package httpx

import (
	"github.com/go-chi/chi/v5"
	"y/adapters/httpx/handlers"
)

// RegisterRoutes mounts Hrms HTTP routes.
// Route prefix ("/hrms") is set by the application in app.go.
func RegisterRoutes(r chi.Router, h *handlers.Handler) {
	r.Route("/hrms", func(r chi.Router) {
		// Public routes (if any) go before the auth group
		// r.Get("/public", h.PublicEndpoint)

		// Protected routes (JWT required — applied by platform middleware)
		r.Get("/",        h.List)
		r.Post("/",       h.Create)
		r.Get("/{id}",   h.Get)
		r.Delete("/{id}", h.Delete)
	})
}
