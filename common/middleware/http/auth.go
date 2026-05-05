package httpmw

import (
	"context"
	"net/http"
	"strings"

	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/luyb177/meow-nook/common/httpresp"
	"github.com/luyb177/meow-nook/common/jwtx"
)

type claimsKey struct{}

func ClaimsFromContext(ctx context.Context) (*jwtx.Claims, bool) {
	c, ok := ctx.Value(claimsKey{}).(*jwtx.Claims)
	return c, ok
}

func AuthMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractBearer(r)
			if tokenStr == "" {
				httpresp.JsonBaseResponseCtx(r.Context(), w, errorx.ErrUnauthorized)
				return
			}

			claims, err := jwtx.Parse(tokenStr, secret)
			if err != nil {
				httpresp.JsonBaseResponseCtx(r.Context(), w, errorx.ErrTokenInvalid)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			next(w, r.WithContext(ctx))
		}
	}
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok || claims.Role != jwtx.RoleAdmin {
			httpresp.JsonBaseResponseCtx(r.Context(), w, errorx.ErrPermissionDenied)
			return
		}
		next(w, r)
	}
}

func extractBearer(r *http.Request) string {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}
