# Variables
APP_NAME = go-proxy-server
DOCKER_IMAGE = $(APP_NAME)
DOCKER_CONTAINER = $(APP_NAME)
GO_FILES = $(wildcard *.go)
SWAGGER_FILES = $(wildcard docs/*)

# Build the Go application
build:
	@echo "Building the Go application..."
	@go build -o $(APP_NAME)

# Run the application locally
run: build
	@echo "Running the application..."
	@./$(APP_NAME)

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init

# Test the application
test:
	@echo "Running tests..."
	@go test -v

# Test the application
test-html:
	@echo "Creation of UI for tests..."
	@go test -coverprofile=cover.txt
	@go tool cover -html=cover.txt


# Build the Docker image
docker-build:
	@echo "Building the Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

# Run the Docker container
docker-run: docker-build
	@echo "Running the Docker container..."
	@docker run -d -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

# Stop the Docker container
docker-stop:
	@echo "Stopping the Docker container..."
	@docker stop $(DOCKER_CONTAINER)

# Remove the Docker container
docker-remove:
	@echo "Removing the Docker container..."
	@docker rm $(DOCKER_CONTAINER)

# Docker Compose up
up:
	@echo "Docker compose up..."
	@docker compose up -d

# Docker Compose down
down:
	@echo "Docker compose down..."
	@docker compose down
	@docker rmi go-proxy-server

# Restart container
restart: swagger down up

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
	# @rm -rf docs

# Full workflow: clean, build, swagger, test, docker-build, and docker-run
all: clean build swagger test docker-build docker-run

.PHONY: build run swagger test docker-build docker-run docker-stop docker-remove clean all