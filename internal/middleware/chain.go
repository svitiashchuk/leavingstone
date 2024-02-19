package middleware

import "net/http"

// Middleware defines the structure of a middleware function
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Chain represents a chain of middlewares
type Chain struct {
	middlewares []Middleware
}

// NewChain creates a new chain of middlewares
func NewChain() *Chain {
	return &Chain{}
}

// Use adds a new middleware to the chain
func (c *Chain) Use(middleware Middleware) *Chain {
	c.middlewares = append(c.middlewares, middleware)

	return c
}

// Then applies the chain of middlewares to the given handler
func (c *Chain) Then(handler http.HandlerFunc) http.HandlerFunc {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}
