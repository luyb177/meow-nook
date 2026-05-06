package middleware

import (
	"net/http"

	httpmd "github.com/luyb177/meow-nook/common/middleware/http"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return httpmd.AdminMiddleware(next)
}
