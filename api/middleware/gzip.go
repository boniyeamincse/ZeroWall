package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw *GzipResponseWriter) Write(p []byte) (int, error) {
	return grw.Writer.Write(p)
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

		w.Header().Set("Content-Encoding", " gzip")

		grw := &GzipResponseWriter{
			Writer:         gz,
			ResponseWriter: w,
		}

		gw := &gzipWriter{
			ResponseWriter: w,
			Writer:         gz,
		}

		next.ServeHTTP(gw, r)
		gz.Close()
	})
}

type gzipWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (gw *gzipWriter) Write(p []byte) (int, error) {
	return gw.Writer.Write(p)
}

func Gzip() func(http.Handler) http.Handler {
	handler := NewGzipHandler()
	return handler.Middleware
}
