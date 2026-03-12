package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	port       = flag.String("port", "8080", "HTTP port")
	tlsPort    = flag.String("tls-port", "443", "HTTPS port")
	certFile   = flag.String("cert", "/etc/zerowall/certs/server.crt", "TLS certificate")
	keyFile    = flag.String("key", "/etc/zerowall/certs/server.key", "TLS private key")
	enableTLS  = flag.Bool("tls", true, "Enable TLS")
	readTimeout  = flag.Int("read-timeout", 15, "Read timeout in seconds")
	writeTimeout = flag.Int("write-timeout", 15, "Write timeout in seconds")
	idleTimeout  = flag.Int("idle-timeout", 60, "Idle timeout in seconds")
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func redirectToHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			host := r.Host
			if strings.Contains(host, ":") {
				host = strings.Split(host, ":")[0]
			}
			httpsURL := "https://" + host + ":" + *tlsPort + r.URL.Path
			if r.URL.RawQuery != "" {
				httpsURL += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "online", "version": "1.0.0-Beta"}`)
	})

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"healthy": true, "timestamp": %d}`, time.Now().Unix())
	})

	mux.HandleFunc("/api/v1/login", handleLogin)

	mux.HandleFunc("/api/v1/firewall/rules", handleFirewallRules)

	mux.Handle("/ws/", handleWebSocket)

	return mux
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "jwt_token_placeholder"}`)
}

func handleFirewallRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"rules": []}`)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"error": "WebSocket not implemented"}`)
}

func newServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(*readTimeout) * time.Second,
		WriteTimeout: time.Duration(*writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(*idleTimeout) * time.Second,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
		},
	}
}

func main() {
	flag.Parse()

	baseHandler := loggingMiddleware(securityHeadersMiddleware(corsMiddleware(createMux())))

	if *enableTLS {
		go func() {
			tlsHandler := newServer(":"+*tlsPort, baseHandler)
			fmt.Printf("ZeroWall API [zwapi] HTTPS server starting on port %s...\n", *tlsPort)
			if err := tlsHandler.ListenAndServeTLS(*certFile, *keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server failed: %v", err)
			}
		}()
	}

	httpHandler := redirectToHTTPS(newServer(":"+*port, baseHandler))
	fmt.Printf("ZeroWall API [zwapi] HTTP server starting on port %s (redirects to HTTPS)...\n", *port)
	if err := httpHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
