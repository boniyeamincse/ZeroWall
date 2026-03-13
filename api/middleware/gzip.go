package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

type GzipHandler struct {
	pool sync.Pool
}

func NewGzipHandler() *GzipHandler {
	return &GzipHandler{
		pool: sync.Pool{
			New: func() interface{} {
				gz, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
				return gz
			},
		},
	}
}

func (gh *GzipHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/ws/" || strings.HasPrefix(r.URL.Path, "/ws/") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gh.pool.Get().(*gzip.Writer)
		defer gh.pool.Put(gz)
		gz.Reset(w)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		grw := &gzipResponseWriter{
			Writer:         gz,
			ResponseWriter: w,
		}

		next.ServeHTTP(grw, r)
	})
}

func Gzip() func(http.Handler) http.Handler {
	handler := NewGzipHandler()
	return handler.Middleware
}
