package services

import (
	"fmt"
	"time"
)

// PortalSession represents an authenticated captive portal user
type PortalSession struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	Mac       string    `json:"mac"`
	Username  string    `json:"username"`
	StartTime time.Time `json:"start_time"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CaptivePortal handles guest redirection and session state
type CaptivePortal struct {
	Interface string `json:"interface"`
	Port      int    `json:"port"` // Redirect port
	Sessions  map[string]PortalSession
}

// NewCaptivePortal initializes a portal on a specific interface
func NewCaptivePortal(iface string) *CaptivePortal {
	return &CaptivePortal{
		Interface: iface,
		Port:      8000,
		Sessions:  make(map[string]PortalSession),
	}
}

// GeneratePFRedirects creates pf rules to force HTTP traffic to the portal
func (p *CaptivePortal) GeneratePFRedirects() string {
	return fmt.Sprintf(`
# Captive Portal Redirects for %s
rdr on %s proto tcp from any to any port 80 -> 127.0.0.1 port %d
pass in quick on %s proto tcp from any to 127.0.0.1 port %d
`, p.Interface, p.Interface, p.Port, p.Interface, p.Port)
}

// CreateSession registers a new authenticated guest
func (p *CaptivePortal) CreateSession(ip, mac, username string, duration time.Duration) PortalSession {
	session := PortalSession{
		ID:        fmt.Sprintf("sess_%d", time.Now().UnixNano()),
		IP:        ip,
		Mac:       mac,
		Username:  username,
		StartTime: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	p.Sessions[ip] = session
	fmt.Printf("Captive Portal: Session created for %s (%s)\n", username, ip)
	return session
}

// CheckSession returns true if the IP has an active non-expired session
func (p *CaptivePortal) CheckSession(ip string) bool {
	session, ok := p.Sessions[ip]
	if !ok {
		return false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(p.Sessions, ip)
		return false
	}
	return true
}
