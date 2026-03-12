package services

import (
	"fmt"
)

// ProxyBackend represents a backend server for the reverse proxy
type ProxyBackend struct {
	Address string
	Port    int
	Weight  int
}

// ReverseProxyConfig handles SSL offloading and load balancing
type ReverseProxyConfig struct {
	ListenPort int
	SSLCert    string
	Backends   []ProxyBackend
}

// GenerateNGINXConfig creates an upstream and server block config
func (c ReverseProxyConfig) GenerateNGINXConfig() string {
	config := "upstream backend_farm {\n"
	for _, b := range c.Backends {
		config += fmt.Sprintf("  server %s:%d weight=%d;\n", b.Address, b.Port, b.Weight)
	}
	config += "}\n\nserver {\n"
	config += fmt.Sprintf("  listen %d ssl;\n", c.ListenPort)
	config += fmt.Sprintf("  ssl_certificate %s;\n", c.SSLCert)
	config += "  location / {\n    proxy_pass http://backend_farm;\n  }\n}\n"
	return config
}

// StartProxy reloads the proxy service
func StartProxy() error {
	fmt.Println("Reloading Reverse Proxy service...")
	return nil
}
