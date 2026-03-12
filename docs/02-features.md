# ZeroWall — Features

ZeroWall is packed with enterprise-grade features designed to provide complete network security, visibility, and control. This document provides an in-depth look at the core capabilities of the system.

---

## 1. Firewall & Core Networking
- **Stateful Packet Inspection (pf):** Advanced connection tracking for IPv4 and IPv6, ensuring only legitimate traffic passes through your network.
- **NAT / PAT:** Support for Source NAT, Destination NAT (Port Forwarding), 1:1 NAT, and Outbound NAT rules.
- **Traffic Shaping (QoS):** Prioritize critical traffic (VoIP, Video) and limit low-priority applications using HFSC, PRIQ, or FairQ algorithms.
- **High Availability (CARP):** Active/Passive failover with configuration and state synchronization, ensuring zero downtime during hardware failure.
- **VLAN Support (802.1Q):** Logical network segmentation on a single physical interface.

## 2. VPN & Secure Remote Access
- **WireGuard:** Ultra-fast, modern VPN protocol utilizing state-of-the-art cryptography.
- **OpenVPN:** Highly flexible VPN supporting various authentication methods including LDAP, RADIUS, and Multi-Factor Authentication (MFA).
- **IPsec:** Standard-based VPN for site-to-site connectivity with third-party hardware.

## 3. Threat Prevention & IDS/IPS
- **Suricata Integration:** Real-time intrusion detection and prevention with Deep Packet Inspection (DPI).
- **Signature Updates:** Automatic updates from Emerging Threats (ET) and Snort rulesets.
- **Inline Mode:** Blocks malicious traffic before it reaches the internal network without requiring a separate network tap.
- **GeoIP Blocking:** Block or allow traffic based on geographic location (country-based firewalling).

## 4. Management & Monitoring
- **Web Dashboard:** A responsive, modern React-based interface for centralized management.
- **REST API:** Fully documented API for automation, orchestration, and integration with third-party tools.
- **Real-Time Telemetry:** Live graphs, connection state monitoring, and log streaming via WebSockets.
- **Role-Based Access Control (RBAC):** Granular permissions for different administrative roles (Admin, Operator, Read-Only).

## 5. Network Services
- **Unbound DNS:** Integrated secure recursive DNS resolver with DNSSEC support.
- **DHCP Server:** Scalable DHCP server with support for static mappings and multiple subnets.
- **Dynamic DNS:** Support for numerous providers (Cloudflare, DuckDNS, etc.) to handle dynamic WAN IPs.
- **NTP Server:** Synchronize time across your entire network.

---

*Previous: [Overview](01-overview.md) | Next: [Firewall Engine](03-firewall-engine.md)*
