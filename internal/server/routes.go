package server

import "github.com/go-chi/chi/v5"

// registerRoutes registers all routes for the application
func registerRoutes(r *chi.Mux, dependencies Dependencies) {
	r.Group(unauthenticatedRoutes(r, dependencies))
	r.Group(publicRoutes(r, dependencies))
	r.Group(protectedRoutes(r, dependencies))
	r.Group(adminRoutes(r, dependencies))
}

// unauthenticatedRoutes defines routes that will not attempt authentication
func unauthenticatedRoutes(r chi.Router, dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		// System Checks
		r.Get("/liveness", dependencies.SysHandler.Liveness)
		r.Get("/readiness", dependencies.SysHandler.Readiness)

		// Auth
		r.Post("/authbus/register", dependencies.AuthHandler.Register)
		r.Post("/authbus/login", dependencies.AuthHandler.Login)
		r.Post("/authbus/refresh", dependencies.AuthHandler.Refresh)
	}
}

// publicRoutes defines routes that do not require authentication
func publicRoutes(r chi.Router, dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(dependencies.AuthMiddleware.Authenticate)
	}
}

// protectedRoutes defines routes that require authentication
func protectedRoutes(r chi.Router, dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(dependencies.AuthMiddleware.Authenticate)
		r.Use(dependencies.AuthMiddleware.ForceAuthentication)

		// Auth
		r.Post("/authbus/logout", dependencies.AuthHandler.Logout)
	}
}

// adminRoutes defines routes that
func adminRoutes(r chi.Router, dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(dependencies.AuthMiddleware.Authenticate)
		r.Use(dependencies.AuthMiddleware.ForceAdmin)
	}
}
