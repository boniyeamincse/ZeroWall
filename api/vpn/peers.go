package vpn

import (
	"fmt"
	"os/exec"
)

// Peer represents a WireGuard remote peer
type Peer struct {
	PublicKey  string   `json:"public_key"`
	AllowedIPs []string `json:"allowed_ips"`
	Endpoint   string   `json:"endpoint"`
	Keepalive  int      `json:"keepalive"`
}

// AddPeer associates a peer with a WireGuard interface
func AddPeer(iface string, peer Peer) error {
	allowedStr := ""
	for _, ip := range peer.AllowedIPs {
		allowedStr += ip + ","
	}

	cmd := exec.Command("wg", "set", iface, "peer", peer.PublicKey, "endpoint", peer.Endpoint, "allowed-ips", allowedStr)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add peer: %s", string(out))
	}

	return nil
}
