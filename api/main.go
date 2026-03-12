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

	"zerowall/api/auth"
	"zerowall/api/firewall"
	"zerowall/api/middleware"
)

var (
	port        = flag.String("port", "8080", "HTTP port")
	tlsPort     = flag.String("tls-port", "443", "HTTPS port")
	certFile    = flag.String("cert", "/etc/zerowall/certs/server.crt", "TLS certificate")
	keyFile     = flag.String("key", "/etc/zerowall/certs/server.key", "TLS private key")
	enableTLS   = flag.Bool("tls", true, "Enable TLS")
	readTimeout  = flag.Int("read-timeout", 15, "Read timeout in seconds")
	writeTimeout = flag.Int("write-timeout", 15, "Write timeout in seconds")
	idleTimeout  = flag.Int("idle-timeout", 60, "Idle timeout in seconds")
	jwtSecret   = flag.String("jwt-secret", "zerowall-secret-key-change-in-production", "JWT secret")
)

var jwt *auth.JWT

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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/login" || r.URL.Path == "/api/v1/health" || r.URL.Path == "/api/v1/status" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "authorization required"}`, http.StatusUnauthorized)
			return
		}

		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			http.Error(w, `{"error": "invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		claims, err := jwt.VerifyToken(token)
		if err != nil {
			http.Error(w, `{"error": "invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User", claims.Username)
		r.Header.Set("X-Role", claims.Role)
		next.ServeHTTP(w, r)
	})
}

func rateLimiterMiddleware(next http.Handler) http.Handler {
	return middleware.RateLimit(100, 200)(next)
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

	fwHandler := firewall.NewFirewallHandler()
	statsHandler := firewall.NewStatsHandler()
	queueHandler := firewall.NewQueueHandler()

	mux.HandleFunc("/api/v1/status", handleStatus)
	mux.HandleFunc("/api/v1/health", handleHealth)
	mux.HandleFunc("/api/v1/login", handleLogin)
	mux.HandleFunc("/api/v1/logout", handleLogout)

	mux.HandleFunc("/api/v1/firewall/rules", fwHandler.GetRules)
	mux.HandleFunc("/api/v1/firewall/rules/", fwHandler.GetRule)
	mux.HandleFunc("/api/v1/firewall/rule", fwHandler.CreateRule)
	mux.HandleFunc("/api/v1/firewall/rule/", fwHandler.UpdateRule)
	mux.HandleFunc("/api/v1/firewall/rule/", fwHandler.DeleteRule)
	mux.HandleFunc("/api/v1/firewall/rules/reorder", fwHandler.ReorderRules)
	mux.HandleFunc("/api/v1/firewall/rule/toggle/", fwHandler.ToggleRule)

	mux.HandleFunc("/api/v1/firewall/nat", fwHandler.GetNATRules)
	mux.HandleFunc("/api/v1/firewall/nat/", fwHandler.CreateNATRule)
	mux.HandleFunc("/api/v1/firewall/nat/", fwHandler.DeleteNATRule)

	mux.HandleFunc("/api/v1/firewall/aliases", fwHandler.GetAliases)

	mux.HandleFunc("/api/v1/firewall/stats", statsHandler.GetStats)
	mux.HandleFunc("/api/v1/firewall/states", statsHandler.GetStateList)
	mux.HandleFunc("/api/v1/firewall/logs", statsHandler.GetLogs)
	mux.HandleFunc("/api/v1/firewall/flush", statsHandler.FlushStates)

	mux.HandleFunc("/api/v1/firewall/queues", queueHandler.GetQueues)
	mux.HandleFunc("/api/v1/firewall/queue", queueHandler.CreateQueue)
	mux.HandleFunc("/api/v1/firewall/queue/", queueHandler.UpdateQueue)
	mux.HandleFunc("/api/v1/firewall/queue/", queueHandler.DeleteQueue)
	mux.HandleFunc("/api/v1/firewall/queue/stats", queueHandler.GetQueueStats)

	mux.HandleFunc("/api/v1/firewall/apply", fwHandler.ApplyFirewall)

	mux.Handle("/ws/", handleWebSocket)

	return mux
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "online", "version": "1.0.0-Beta"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"healthy": true, "timestamp": %d}`, time.Now().Unix())
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := decodeJSON(r.Body, &creds); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, `{"error": "username and password required"}`, http.StatusBadRequest)
		return
	}

	token, err := jwt.GenerateToken(creds.Username, "admin", 24*time.Hour)
	if err != nil {
		http.Error(w, `{"error": "failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success": true, "token": "%s", "user": {"username": "%s", "role": "admin"}}`, token, creds.Username)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success": true, "message": "logged out"}`)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"error": "WebSocket not implemented"}`)
}

func decodeJSON(body interface{}, v interface{}) error {
	return nil
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

	jwt = auth.NewJWT(*jwtSecret)

	baseHandler := loggingMiddleware(
		securityHeadersMiddleware(
			corsMiddleware(
				authMiddleware(
					rateLimiterMiddleware(
						createMux())))))

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
