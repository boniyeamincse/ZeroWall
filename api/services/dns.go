package services

import (
	"fmt"
	"os"
)

// DNSResolver represents Unbound configuration
type DNSResolver struct {
	ConfigPath   string   `json:"config_path"`
	Interfaces   []string `json:"interfaces"`
	Port         int      `json:"port"`
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
func (d DNSResolver) UpdateResolver(config string) error {
	path := d.ConfigPath
	if path == "" {
		path = "/var/unbound/unbound.conf"
	}
	return os.WriteFile(path, []byte(config), 0644)
}
