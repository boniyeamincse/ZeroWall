package system

import (
	"bytes"
	"fmt"
	"net/http"
)

// SyncManager handles XMLRPC configuration synchronization between HA nodes
type SyncManager struct {
	PeerIP   string
	Password string
}

// SyncConfig pushes local configuration to the remote peer
func (s SyncManager) SyncConfig(configXML []byte) error {
	fmt.Printf("Syncing configuration to %s via XMLRPC...\n", s.PeerIP)
	
	// Construct basic XMLRPC envelope (simplified)
	payload := fmt.Sprintf(`<?xml version="1.0"?>
<methodCall>
  <methodName>config.update</methodName>
  <params>
    <param><value><string>%s</string></value></param>
    <param><value><base64>%s</base64></value></param>
  </params>
</methodCall>`, s.Password, configXML)

	resp, err := http.Post(fmt.Sprintf("https://%s:8443/xmlrpc", s.PeerIP), "text/xml", bytes.NewBufferString(payload))
	if err != nil {
		return fmt.Errorf("sync request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sync failed with status: %d", resp.StatusCode)
	}

	return nil
}
