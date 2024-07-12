package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Test the proxyHandler
func TestProxyHandler(t *testing.T) {
	// Test case: Valid request
	t.Run("valid request", func(t *testing.T) {
		requestData := RequestData{
			Method: "GET",
			URL:    "http://google.com",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		requestBody, _ := json.Marshal(requestData)
		req := httptest.NewRequest(http.MethodPost, "/proxy", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(proxyHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

		var responseData ResponseData
		err := json.NewDecoder(rr.Body).Decode(&responseData)
		assert.NoError(t, err, "error decoding response body")

		assert.NotEqual(t, 0, responseData.Status, "expected non-zero status code")
		assert.NotEmpty(t, responseData.Headers, "expected non-empty headers")
		assert.NotEqual(t, 0, responseData.Length, "expected non-zero length")
	})

	// Test case: Invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		invalidJSON := `{"method": "GET", "url": "http://google.com", "headers": {}`
		req := httptest.NewRequest(http.MethodPost, "/proxy", bytes.NewBufferString(invalidJSON))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(proxyHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Equal(t, "Invalid JSON format\n", rr.Body.String(), "handler returned unexpected body")
	})

	// Test case: Method not allowed
	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/proxy", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(proxyHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "handler returned wrong status code")
		assert.Equal(t, "Only POST method is allowed\n", rr.Body.String(), "handler returned unexpected body")
	})

	// Test case: Can't read request body
	t.Run("can't read request body", func(t *testing.T) {
		// Create a request with a nil body
		req := httptest.NewRequest(http.MethodPost, "/proxy", nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(proxyHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "handler returned wrong status code")
		assert.Equal(t, "Invalid JSON format\n", rr.Body.String(), "handler returned unexpected body")
	})
}

// Test main function
func TestMainFunction(t *testing.T) {
	// Set up a chi router with the same routes and middleware as main
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/proxy", proxyHandler)

	// Test case: Swagger route
	t.Run("swagger route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	// Test case: Proxy route
	t.Run("proxy route", func(t *testing.T) {
		requestData := RequestData{
			Method: "GET",
			URL:    "http://google.com",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		requestBody, _ := json.Marshal(requestData)
		req := httptest.NewRequest(http.MethodPost, "/proxy", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	// Test case: Middleware application (Logger and Recoverer)
	t.Run("middleware application", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		// The Logger middleware prints to the console, so we can't test its output directly here.
		// However, if the request was handled without panicking, the Recoverer middleware worked.
		assert.Equal(t, http.StatusOK, rr.Code, "middleware did not recover from a request")
	})
}
