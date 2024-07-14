package main

import (
	"go-proxy/internal/app"
)

// go-proxy-1fo6.onrender.com
// localhost:8080

// @title HTTP Proxy Server API
// @version 1.0
// @description This is a simple HTTP proxy server.
// @BasePath /
// @host go-proxy-1fo6.onrender.com
// @schemes http https
func main() {
	app.Run()
}
