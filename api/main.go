package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// zwapi server entry point
func main() {
	port := os.Getenv("ZW_API_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Generic API routes
	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status": "online", "version": "1.0.0-Beta"}`)
	})

	// TODO: Add complex routes for firewall, vpn, system modules

	fmt.Printf("ZeroWall API [zwapi] starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("API failed to start: %v", err)
	}
}
