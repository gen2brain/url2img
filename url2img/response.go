package main

import (
	"net/http"
)

// responseWriter is http.ResponseWriter with size and status
type responseWriter struct {
	http.ResponseWriter

	size   int
	status int
}

// NewResponseWriter returns new responseWriter
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, 0, 0}
}

// Write writes the data to the connection as part of an HTTP reply
func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// WriteHeader sends an HTTP response header with status code
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Size returns response size
func (w *responseWriter) Size() int {
	return w.size
}

// Status returns response status code
func (w *responseWriter) Status() int {
	return w.status
}
