package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	_ "go-proxy/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// RequestData represents the JSON structure of the client's request
type RequestData struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

// ResponseData represents the JSON structure of the server's response
type ResponseData struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int               `json:"length"`
}

var (
	requestStore  sync.Map
	responseStore sync.Map
)

// @title HTTP Proxy Server API
// @version 1.0
// @description This is a simple HTTP proxy server.
// @host localhost:8080
// @BasePath /
func main() {
	// Serve Swagger UI and API documentation
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Handle proxy requests
	http.HandleFunc("/proxy", proxyHandler)

	// Start HTTP server
	log.Println("Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// proxyHandler godoc
// @Summary Proxy HTTP request
// @Description Proxy an HTTP request to an external service and return the response.
// @Tags proxy
// @Accept json
// @Produce json
// @Param request body RequestData true "Request Data"
// @Success 200 {object} ResponseData
// @Failure 400 {string} string "Invalid JSON format"
// @Failure 405 {string} string "Only POST method is allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /proxy [post]
// Handler for the "/proxy" route
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData

	// Check the request method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reading the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Request Body: %s\n", body) // Print request body for debugging

	// Parsing the JSON body request
	if err := json.Unmarshal(body, &requestData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err) // Print JSON parsing error for debugging
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	fmt.Printf("Parsed Request Data: %+v\n", requestData) // Print parsed request data for debugging

	// Creating new http request for the external service
	client := &http.Client{}
	req, err := http.NewRequest(requestData.Method, requestData.URL, nil)
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err) // Print request creation error for debugging
		http.Error(w, "Can't create request", http.StatusInternalServerError)
		return
	}

	// Adding headers to the request
	for k, v := range requestData.Headers {
		req.Header.Add(k, v)
	}

	// Executing the request to external service
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request to external service: %v\n", err) // Print request execution error for debugging
		http.Error(w, "Error making request to external service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Reading response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err) // Print response body read error for debugging
		http.Error(w, "Can't read response body", http.StatusInternalServerError)
		return
	}

	// Constructing response data
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		responseHeaders[key] = values[0]
	}

	responseData := ResponseData{
		ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
		Status:  resp.StatusCode,
		Headers: responseHeaders,
		Length:  len(responseBody),
	}

	// Returning JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}
