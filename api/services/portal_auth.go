package services

import (
	"fmt"
)

// PortalAuthBackend defines how the portal verifies credentials
type PortalAuthBackend string

const (
	AuthLocal  PortalAuthBackend = "local"
	AuthRADIUS PortalAuthBackend = "radius"
)

// AuthenticateGuest verifies guest credentials against the chosen backend
func AuthenticateGuest(username, password string, backend PortalAuthBackend) (bool, error) {
	fmt.Printf("Authenticating guest '%s' via %s backend...\n", username, backend)
	
	switch backend {
	case AuthLocal:
		// Logic to check local SQLite/XML user DB
		return true, nil
	case AuthRADIUS:
		// Logic to send ACCESS-REQUEST to RADIUS server
		return true, nil
	default:
		return false, fmt.Errorf("unsupported authentication backend: %s", backend)
	}
}
