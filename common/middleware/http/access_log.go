package httpmw

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/luyb177/meow-nook/common/logger"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Write(b []byte) (int, error) {
	// 如果业务没显式 WriteHeader，默认 200
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func AccessLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rid := strings.TrimSpace(r.Header.Get("X-Request-Id"))
		if rid == "" {
			rid = newRequestID()
		}
		w.Header().Set("X-Request-Id", rid)

		ctx := logger.WithRequestID(r.Context(), rid)

		if tid := logger.TraceIDFromOTel(ctx); tid != "" {
			ctx = logger.WithTraceID(ctx, tid)
		}

		r = r.WithContext(ctx)

		rec := &statusRecorder{ResponseWriter: w}
		next(rec, r)

		logger.FromContext(ctx).Info("http access",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", rec.status),
			zap.Int("bytes", rec.bytes),
			zap.Int64("latency_ms", time.Since(start).Milliseconds()),
			zap.String("remote", r.RemoteAddr),
		)
	}
}

func newRequestID() string {
	// 16 bytes => 32 hex chars
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
