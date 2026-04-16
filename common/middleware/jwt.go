// Package middleware provides reusable HTTP middleware for the meow-nook
// gateway, including JWT authentication and Casbin-based authorisation.
package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// ──────────────────────────────────────────────
// JWT
// ──────────────────────────────────────────────

// Claims is the custom JWT payload used by meow-nook.
type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates the Bearer token in the Authorization header and
// injects the parsed claims into the request context under claimsKey.
func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				httpx.Error(w, errorx.ErrUnauthorized)
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errorx.ErrTokenInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				httpx.Error(w, errorx.ErrTokenInvalid)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

func extractToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}
	return ""
}
