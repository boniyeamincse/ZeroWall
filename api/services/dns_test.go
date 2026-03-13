package services

import (
	"strings"
	"testing"
)

func TestDNSResolver_CreateDNSConfig(t *testing.T) {
	resolver := DNSResolver{
		Port:         53,
		Interfaces:   []string{"127.0.0.1", "192.168.1.1"},
		EnableDNSSEC: true,
	}

	got := resolver.CreateDNSConfig()

	wants := []string{
		"port: 53",
		"interface: 127.0.0.1",
		"interface: 192.168.1.1",
		"validator-enable: yes",
		"access-control: 0.0.0.0/0 allow",
	}

	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("CreateDNSConfig() = %q, want to contain %q", got, want)
		}
	}
}
