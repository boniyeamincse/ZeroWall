# ZeroWall — Modules

ZeroWall is designed with a highly modular architecture, allowing core features to be managed as independent components. This document provides an overview of the primary modules that make up the ZeroWall system.

---

## 1. Network Interface Module
Manages physical and logical network interfaces.
- **Ethernet:** Configuration of WAN, LAN, and OPT interfaces.
- **VLANs:** Creation and management of 802.1Q tagged interfaces.
- **LAGG:** Support for Link Aggregation (LACP, Load Balance, Failover).
- **Wireless:** Access point and client mode configuration for supported wireless cards.

## 2. Firewall Module
The core filtering engine based on `pf`.
- **Rules Engine:** Management of inbound and outbound filtering rules.
- **Aliases:** Grouping of IPs, networks, and ports.
- **NAT:** Port forwarding and outbound network address translation.
- **Schedules:** Time-based rule activation.

## 3. VPN Module
Secure communication components.
- **WireGuard Module:** High-performance kernel-mode VPN.
- **OpenVPN Module:** Feature-rich SSL/TLS VPN.
- **IPsec Module:** Standards-based IKEv2 site-to-site connectivity.

## 4. Intrusion Detection & Prevention (IDS/IPS) Module
Threat analysis and mitigation.
- **Suricata Engine:** Deep Packet Inspection (DPI) and signature matching.
- **Rule Manager:** Automatic and manual ruleset updates.
- **GeoIP Module:** Filter traffic by country of origin or destination.

## 5. Network Services Module
Essential services for network operations.
- **Unbound DNS:** Secure recursive DNS resolver.
- **DHCP Server:** Dynamic IP assignment and static mappings.
- **NTP:** Network time synchronization.
- **Dynamic DNS:** Automatic IP updates for remote access.

## 6. High Availability (HA) Module
Redundancy and failover capabilities.
- **CARP:** Common Address Redundancy Protocol for virtual IPs.
- **pfsync:** State table synchronization between HA nodes.
- **XMLRPC Sync:** Configuration synchronization between primary and secondary units.

## 7. Management Module
User and API interfaces.
- **Web UI:** React-based management dashboard.
- **REST API:** Programmatic interface for system management.
- **zwconfigd:** The configuration daemon that orchestrates all other modules.

---

*Previous: [Network Flow](05-network-flow.md) | Next: [Security Model](07-security-model.md)*
