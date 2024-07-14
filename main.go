package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "strings"
    "sync"

    _ "go-proxy/docs"
    "go-proxy/internal/models"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/google/uuid"
    httpSwagger "github.com/swaggo/http-swagger"
)

var (
    client        = &http.Client{}
    requestStore  sync.Map
    responseStore sync.Map
)

// @title HTTP Proxy Server API
// @version 1.0
// @description This is a simple HTTP proxy server.
// @BasePath /
// @host localhost:8080
// 
// @schemes http https
func main() {
    // Create a new chi router
    r := chi.NewRouter()

    // Add middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // CORS configuration
    cors := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    })
    r.Use(cors.Handler)

    // Serve Swagger UI and API documentation
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    // Handle proxy requests
    r.Post("/proxy", proxyHandler)

    // Start HTTP server
    log.Println("Starting HTTP server on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}

// proxyHandler godoc
// @Summary Proxy HTTP request
// @Description Proxy an HTTP request to an external service and return the response.
// @Tags proxy
// @Accept json
// @Produce json
// @Param request body models.RequestData true "Request Data"
// @Success 200 {object} models.ResponseData
// @Failure 400 {string} string "Invalid JSON format"
// @Failure 405 {string} string "Only POST method is allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /proxy [post]
func proxyHandler(w http.ResponseWriter, r *http.Request) {
    var requestData models.RequestData

    // Check the request method
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Reading the request body
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading request body: %v", err)
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // Parsing the JSON body request
    if err := json.Unmarshal(body, &requestData); err != nil {
        log.Printf("Error unmarshaling request body: %v", err)
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    // Creating new http request for the external service
    req, err := http.NewRequest(requestData.Method, requestData.URL, strings.NewReader(requestData.Body))
    if err != nil {
        log.Printf("Error creating new request: %v", err)
        http.Error(w, "Can't create request", http.StatusInternalServerError)
        return
    }

    // Adding headers to the request
    for k, v := range requestData.Headers {
        req.Header.Add(k, v)
    }

    // Setting Content-Type if provided
    if contentType, exists := requestData.Headers["Content-Type"]; exists {
        req.Header.Set("Content-Type", contentType)
    }

    // Executing the request to external service
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making request to external service: %v", err)
        http.Error(w, "Error making request to external service", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Reading response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        http.Error(w, "Can't read response body", http.StatusInternalServerError)
        return
    }

    // Constructing response data
    responseHeaders := make(map[string]string)
    for key, values := range resp.Header {
        responseHeaders[key] = strings.Join(values, ", ")
    }

    responseID := uuid.New().String()
    responseData := models.ResponseData{
        ID:      responseID,
        Status:  resp.StatusCode,
        Headers: responseHeaders,
        Length:  len(responseBody),
        Body:    string(responseBody),
    }

    // Save request and response in sync.Map
    requestStore.Store(responseID, requestData)
    responseStore.Store(responseID, responseData)

    // Returning JSON response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(responseData); err != nil {
        log.Printf("Error encoding response data: %v", err)
    }
}
