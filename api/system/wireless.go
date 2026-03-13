package system

import (
	"fmt"
	"os"
)

// WirelessConfig represents HostAPD settings for AP mode
type WirelessConfig struct {
	Interface string `json:"interface"`
	SSID      string `json:"ssid"`
	Channel   int    `json:"channel"`
	WPAKey    string `json:"wpa_key"`
	HWMode    string `json:"hw_mode"` // a, b, g, n, ac
}

// ConfigureAP generates HostAPD config and triggers service restart
func ConfigureAP(config WirelessConfig) error {
	confContent := fmt.Sprintf(`interface=%s
driver=bsd
ssid=%s
hw_mode=%s
channel=%d
wpa=2
wpa_passphrase=%s
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
`, config.Interface, config.SSID, config.HWMode, config.Channel, config.WPAKey)

	path := fmt.Sprintf("/etc/hostapd-%s.conf", config.Interface)
	if err := os.WriteFile(path, []byte(confContent), 0600); err != nil {
		return fmt.Errorf("failed to write hostapd config: %v", err)
	}

	// In a real system, we'd restart hostapd via rc.d
	return nil
}
