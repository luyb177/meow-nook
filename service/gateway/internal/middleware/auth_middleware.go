package middleware

import (
	"net/http"

	httpmd "github.com/luyb177/meow-nook/common/middleware/http"
)

type AuthMiddleware struct {
	Secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		Secret: secret,
	}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return httpmd.AuthMiddleware(m.Secret)(next)
}
