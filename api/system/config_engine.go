package system

import (
	"fmt"
	"os"
)

// ConfigEngine handles atomic operations on config.xml
type ConfigEngine struct {
	ConfigPath string
}

// SaveAtomic writes configuration securely, maintaining a backup
func (e ConfigEngine) SaveAtomic(data []byte) error {
	backupPath := e.ConfigPath + ".bak"
	tempPath := e.ConfigPath + ".tmp"

	// 1. Create backup of current config
	if _, err := os.Stat(e.ConfigPath); err == nil {
		if err := os.Rename(e.ConfigPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
	}

	// 2. Write to temp file
	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temp config: %v", err)
	}

	// 3. Atomic rename (at least at filesystem level)
	if err := os.Rename(tempPath, e.ConfigPath); err != nil {
		return fmt.Errorf("atomic rename failed: %v", err)
	}

	fmt.Println("Configuration saved atomically.")
	return nil
}

// Rollback restores the previous configuration
func (e ConfigEngine) Rollback() error {
	backupPath := e.ConfigPath + ".bak"
	return os.Rename(backupPath, e.ConfigPath)
}
