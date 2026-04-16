// Package middleware contains gateway-level HTTP middleware.
package middleware

import (
	"net/http"

	"github.com/luyb177/meow-nook/common/middleware"
)

// AuthMiddleware wraps the shared JWT middleware for use inside the go-zero
// service context.  The secret is injected at start-up via NewAuthMiddleware.
type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: secret}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return middleware.JWTMiddleware(m.secret)(next).ServeHTTP
}
