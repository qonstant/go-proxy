// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/proxy": {
            "post": {
                "description": "Proxy an HTTP request to an external service and return the response.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "proxy"
                ],
                "summary": "Proxy HTTP request",
                "parameters": [
                    {
                        "description": "Request Data",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.RequestData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseData"
                        }
                    },
                    "400": {
                        "description": "Invalid JSON format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "405": {
                        "description": "Only POST method is allowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.RequestData": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "headers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "method": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "main.ResponseData": {
            "type": "object",
            "properties": {
                "headers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "length": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "HTTP Proxy Server API",
	Description:      "This is a simple HTTP proxy server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}