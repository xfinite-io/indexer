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
