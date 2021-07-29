// Package common provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/algorand/oapi-codegen DO NOT EDIT.
package common

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns 200 if healthy.
	// (GET /health)
	MakeHealthCheck(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// MakeHealthCheck converts echo context to params.
func (w *ServerInterfaceWrapper) MakeHealthCheck(ctx echo.Context) error {

	validQueryParams := map[string]bool{
		"pretty": true,
	}

	// Check for unknown query parameters.
	for name, _ := range ctx.QueryParams() {
		if _, ok := validQueryParams[name]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unknown parameter detected: %s", name))
		}
	}

	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.MakeHealthCheck(ctx)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface, m ...echo.MiddlewareFunc) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/health", wrapper.MakeHealthCheck, m...)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
