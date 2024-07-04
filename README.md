# Go Proxy Server

This is a simple HTTP proxy server implemented in Go. It allows users to proxy HTTP requests to external services and return the responses. The server also provides a Swagger UI for API documentation.

## Features

- Proxy HTTP requests and return the responses.
- Swagger UI for API documentation.
- Docker support for containerization.
- Unit tests for validating the proxy handler.

## Prerequisites

- Go 1.22 or later
- Docker (for containerization)

## Getting Started

### Clone the Repository

```bash
https://github.com/qonstant/go-proxy.git
cd go-proxy
```
## Build and Run Locally

### Build the application:

```bash
make build
```

### Run the application:

```bash
make run
```
After running make run, the server will start and be accessible at http://localhost:8080.

Link to the deployment: https://go-proxy-1fo6.onrender.com/proxy

### Generate Swagger Documentation

```bash
make swagger
```
### Run Tests
```bash
make test
```
## Docker
### Build the Docker Image
```bash
make docker-build
```
### Run the Docker Container
```bash
make docker-run
```
### Stop the Docker Container
```bash
make docker-stop
```
### Remove the Docker Container
```bash
make docker-remove
```

## API Endpoints

### Proxy Request

- URL: /proxy
- Method: POST
- Request Body:
```json
{
    "method": "GET",
    "url": "http://example.com",
    "headers": {
        "Content-Type": "application/json"
    }
}
```

For instance:
```json
{
    "method": "GET",
    "url": "http://jsonplaceholder.typicode.com/posts/1",
    "headers": {
        "Content-Type": "application/json"
    }
}
```
- Response:
```json
{
    "id": "1627563890765102000",
    "status": 200,
    "headers": {
        "Content-Type": "application/json"
    },
    "length": 1270
}
```

From the example above:
```json
{
    "id": "d7f4fb3d-6def-4378-857e-91711eb018c6",
    "status": 200,
    "headers": {
        "Access-Control-Allow-Credentials": "true",
        "Alt-Svc": "h3=\":443\"; ma=86400",
        "Cache-Control": "max-age=43200",
        "Cf-Cache-Status": "REVALIDATED",
        "Cf-Ray": "89db0c3a8d6d5d50-FRA",
        "Connection": "keep-alive",
        "Content-Type": "application/json; charset=utf-8",
        "Date": "Thu, 04 Jul 2024 00:37:37 GMT",
        "Etag": "W/\"124-yiKdLzqO5gfBrJFrcdJ8Yq0LGnU\"",
        "Expires": "-1",
        "Nel": "{\"report_to\":\"heroku-nel\",\"max_age\":3600,\"success_fraction\":0.005,\"failure_fraction\":0.05,\"response_headers\":[\"Via\"]}",
        "Pragma": "no-cache",
        "Report-To": "{\"group\":\"heroku-nel\",\"max_age\":3600,\"endpoints\":[{\"url\":\"https://nel.heroku.com/reports?ts=1719290587&sid=e11707d5-02a7-43ef-b45e-2cf4d2036f7d&s=eAlPGj2psKwqFTi3aRIeAycEDJsdhwHLI%2F0cXgblPNM%3D\"}]}",
        "Reporting-Endpoints": "heroku-nel=https://nel.heroku.com/reports?ts=1719290587&sid=e11707d5-02a7-43ef-b45e-2cf4d2036f7d&s=eAlPGj2psKwqFTi3aRIeAycEDJsdhwHLI%2F0cXgblPNM%3D",
        "Server": "cloudflare",
        "Vary": "Origin, Accept-Encoding",
        "Via": "1.1 vegur",
        "X-Content-Type-Options": "nosniff",
        "X-Powered-By": "Express",
        "X-Ratelimit-Limit": "1000",
        "X-Ratelimit-Remaining": "999",
        "X-Ratelimit-Reset": "1719290646"
    },
    "length": 292,
    "body": "{\n  \"userId\": 1,\n  \"id\": 1,\n  \"title\": \"sunt aut facere repellat provident occaecati excepturi optio reprehenderit\",\n  \"body\": \"quia et suscipit\\nsuscipit recusandae consequuntur expedita et cum\\nreprehenderit molestiae ut ut quas totam\\nnostrum rerum est autem sunt rem eveniet architecto\"\n}"
}
```

### Swagger Documentation

- URL: /swagger/

- Link: https://go-proxy-1fo6.onrender.com/swagger/index.html

## Project Structure

- main.go: The main server implementation.
- main_test.go: Unit tests for the proxy handler.
- Makefile: Makefile for building, running, testing, and Docker tasks.
- Dockerfile: Dockerfile for containerizing the application.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements.
