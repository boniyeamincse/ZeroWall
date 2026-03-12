# ZeroWall — Performance & Security Improvements

This document outlines recommended improvements for performance optimization and security hardening based on code analysis.

---

## 1. Security Improvements

### Critical

#### 1.1 API Server Lacks TLS Encryption
**Location:** `api/main.go:27`

**Issue:** API server uses plain HTTP without TLS.

```go
// Current (insecure)
http.ListenAndServe(":"+port, mux)
```

**Recommendation:**
- Enable TLS with proper certificates
- Force HTTPS only via redirect middleware
- Add `http.Server` with configured timeouts

```go
// Recommended
srv := &http.Server{
    Addr:         ":" + port,
    Handler:      mux,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}
srv.ListenAndServeTLS("/etc/zerowall/certs/server.crt", "/etc/zerowall/certs/server.key")
```

---

#### 1.2 Mock JWT Implementation
**Location:** `api/auth/auth.go:15-24`

**Issue:** Authentication uses hardcoded mock tokens instead of proper JWT.

```go
// Current (insecure)
func GenerateToken(username string, role string, secret string) (string, error) {
    return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.mock_token.signature", nil
}

func VerifyToken(tokenString string, secret string) (*Claims, error) {
    return &Claims{Username: "admin", Role: "admin"}, nil
}
```

**Recommendation:**
- Implement proper JWT signing using HMAC or RSA
- Add proper token expiration
- Use library: `github.com/golang-jwt/jwt/v5`

---

#### 1.3 No Rate Limiting
**Issue:** API endpoints lack rate limiting, vulnerable to brute force and DDoS.

**Recommendation:**
- Add middleware for request throttling per IP/user
- Use token bucket or sliding window algorithm

```go
func rateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Implement rate limiting logic
    })
}
```

---

### High Priority

#### 1.4 Permissive File Permissions
**Location:** `api/firewall/pf_engine.go:73`

**Issue:** Config files written with `0644` (readable by all).

```go
os.WriteFile(tempPath, []byte(content), 0644)
```

**Recommendation:**
```go
os.WriteFile(tempPath, []byte(content), 0600)  // Owner read/write only
```

---

#### 1.5 No Input Validation
**Location:** `api/firewall/pf_engine.go:51-64`

**Issue:** Firewall rules accept raw strings without sanitization - risk of rule injection.

**Recommendation:**
- Validate all input against strict schemas
- Sanitize rule parameters (interface names, ports, IP addresses)
- Use regex whitelist for allowed characters

---

#### 1.6 Missing Security Headers
**Issue:** Web server doesn't set security headers.

**Recommendation:** Add headers:
- `Strict-Transport-Security` (HSTS)
- `Content-Security-Policy` (CSP)
- `X-Frame-Options: DENY`
- `X-Content-Type-Options: nosniff`

---

## 2. Performance Improvements

### High Priority

#### 2.1 Missing HTTP Server Timeouts
**Location:** `api/main.go:27`

**Issue:** No read/write timeouts configured - potential for slowloris attacks.

**Recommendation:**
```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  30 * time.Second,
}
```

---

#### 2.2 No Response Compression
**Issue:** API responses not compressed, wasting bandwidth.

**Recommendation:** Add gzip middleware:

```go
var gz http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
    if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()
        // Wrap response writer
    }
}
```

Or use `github.com/klauspost/compress/gzhttp`.

---

#### 2.3 Sysctl Tuning Enhancements
**Location:** `kernel/sysctl.conf`

**Current:** Basic TCP tuning present.

**Recommendations - Add:**

```sh
# TCP timestamps for RTT measurement
net.inet.tcp.timestamps=1

# SYN cache tuning for high connection rates
net.inet.tcp.syncache.hashsize=4096
net.inet.tcp.syncache.bucketlimit=100

# pf state limits
net.pf.request_maxcount=100000
net.pf.states_hashsize=65536

# TCP buffer auto-tuning
net.inet.tcp.recvbuf_auto=1
net.inet.tcp.sendbuf_auto=1

# Connection tracking
net.netflow.netflow_enable=1
```

---

#### 2.4 pf Optimization
**Issue:** No optimization settings for connection timeouts.

**Recommendation in generated pf.conf:**

```
# Connection tracking timeouts
set optimization aggressive

# State limits
set limit states 100000
set limit frags 50000

# Logging
set loginterface egress
set block-policy return
```

---

### Medium Priority

#### 2.5 No Caching Layer
**Issue:** Frequently accessed data (aliases, rules) recomputed on every request.

**Recommendation:**
- Add in-memory cache for static lookups
- Use `github.com/hashicorp/golang-lru` or `groupcache`
- Cache compiled rule sets

---

#### 2.6 Sequential Rule Processing
**Location:** `api/firewall/pf_engine.go:51-64`

**Issue:** Rules generated in sequential loop.

**Recommendation:**
- Pre-compile rule templates
- Use strings.Builder more efficiently
- Consider parallel validation of rule syntax

---

#### 2.7 No Connection Pooling
**Issue:** External service connections not pooled.

**Recommendation:**
- Use persistent HTTP clients
- Configure `MaxIdleConns`, `MaxIdleConnsPerHost`

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
}
client := &http.Client{Transport: transport}
```

---

## 3. Summary of Changes

| Priority | Component | Issue | Fix |
|----------|-----------|-------|-----|
| Critical | api/main.go | No TLS | Enable HTTPS with certificates |
| Critical | api/auth/auth.go | Mock JWT | Implement proper JWT signing |
| Critical | api/main.go | No rate limiting | Add throttling middleware |
| High | api/firewall/pf_engine.go | 0644 permissions | Use 0600 |
| High | api/* | No input validation | Add strict validation |
| High | api/main.go | No HTTP timeouts | Add timeouts to Server |
| High | kernel/sysctl.conf | Limited tuning | Add TCP/state limits |
| Medium | api/* | No compression | Add gzip middleware |
| Medium | api/* | No caching | Add in-memory cache |
| Medium | api/* | No connection pooling | Use persistent connections |

---

## 4. Testing Recommendations

After implementing changes:

1. **Security:**
   - Run OWASP ZAP against API endpoints
   - Test JWT token expiration/revocation
   - Verify rate limiting triggers correctly

2. **Performance:**
   - Load test with `wrk` or `vegeta`
   - Measure rule compilation time
   - Profile memory usage under load

---

*Generated: 2026-03-12*
