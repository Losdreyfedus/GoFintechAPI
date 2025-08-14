package middleware

import (
	"net/http"
	"time"

	"backend_path/pkg/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TracingMiddleware adds OpenTelemetry tracing to HTTP requests
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract trace context from headers
		ctx := tracing.ExtractTraceContext(r.Context(), headersToMap(r.Header))

		// Start span for this request
		ctx, span := tracing.StartSpanFromRequest(ctx, r.Method+" "+r.URL.Path)
		defer span.End()

		// Add request attributes to span
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.remote_addr", r.RemoteAddr),
			attribute.String("http.request_id", r.Header.Get("X-Request-ID")),
		)

		// Create response writer wrapper to capture status code
		wrapped := &tracingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Add trace context to request
		r = r.WithContext(ctx)

		// Process request
		start := time.Now()
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		// Add response attributes to span
		span.SetAttributes(
			attribute.Int("http.status_code", wrapped.statusCode),
			attribute.String("http.duration", duration.String()),
		)

		// Set span status based on HTTP status code
		if wrapped.statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}
	})
}

// tracingResponseWriter wraps http.ResponseWriter to capture status code
type tracingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *tracingResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *tracingResponseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

// headersToMap converts http.Header to map[string]string
func headersToMap(headers http.Header) map[string]string {
	result := make(map[string]string)
	for key, values := range headers {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result
}
