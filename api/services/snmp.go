package services

import (
	"fmt"
	"os"
)

// SNMPAgent represents bsnmpd settings
type SNMPAgent struct {
	Community string `json:"community"`
	Location  string `json:"location"`
	Contact   string `json:"contact"`
}

// ApplySNMPConfig generates snmpd.config
func (s SNMPAgent) ApplySNMPConfig() error {
	content := fmt.Sprintf(`read := "public"
location := "%s"
contact := "%s"
system := 1
`, s.Location, s.Contact)
	
	return os.WriteFile("/etc/snmpd.config", []byte(content), 0644)
}
