package vpn

import (
	"fmt"
	"os"
)

// OpenVPNServer represents an OpenVPN SSL/TLS server instance
type OpenVPNServer struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"` // udp, tcp
	Port          int    `json:"port"`
	DeviceMode    string `json:"dev_mode"` // tun, tap
	InterfaceName string `json:"interface"`
	TunnelNetwork string `json:"tunnel_network"` // e.g., 10.0.8.0/24
	LocalNetwork  string `json:"local_network"`  // e.g., 192.168.1.0/24
	Ciphers       string `json:"ciphers"`
	AuthDigest    string `json:"auth_digest"`
	EnableLZO     bool   `json:"enable_lzo"`
	CACertPath    string `json:"ca_cert"`
	ServerCert    string `json:"server_cert"`
	ServerKey     string `json:"server_key"`
	DHParamPath   string `json:"dh_param"`
}

// GenerateConfig creates the openvpn.conf content
func (s OpenVPNServer) GenerateConfig() string {
	lzoStr := ""
	if s.EnableLZO {
		lzoStr = "comp-lzo yes"
	}

	return fmt.Sprintf(`dev %s
dev-type %s
proto %s
port %d
ca %s
cert %s
key %s
dh %s
topology subnet
server %s
push "route %s"
cipher %s
auth %s
user nobody
group nobody
persist-key
persist-tun
status /var/log/openvpn-status-%s.log
verb 3
%s
`, s.InterfaceName, s.DeviceMode, s.Protocol, s.Port, s.CACertPath, s.ServerCert, s.ServerKey, s.DHParamPath, 
   s.TunnelNetwork, s.LocalNetwork, s.Ciphers, s.AuthDigest, s.Name, lzoStr)
}

// SaveConfig writes the OpenVPN configuration to a specific path
func (s OpenVPNServer) SaveConfig(configContent string) error {
	configDir := "/usr/local/etc/openvpn"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create openvpn config dir: %v", err)
	}
	
	path := fmt.Sprintf("%s/%s.conf", configDir, s.Name)
	return os.WriteFile(path, []byte(configContent), 0600)
}

// StartService would normally interact with FreeBSD rc.d
func (s OpenVPNServer) StartService() error {
	fmt.Printf("Starting OpenVPN instance: %s\n", s.Name)
	// exec.Command("service", "openvpn", "start", s.Name).Run()
	return nil
}
