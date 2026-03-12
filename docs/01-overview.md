# ZeroWall — Project Overview

> **ZeroWall** is a free, open-source, enterprise-grade network firewall and security gateway built on FreeBSD.  
> It delivers stateful packet inspection, VPN termination, intrusion detection/prevention, and centralized network security management through a unified web dashboard.

---

## Table of Contents

1. [What is ZeroWall?](#what-is-zerowall)
2. [Design Philosophy](#design-philosophy)
3. [Key Capabilities](#key-capabilities)
4. [Target Audience](#target-audience)
5. [Project Status](#project-status)
6. [How ZeroWall Compares](#how-zerowall-compares)
7. [License](#license)
8. [Community & Support](#community--support)

---

## What is ZeroWall?

ZeroWall is a **network operating system** purpose-built for firewall and gateway roles. It transforms commodity x86 hardware (or virtual machines) into a hardened, feature-rich security appliance capable of protecting networks of any size — from small home labs to multi-site enterprise deployments.

ZeroWall is not a software package installed on top of a general-purpose OS. It is a **complete, integrated system image** that boots directly from disk or USB, presenting administrators with a secure, minimal attack surface and a consistent operational environment.

```
┌─────────────────────────────────────────────────────────────┐
│                        ZeroWall System                      │
│                                                             │
│  ┌───────────┐  ┌────────────┐  ┌──────────┐  ┌────────┐  │
│  │  Firewall │  │    VPN     │  │  IDS/IPS │  │  DHCP  │  │
│  │  Engine   │  │  (WG/OVPN) │  │(Suricata)│  │  DNS   │  │
│  └───────────┘  └────────────┘  └──────────┘  └────────┘  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Web Management Dashboard               │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │         FreeBSD 14.x Hardened Base System           │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## Design Philosophy

ZeroWall is built around five core principles:

| Principle | Description |
|-----------|-------------|
| **Security First** | Every component operates with least privilege. Attack surface is minimized by default. |
| **Transparency** | All source code is open. Configuration is human-readable. Behavior is auditable. |
| **Simplicity** | Complex security tasks are accessible without sacrificing depth or control. |
| **Performance** | Packet processing leverages FreeBSD's native network stack and kernel-level filtering. |
| **Extensibility** | A modular plugin architecture lets operators add capabilities without forking the project. |

---

## Key Capabilities

- **Stateful Packet Filtering** — Full connection tracking with IPv4/IPv6 support via `pf`.
- **NAT / PAT** — Network and port address translation with hairpin NAT support.
- **VPN Gateway** — WireGuard and OpenVPN server/client with road-warrior and site-to-site configs.
- **IDS/IPS** — Inline intrusion detection powered by Suricata with Emerging Threats rulesets.
- **High Availability** — CARP-based active/passive failover with configuration synchronization.
- **Traffic Shaping** — HFSC/PRIQ/CODEL queuing for QoS enforcement.
- **DNS & DHCP** — Integrated Unbound DNS resolver and ISC DHCP server.
- **Web UI** — Role-based, responsive management dashboard with real-time monitoring.
- **REST API** — Full programmatic control for automation and integration.
- **Logging & SIEM** — Structured syslog, Elasticsearch export, and Netflow/IPFIX support.

---

## Target Audience

### System Administrators
Operators deploying and maintaining network security at the perimeter or between segments. ZeroWall provides GUI-driven workflows for common tasks and CLI access for advanced operations.

### Security Engineers
Professionals designing defense-in-depth architectures. ZeroWall exposes deep packet inspection, rule tuning, and threat intelligence integration.

### Developers & Contributors
Engineers extending ZeroWall's capabilities or integrating it with third-party platforms. A documented plugin API and developer guide are provided.

### Home Lab Enthusiasts
Advanced users building private networks who want enterprise-class features without enterprise pricing.

---

## Project Status

| Component | Status |
|-----------|--------|
| Firewall Engine (pf) | Stable |
| Web Dashboard | Stable |
| WireGuard VPN | Stable |
| OpenVPN | Stable |
| IDS/IPS (Suricata) | Stable |
| REST API v1 | Stable |
| High Availability (CARP) | Stable |
| Traffic Shaping | Stable |
| Plugin System | Beta |
| Zero Trust Network Access | Roadmap |
| SD-WAN Integration | Roadmap |

---

## How ZeroWall Compares

| Feature | ZeroWall | pfSense CE | OPNsense |
|---------|----------|-----------|---------|
| Base OS | FreeBSD 14 | FreeBSD 14 | FreeBSD 13 |
| License | MIT | Apache 2.0 | BSD 2-Clause |
| WireGuard Native | Yes | Yes | Yes |
| REST API (full) | Yes | Partial | Yes |
| Plugin SDK | Yes | Limited | Yes |
| Inline IPS | Yes | Yes | Yes |
| Zero Trust (ZTNA) | Roadmap | No | No |
| Configuration as Code | Yes | No | No |

---

## License

ZeroWall is released under the **MIT License**.  
See `LICENSE` for full terms.

Third-party components (FreeBSD, Suricata, WireGuard, Unbound, etc.) retain their respective licenses.

---

## Community & Support

| Channel | Location |
|---------|----------|
| GitHub | https://github.com/zerowall/zerowall |
| Documentation | https://docs.zerowall.io |
| Forums | https://forum.zerowall.io |
| Bug Reports | GitHub Issues |
| Security Issues | security@zerowall.io |
| Chat | #zerowall on Libera.Chat / Matrix |

---

*Next: [Features](02-features.md)*