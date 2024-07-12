package handler

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "strings"
    "sync"

    "go-proxy/internal/model"

    "github.com/google/uuid"
)

var (
    client        = &http.Client{}
    requestStore  sync.Map
    responseStore sync.Map
)

// ProxyHandler godoc
// @Summary Proxy HTTP request
// @Description Proxy an HTTP request to an external service and return the response.
// @Tags proxy
// @Accept json
// @Produce json
// @Param request body model.RequestData true "Request Data"
// @Success 200 {object} model.ResponseData
// @Failure 400 {string} string "Invalid JSON format"
// @Failure 405 {string} string "Only POST method is allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /proxy [post]
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
    var requestData model.RequestData

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

    // Parse the provided URL and update the host and scheme from environment variables
    parsedURL, err := url.Parse(requestData.URL)
    if err != nil {
        log.Printf("Error parsing URL: %v", err)
        http.Error(w, "Invalid URL format", http.StatusBadRequest)
        return
    }
    parsedURL.Host = os.Getenv("PROXY_HOST")
    parsedURL.Scheme = os.Getenv("PROXY_SCHEME")
    requestData.URL = parsedURL.String()

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
    responseData := model.ResponseData{
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
