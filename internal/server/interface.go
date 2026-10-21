package server

import "github.com/labstack/echo/v5"

type Router interface {
	// Use adds middleware to the chain.
	Use(middleware ...echo.MiddlewareFunc)
	// CONNECT register a new CONNECT route for a path.
	CONNECT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// DELETE register a new DELETE route for a path.
	DELETE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// GET register a new GET route for a path.
	GET(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// HEAD register a new HEAD route for a path.
	HEAD(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// OPTIONS register a new OPTIONS route for a path.
	OPTIONS(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// PATCH register a new PATCH route for a path.
	PATCH(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// POST register a new POST route for a path.
	POST(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// PUT register a new PUT route for a path.
	PUT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// QUERY register a new QUERY route for a path.
	QUERY(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// TRACE register a new TRACE route for a path.
	TRACE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// Any register a new route for a path.
	Any(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// Add registers a new route for an HTTP method and path.
	Add(method, path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) echo.RouteInfo
	// Group creates a new router group with prefix and optional group-level middleware.
	Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
}
