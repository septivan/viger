package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"runtime/debug"
	"slices"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request-id"

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (writer *statusWriter) WriteHeader(status int) {
	writer.status = status
	writer.ResponseWriter.WriteHeader(status)
}

func (writer *statusWriter) Unwrap() http.ResponseWriter { return writer.ResponseWriter }

func (server *api) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		writer.Header().Set("Referrer-Policy", "no-referrer")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "DENY")
		writer.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(writer, request)
	})
}

func (server *api) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if origin != "" && slices.Contains(server.settings.AllowedOrigins, origin) {
			writer.Header().Set("Access-Control-Allow-Origin", origin)
			writer.Header().Set("Vary", "Origin")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
		}
		if request.Method == http.MethodOptions {
			if origin == "" || !slices.Contains(server.settings.AllowedOrigins, origin) {
				writeError(writer, request, http.StatusForbidden, "ORIGIN_NOT_ALLOWED", "The request origin is not allowed.", nil)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (server *api) requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestID := request.Header.Get("X-Request-ID")
		if requestID == "" || len(requestID) > 128 {
			requestID = newRequestID()
		}
		writer.Header().Set("X-Request-ID", requestID)
		started := time.Now()
		tracked := &statusWriter{ResponseWriter: writer, status: http.StatusOK}
		server.settings.Metrics.RecordHTTPRequest()
		next.ServeHTTP(tracked, request.WithContext(context.WithValue(request.Context(), requestIDKey, requestID)))
		server.settings.Logger.Info("HTTP request", "request_id", requestID, "method", request.Method, "path", request.URL.Path, "status", tracked.status, "duration_ms", time.Since(started).Milliseconds())
	})
}

func (server *api) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				server.settings.Logger.Error("panic recovered", "request_id", requestID(request.Context()), "error", recovered, "stack", string(debug.Stack()))
				writeError(writer, request, http.StatusInternalServerError, "INTERNAL_ERROR", "The request could not be completed.", nil)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func requestID(context context.Context) string {
	if value, ok := context.Value(requestIDKey).(string); ok {
		return value
	}
	return "unavailable"
}

func newRequestID() string {
	value := make([]byte, 12)
	if _, err := rand.Read(value); err != nil {
		return "request-id-unavailable"
	}
	return hex.EncodeToString(value)
}
