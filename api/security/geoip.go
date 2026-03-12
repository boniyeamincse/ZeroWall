package security

import (
	"fmt"
	"os"
	"os/exec"
)

// GeoIPManager handles MaxMind-based GeoIP blocking
type GeoIPManager struct {
	DBPath    string
	PFTable   string // pf table name, e.g., <blocked_countries>
}

// CreateGeoIPTable initializes a pf table from a CIDR list for a country
func (m GeoIPManager) CreateGeoIPTable(countryCode string, cidrListPath string) error {
	// 1. Validate CIDR list exists
	if _, err := os.Stat(cidrListPath); os.IsNotExist(err) {
		return fmt.Errorf("cidr list for %s not found: %s", countryCode, cidrListPath)
	}

	// 2. Load into pf table
	cmd := exec.Command("pfctl", "-t", m.PFTable, "-T", "replace", "-f", cidrListPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to populate geoip table: %s", string(out))
	}

	return nil
}

// BlockCountry adds a country to the active blocklist policy
func (m GeoIPManager) BlockCountry(countryCode string) error {
	fmt.Printf("Blocking all traffic from %s using table %s\n", countryCode, m.PFTable)
	// Implementation would append to a persistent config and reload pf
	return nil
}
