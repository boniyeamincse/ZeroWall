# ZeroWall — Configuration Guide

This guide covers advanced configuration of ZeroWall after initial installation, including interface fine-tuning, firewall policy design, and service optimization.

---

## 1. Interface Management
Beyond basic assignments, ZeroWall allows for detailed interface customization:

- **Speed/Duplex:** Manual settings for legacy hardware or specific ISP handoffs.
- **MTU Manipulation:** Essential for certain VPN tunnels and PPPoE connections (typical WAN MTU: 1492 or 1500; typical WireGuard MTU: 1420).
- **MAC Spoofing:** Change the hardware address of an interface to satisfy ISP requirements.
- **MSS Clamping:** Automatically adjust TCP Maximum Segment Size to prevent fragmentation issues on restricted MTU links.

## 2. Advanced Firewall Policy
ZeroWall's `pf` implementation supports complex filtering logic:

- **Floating Rules:** Rules that apply globally before interface-specific rules. Used for blocking "Known Bad" IPs across all ports.
- **Policy Routing:** Direct traffic from specific internal hosts or subnets out of a specific gateway (e.g., routing all IoT traffic through a VPN).
- **Limiters:** Implement per-user or per-service bandwidth limits (Traffic Shaping).
- **Time-Based Rules:** Disable Guest WiFi access or specific server ports during off-hours.

## 3. High Availability (HA) Setup
For mission-critical environments, ZeroWall supports active/passive failover:

1. **CARP Base:** Configure an identical ZeroWall unit.
2. **Virtual IPs:** Create a shared Virtual IP (VIP) for WAN and LAN.
3. **Synchronization (XMLRPC):** Configure the master node to sync its configuration to the slave's IP.
4. **pfsync:** Enable state table synchronization so connections are not dropped during failover.

## 4. Intrusion Prevention Configuration
Fine-tuning Suricata for your network:

- **Promiscuous Mode:** Enable on interfaces where you want to monitor all traffic, not just traffic destined for the firewall.
- **Ruleset Categories:** Enable only relevant rulesets (e.g., Web Server rules for a DMZ interface) to improve performance.
- **Suppression:** Silence false positives by Suppressing specific SIDs (Signature IDs) in the alert log.

## 5. System Tunables (Sysctl)
ZeroWall provides a GUI for adjusting FreeBSD kernel parameters (`sysctl`).
- **Standard Tunables:** Buffer sizes, TCP timing, and security hardening.
- **Hardware-Specific:** Intel NIC offloading (TSO/LRO) settings.
- **Warning:** Incorrect sysctl values can lead to system instability; only adjust if instructed by documentation or support.

---

*Previous: [Installation Guide](12-installation-guide.md) | Next: [Developer Guide](14-developer-guide.md)*
