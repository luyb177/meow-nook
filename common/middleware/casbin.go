package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/luyb177/meow-nook/common/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// CasbinMiddleware uses a Casbin enforcer to check whether the caller is
// allowed to access the requested (path, method) resource.
//
// The subject is derived from the JWT claims stored in the request context by
// JWTMiddleware.  Requests that have no claims attached are rejected with 401.
func CasbinMiddleware(enforcer *casbin.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				httpx.Error(w, errorx.ErrUnauthorized)
				return
			}

			obj := r.URL.Path
			act := r.Method

			ok, err := enforcer.Enforce(claims.Role, obj, act)
			if err != nil {
				httpx.Error(w, errorx.Wrap(errorx.CodeInternalError, "权限校验失败", err))
				return
			}
			if !ok {
				httpx.Error(w, errorx.ErrPermissionDenied)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
