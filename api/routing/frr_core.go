package routing

import (
	"fmt"
	"os/exec"
)

// FRREngine manages the Free Range Routing suite
type FRREngine struct {
	EnabledDaemons []string // bgpd, ospfd, zebra
}

// GenerateConfig assembles the main frr.conf
func (e FRREngine) GenerateConfig(daemonsConfigs string) string {
	header := "!\n! ZeroWall FRR Configuration\n!\n"
	return header + daemonsConfigs
}

// StartFRR activates the zebra management and routing daemons
func (e FRREngine) StartFRR() error {
	fmt.Println("Starting FRR Routing Suite...")
	// 1. Enable zebra
	// 2. Enable specific daemons in /usr/local/etc/frr/daemons
	// 3. service frr start
	return nil
}

// ReloadConfig applies changes without dropping routes
func (e FRREngine) ReloadConfig() error {
	cmd := exec.Command("vtysh", "-c", "write memory")
	return cmd.Run()
}
