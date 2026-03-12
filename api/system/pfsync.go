package system

import (
	"fmt"
	"os/exec"
)

// PfSyncConfig handles firewall state synchronization
type PfSyncConfig struct {
	SyncInterface string `json:"interface"`
	SyncPeer      string `json:"peer_ip"`
}

// EnablePfSync activates state synchronization across nodes
func EnablePfSync(config PfSyncConfig) error {
	// 1. Create pfsync interface
	exec.Command("ifconfig", "pfsync0", "create").Run()
	
	// 2. Set syncdev and peer
	cmd := exec.Command("ifconfig", "pfsync0", "syncdev", config.SyncInterface, "syncpeer", config.SyncPeer, "up")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to enable pfsync: %s", string(out))
	}
	
	return nil
}

// DisablePfSync tears down the sync interface
func DisablePfSync() error {
	cmd := exec.Command("ifconfig", "pfsync0", "destroy")
	return cmd.Run()
}
