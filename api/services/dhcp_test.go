package services

import (
	"strings"
	"testing"
)

func TestDHCPServer_GenerateConfig(t *testing.T) {
	server := DHCPServer{
		Subnet:     "192.168.1.0",
		Netmask:    "255.255.255.0",
		RangeStart: "192.168.1.100",
		RangeEnd:   "192.168.1.200",
		Gateway:    "192.168.1.1",
		DNS:        []string{"8.8.8.8"},
	}

	got := server.GenerateConfig()

	wants := []string{
		"subnet 192.168.1.0 netmask 255.255.255.0",
		"range 192.168.1.100 192.168.1.200",
		"option routers 192.168.1.1",
		"option domain-name-servers 192.168.1.1",
	}

	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("GenerateConfig() = %q, want to contain %q", got, want)
		}
	}
}
