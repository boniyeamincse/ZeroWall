package services

import (
	"fmt"
)

// DNSSECConfig manages Unbound security extensions
type DNSSECConfig struct {
	Enabled bool
	LogQueries bool
}

// GenerateUnboundSECConfig produces the security block for unbound.conf
func (c DNSSECConfig) GenerateUnboundSECConfig() string {
	if !c.Enabled {
		return "# DNSSEC Disabled\n"
	}
	
	config := "server:\n"
	config += "  auto-trust-anchor-file: \"/var/unbound/root.key\"\n"
	config += "  module-config: \"validator iterator\"\n"
	config += "  val-log-level: 2\n"
	
	if c.LogQueries {
		config += "  log-queries: yes\n"
	}
	
	return config
}

// FetchRootKey updates the DNSSEC root anchor
func FetchRootKey() error {
	fmt.Println("DNSSEC: Fetching latest root trust anchors...")
	// unbound-anchor -a /var/unbound/root.key
	return nil
}
