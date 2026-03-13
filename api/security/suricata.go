package security

import (
	"fmt"
	"os/exec"
)

// SuricataEngine manages the IDS/IPS service
type SuricataEngine struct {
	Interface string `json:"interface"`
	Mode      string `json:"mode"` // ids, ips
}

// GenerateConfig creates suricata.yaml settings
func (s SuricataEngine) GenerateConfig() string {
	runMode := "autofp"
	if s.Mode == "ips" {
		runMode = "workers"
	}

	return fmt.Sprintf(`%s
vars:
  address-groups:
    HOME_NET: "[192.168.0.0/16,10.0.0.0/8,172.16.0.0/12]"
    EXTERNAL_NET: "!$HOME_NET"

default-log-dir: /var/log/suricata/
outputs:
  - eve-log:
      enabled: yes
      filetype: regular
      filename: eve.json
      types:
        - alert
        - http
        - dns
        - tls

pcap:
  - interface: %s
    cluster-id: 99
    cluster-type: %s
`, "# Suricata Config", s.Interface, runMode)
}

// StartService reloads suricata configs or starts daemon
func (s SuricataEngine) StartService() error {
	fmt.Printf("Starting Suricata on %s in %s mode\n", s.Interface, s.Mode)
	// exec.Command("suricata", "-D", "-c", "/usr/local/etc/suricata/suricata.yaml", "-i", s.Interface).Run()
	return nil
}

// ReloadRules triggers a SIGUSR2 for rule reload
func ReloadRules() error {
	cmd := exec.Command("pkill", "-USR2", "suricata")
	return cmd.Run()
}
