# Use the official Golang image as the base image
FROM golang:1.22 as builder

# Set the working directory within the builder container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the builder container
COPY . .

# Build the Go app for Linux (amd64) with CGO disabled
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-proxy-server

# Create a new stage for the final application image (based on Alpine Linux)
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/go-proxy-server ./go-proxy-server

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./go-proxy-server"]
