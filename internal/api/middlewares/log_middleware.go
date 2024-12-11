package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/safepass/server/internal/logging"
)

type LogMiddleware struct {
	logger *logging.Logger
}

type ResponseCaptureWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseCaptureWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func NewLogMiddleware(logger *logging.Logger) *LogMiddleware {
	return &LogMiddleware{
		logger: logger,
	}
}

func (m *LogMiddleware) LogMiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &ResponseCaptureWriter{
			ResponseWriter: w,
		}

		m.logger.Info(fmt.Sprintf("%s%s %s %s%s", logging.BLACK, r.Method, r.Proto, r.URL.Path, logging.RESET))

		next.ServeHTTP(wrappedWriter, r)

		m.logger.Info(fmt.Sprintf("%s%d %s %v%s", logging.BLACK, wrappedWriter.statusCode, r.URL.Path, time.Since(start), logging.RESET))
	})
}
