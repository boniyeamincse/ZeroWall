package services

import (
	"fmt"
	"os"
)

// DHCPServer represents an ISC DHCP server instance
type DHCPServer struct {
	Interface string   `json:"interface"`
	Subnet    string   `json:"subnet"`
	Netmask   string   `json:"netmask"`
	RangeStart string  `json:"range_start"`
	RangeEnd   string  `json:"range_end"`
	DNS        []string `json:"dns"`
	Gateway    string   `json:"gateway"`
}

// GenerateConfig creates dhcpd.conf content
func (s DHCPServer) GenerateConfig() string {
	dnsStr := ""
	if len(s.DNS) > 0 {
		dnsStr = "  option domain-name-servers " + s.Gateway + ";"
	}

	return fmt.Sprintf(`subnet %s netmask %s {
  range %s %s;
  option routers %s;
%s
  default-lease-time 600;
  max-lease-time 7200;
}
`, s.Subnet, s.Netmask, s.RangeStart, s.RangeEnd, s.Gateway, dnsStr)
}

// SaveConfig writes the configuration to disk
func SaveConfig(content string) error {
	return os.WriteFile("/usr/local/etc/dhcpd.conf", []byte(content), 0644)
}
