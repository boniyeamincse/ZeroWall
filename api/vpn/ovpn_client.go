package vpn

import (
	"fmt"
	"os"
)

// OpenVPNClient represents a connection to a remote OpenVPN server
type OpenVPNClient struct {
	Name            string `json:"name"`
	Protocol        string `json:"protocol"` // udp, tcp
	ServerAddress   string `json:"server_addr"`
	ServerPort      int    `json:"server_port"`
	InterfaceName   string `json:"interface"` // e.g., ovpnc1
	Description     string `json:"description"`
	AuthType        string `json:"auth_type"` // tls, user_pass
	CACertPath      string `json:"ca_cert"`
	ClientCertPath  string `json:"client_cert"`
	ClientKeyPath   string `json:"client_key"`
	TLSAuthKeyPath  string `json:"tls_auth_key"`
	Compression     bool   `json:"compression"`
	ProxyType       string `json:"proxy_type"` // none, socks, http
	ProxyAddr       string `json:"proxy_addr"`
	ProxyPort       int    `json:"proxy_port"`
}

// GenerateClientConfig creates the client side .conf content
func (c OpenVPNClient) GenerateClientConfig() string {
	compStr := "comp-lzo no"
	if c.Compression {
		compStr = "comp-lzo yes"
	}

	proxyStr := ""
	if c.ProxyType != "none" {
		proxyStr = fmt.Sprintf("%s-proxy %s %d", c.ProxyType, c.ProxyAddr, c.ProxyPort)
	}

	return fmt.Sprintf(`client
dev %s
dev-type tun
proto %s
remote %s %d
resolv-retry infinite
nobind
persist-key
persist-tun
ca %s
cert %s
key %s
tls-auth %s 1
remote-cert-tls server
cipher AES-256-GCM
auth SHA256
verb 3
%s
%s
`, c.InterfaceName, c.Protocol, c.ServerAddress, c.ServerPort, 
   c.CACertPath, c.ClientCertPath, c.ClientKeyPath, c.TLSAuthKeyPath,
   compStr, proxyStr)
}

// SaveConfig writes the client configuration to disk
func (c OpenVPNClient) SaveConfig(content string) error {
	configDir := "/usr/local/etc/openvpn"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	path := fmt.Sprintf("%s/client-%s.conf", configDir, c.Name)
	return os.WriteFile(path, []byte(content), 0600)
}

// ConnectionStatus would check the actual status of the ovpncX interface
func (c OpenVPNClient) ConnectionStatus() (string, error) {
	// Logic to parse logs or ifconfig status
	return "unknown", nil
}
