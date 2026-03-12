package security

import (
	"fmt"
)

// ContentFilter handles URL/Category blocking
type ContentFilter struct {
	Blacklist []string
	Enabled   bool
}

// BlockDomain adds a domain to the blacklist
func (f *ContentFilter) BlockDomain(domain string) {
	f.Blacklist = append(f.Blacklist, domain)
	fmt.Printf("Content Filter: Added %s to blacklist\n", domain)
}

// SyncToDNS applies the blacklist to Unbound/DNS server
func (f ContentFilter) SyncToDNS() error {
	fmt.Println("Syncing content filter rules to Unbound DNS...")
	// Implementation: write local-zone overrides to unbound.conf
	return nil
}

// UpdateCategories pulls global blacklist updates
func UpdateCategories() {
	fmt.Println("Content Filter: Pulling category-based blacklist (Adult, Gambling, Malware)...")
}
