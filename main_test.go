package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestProxyHandler(t *testing.T) {
    // Test case: Valid request
    t.Run("valid request", func(t *testing.T) {
        requestData := RequestData{
            Method: "GET",
            URL:    "http://example.com",
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

        if status := rr.Code; status != http.StatusOK {
            t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
        }

        var responseData ResponseData
        err := json.NewDecoder(rr.Body).Decode(&responseData)
        if err != nil {
            t.Errorf("error decoding response body: %v", err)
        }

        if responseData.Status == 0 {
            t.Errorf("expected non-zero status code")
        }

        if len(responseData.Headers) == 0 {
            t.Errorf("expected non-empty headers")
        }

        if responseData.Length == 0 {
            t.Errorf("expected non-zero length")
        }
    })

    // Test case: Invalid JSON
    t.Run("invalid JSON", func(t *testing.T) {
        invalidJSON := `{"method": "GET", "url": "http://example.com", "headers": {}`
        req := httptest.NewRequest(http.MethodPost, "/proxy", bytes.NewBufferString(invalidJSON))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(proxyHandler)
        handler.ServeHTTP(rr, req)

        if status := rr.Code; status != http.StatusBadRequest {
            t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
        }

        expected := "Invalid JSON format\n"
        if body := rr.Body.String(); body != expected {
            t.Errorf("handler returned unexpected body:\ngot:\n%v\nwant:\n%v\n", body, expected)
        }
    })

    // Test case: Method not allowed
    t.Run("method not allowed", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/proxy", nil)

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(proxyHandler)
        handler.ServeHTTP(rr, req)

        if status := rr.Code; status != http.StatusMethodNotAllowed {
            t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
        }

        expected := "Only POST method is allowed\n"
        if body := rr.Body.String(); body != expected {
            t.Errorf("handler returned unexpected body:\ngot:\n%v\nwant:\n%v\n", body, expected)
        }
    })
}
