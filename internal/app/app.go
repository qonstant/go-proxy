package app

import (
	"log"
	"net/http"

	_ "go-proxy/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"go-proxy/internal/handler"
)

func Run() {
	// Create a new chi router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve Swagger UI and API documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Handle proxy requests
	r.Post("/proxy", handler.ProxyHandler)

	// Start HTTP server
	log.Println("Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
