package main

import (
	"fmt"
	"net/http"
)

// InitSwagger sets up the /api/docs endpoint
func InitSwagger(mux *http.ServeMux) {
	fmt.Println("API Docs: Initializing Swagger/OpenAPI exporter...")
	
	mux.HandleFunc("/api/v1/docs.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"swagger": "2.0",
			"info": { "title": "ZeroWall API", "version": "1.0" },
			"paths": { "/api/v1/status": { "get": { "responses": { "200": { "description": "OK" } } } } }
		}`)
	})
}
