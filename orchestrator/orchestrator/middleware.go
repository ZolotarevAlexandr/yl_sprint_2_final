package orchestrator

import (
	"log"
	"net/http"
	"time"
)

// ResponseWriterWrapper is a wrapper around http.ResponseWriter to capture the status code.
type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and delegates to the original WriteHeader.
func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs the request details and response status code.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Request: %s %s", r.Method, r.URL.Path)

		rw := &ResponseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		log.Printf("Completed with status %d in %v", rw.statusCode, time.Since(start))
	})
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and send a user-friendly message
				log.Printf("Error occurred: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
