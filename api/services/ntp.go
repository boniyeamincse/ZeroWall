package services

import (
	"fmt"
	"os"
)

// NTPServer represents Chrony or OpenNTPD settings
type NTPServer struct {
	Pools []string `json:"pools"`
}

// WriteNTPConfig applies settings to ntpd.conf
func (n NTPServer) WriteNTPConfig() error {
	content := "# ZeroWall NTP Config\n"
	for _, pool := range n.Pools {
		content += fmt.Sprintf("server %s iburst\n", pool)
	}
	return os.WriteFile("/etc/ntp.conf", []byte(content), 0644)
}
