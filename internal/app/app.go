package app

import (
	"log"
	"net/http"

	"go-proxy/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title HTTP Proxy Server API
// @version 1.0
// @description This is a simple HTTP proxy server.
// @BasePath /
// @host go-proxy-1fo6.onrender.com
// @schemes https
func Run() {
	// Load environment variables from app.env
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

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
