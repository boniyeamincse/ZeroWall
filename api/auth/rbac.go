package auth

import "net/http"

// RBACMiddleware checks permissions for a given route and role
func RBACMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Get token from Header
		// 2. Decode claims
		// 3. Check role
		userRole := "admin" // Mocked
		
		if userRole != requiredRole && userRole != "admin" {
			http.Error(w, "Forbidden: Insufficient Permissions", http.StatusForbidden)
			return
		}
		
		next(w, r)
	}
}
