package system

import (
	"fmt"
	"os/exec"
)

// Diagnostics provides network troubleshooting tools
type Diagnostics struct{}

// Ping performs an ICMP echo request
func (d Diagnostics) Ping(host string, count int) (string, error) {
	cmd := exec.Command("ping", "-c", fmt.Sprintf("%d", count), host)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Traceroute traces the hop path to a host
func (d Diagnostics) Traceroute(host string) (string, error) {
	cmd := exec.Command("traceroute", host)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// CaptureTraffic initiates a brief tcpdump on an interface
func (d Diagnostics) CaptureTraffic(iface string, durationSec int) error {
	fmt.Printf("Capturing traffic on %s for %d seconds...\n", iface, durationSec)
	cmd := exec.Command("tcpdump", "-i", iface, "-w", "/tmp/capture.pcap")
	// In reality: run for duration and then kill
	return cmd.Start()
}
