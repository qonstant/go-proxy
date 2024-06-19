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

### Swagger Documentation

- URL: /swagger/

## Project Structure

- main.go: The main server implementation.
- main_test.go: Unit tests for the proxy handler.
- Makefile: Makefile for building, running, testing, and Docker tasks.
- Dockerfile: Dockerfile for containerizing the application.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements.

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). See the [LICENSE](LICENSE) file for details.
