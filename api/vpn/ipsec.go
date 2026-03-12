package vpn

import (
	"fmt"
	"os"
)

// IPsecConfig represents a StrongSwan IKEv2 connection
type IPsecConfig struct {
	Name            string `json:"name"`
	RemoteGateway   string `json:"remote_gw"`
	LocalID         string `json:"local_id"`
	RemoteID        string `json:"remote_id"`
	PSK             string `json:"psk"`
	LocalSubnets    []string `json:"local_ts"`
	RemoteSubnets   []string `json:"remote_ts"`
	EncryptionProto string `json:"enc_proto"` // aes256-sha256-modp2048
}

// GenerateSwanConfig creates swanctl.conf content
func (i IPsecConfig) GenerateSwanConfig() string {
	tsLocal := ""
	for _, s := range i.LocalSubnets { tsLocal += s + "," }
	
	tsRemote := ""
	for _, s := range i.RemoteSubnets { tsRemote += s + "," }

	return fmt.Sprintf(`connections {
    %s {
        remote_addrs = %s
        local {
            auth = psk
            id = %s
        }
        remote {
            auth = psk
            id = %s
        }
        children {
            %s-child {
                local_ts = %s
                remote_ts = %s
                esp_proposals = %s
            }
        }
    }
}
secrets {
    ike-%s {
        id = %s
        secret = "%s"
    }
}
`, i.Name, i.RemoteGateway, i.LocalID, i.RemoteID, i.Name, tsLocal, tsRemote, i.EncryptionProto, i.Name, i.RemoteID, i.PSK)
}

// SaveConfig writes the strongSwan configuration
func (i IPsecConfig) SaveConfig(config string) error {
	path := fmt.Sprintf("/usr/local/etc/swanctl/conf.d/%s.conf", i.Name)
	return os.WriteFile(path, []byte(config), 0600)
}
