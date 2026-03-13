package services

import (
	"fmt"
	"os/exec"
)

// DDNSClient handles Dynamic DNS updates
type DDNSClient struct {
	Provider string `json:"provider"`
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateAddress triggers an update to the provider
func (c DDNSClient) UpdateAddress(ip string) error {
	fmt.Printf("Updating %s for %s with IP %s via %s\n", c.Provider, c.Hostname, ip, c.Username)
	
	// Implementation calls a generic update URL using curl
	// Example: https://dyndns.example.com/update?hostname=myhost&myip=1.2.3.4
	updateURL := fmt.Sprintf("https://%s/update?hostname=%s&myip=%s", c.Provider, c.Hostname, ip)
	
	cmd := exec.Command("curl", "-u", fmt.Sprintf("%s:%s", c.Username, c.Password), updateURL)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("DDNS update failed: %v, output: %s", err, string(out))
	}
	
	return nil
}
