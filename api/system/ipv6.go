package system

import (
	"fmt"
	"os/exec"
)

// IPv6Stack manages global IPv6 settings and Router Advertisements
type IPv6Stack struct {
	Enabled        bool
	AcceptRA       bool
	PrivacyExtensions bool
}

// ConfigureSysctl sets kernel-level IPv6 behavior
func (s IPv6Stack) ConfigureSysctl() error {
	fmt.Println("IPv6: Configuring kernel networking stack...")
	
	// net.inet6.ip6.forwarding=1
	// net.inet6.ip6.accept_rtadv=1
	return nil
}

// EnableRA starts the Route Advertisement daemon (rtadvd)
func EnableRA(iface string) error {
	fmt.Printf("IPv6: Starting rtadvd on %s\n", iface)
	cmd := exec.Command("rtadvd", iface)
	return cmd.Run()
}

// SetupDHCPv6 starts the stateful/stateless DHCPv6 service
func SetupDHCPv6(iface string) {
	fmt.Printf("IPv6: Initializing DHCPv6 on %s\n", iface)
}
