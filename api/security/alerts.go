package security

import (
	"encoding/json"
	"fmt"
	"os"
)

// EveAlert represents a Suricata EVE JSON alert entry
type EveAlert struct {
	Timestamp string `json:"timestamp"`
	FlowID    int64  `json:"flow_id"`
	SrcIP     string `json:"src_ip"`
	SrcPort   int    `json:"src_port"`
	DestIP    string `json:"dest_ip"`
	DestPort  int    `json:"dest_port"`
	Proto     string `json:"proto"`
	Alert     struct {
		Action      string `json:"action"`
		Gid         int    `json:"gid"`
		SignatureID int    `json:"signature_id"`
		Rev         int    `json:"rev"`
		Signature   string `json:"signature"`
		Category    string `json:"category"`
		Severity    int    `json:"severity"`
	} `json:"alert"`
}

// ParseAlerts reads and decodes Suricata EVE JSON alerts
func ParseAlerts(filename string) ([]EveAlert, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open alert log: %v", err)
	}
	defer file.Close()

	var alerts []EveAlert
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var alert EveAlert
		if err := decoder.Decode(&alert); err != nil {
			// Skip malformed entries
			continue
		}
		if alert.Alert.Signature != "" {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}
