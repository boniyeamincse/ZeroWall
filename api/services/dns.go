package services

import (
	"fmt"
	"os"
)

// DNSResolver represents Unbound configuration
type DNSResolver struct {
	Interfaces []string `json:"interfaces"`
	Port        int      `json:"port"`
	EnableDNSSEC bool     `json:"enable_dnssec"`
}

// CreateDNSConfig generates unbound.conf
func (d DNSResolver) CreateDNSConfig() string {
	ifaceStr := ""
	for _, iface := range d.Interfaces {
		ifaceStr += fmt.Sprintf("  interface: %s\n", iface)
	}

	dnssecStr := "no"
	if d.EnableDNSSEC {
		dnssecStr = "yes"
	}

	return fmt.Sprintf(`server:
  port: %d
  do-ip4: yes
  do-ip6: yes
  do-udp: yes
  do-tcp: yes
%s
  validator-enable: %s
  access-control: 0.0.0.0/0 allow
`, d.Port, ifaceStr, dnssecStr)
}

// UpdateResolver applies the configuration
func UpdateResolver(config string) error {
	return os.WriteFile("/var/unbound/unbound.conf", []byte(config), 0644)
}
