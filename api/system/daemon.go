package system

import (
	"fmt"
	"time"
)

// ConfigDaemon handles backend job execution
type ConfigDaemon struct {
	RefreshInterval int // Seconds
}

// Run starts the zwconfigd loop
func (d ConfigDaemon) Run() {
	fmt.Println("ZeroWall Configuration Daemon [zwconfigd] started.")
	
	// Main loop for applying periodic tasks or queueing changes
	for {
		// 1. Check change queue
		// 2. Apply pending pf rules
		// 3. Rotate logs
		time.Sleep(time.Duration(d.RefreshInterval) * time.Second)
	}
}

// TriggerReload sends a Signal or triggers a re-configuration event
func TriggerReload(module string) {
	fmt.Printf("Triggering reload for module: %s\n", module)
}
