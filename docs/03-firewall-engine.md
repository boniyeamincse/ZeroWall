# ZeroWall — Firewall Engine

The heart of ZeroWall's security is its stateful firewall engine, powered by the FreeBSD **OpenBSD Packet Filter (pf)**. This document explains how the engine operates, how rules are structured, and how ZeroWall manages filtering logic.

---

## Architecture Overview

ZeroWall leverages the native FreeBSD `pf` kernel module for high-performance packet processing. The configuration is managed by the `zwconfigd` daemon, which translates user-defined rules in `config.xml` into a standard `pf.conf` format.

```
[ Web UI ] ──► [ config.xml ] ──► [ zwconfigd ] ──► [ /etc/pf.conf ] ──► [ pfctl ] ──► [ Kernel ]
```

## Stateful Inspection

Unlike simple packet filters, ZeroWall tracks the state of every connection. When a packet matches a "pass" rule, a state entry is created. Subsequent packets for that connection are automatically permitted, bypassing further rule evaluation.

### State Table Management
- **States:** Tracked by source/destination IP, protocol, ports, and sequence numbers.
- **Limits:** Configurable maximum state table size (default scales based on RAM).
- **Timeouts:** Automatic cleaning of idle or half-closed connections.

## Rule Structure & Processing

ZeroWall processes rules in a specific order:
1. **Normalization (Scrub):** Packets are reconstructed to prevent fragmentation attacks and normalized for processing.
2. **NAT / Binat / RDR:** Address and port translations are applied.
3. **Filtering Rules:**
   - Rules are evaluated on a **first-match** basis (unless the `quick` keyword is omitted).
   - ZeroWall uses "Quick" by default for all GUI-generated rules.

### Rule Types
| Type | Description |
|------|-------------|
| **Pass** | Permits traffic matching the criteria. |
| **Block** | Silently drops the packet. |
| **Reject** | Drops the packet and sends an ICMP Unreachable or TCP RST response. |
| **Floating** | Global rules that apply to multiple interfaces or directions. |

## Aliases and Groups

To simplify management, ZeroWall supports:
- **Host Aliases:** Single IPs or FQDNs.
- **Network Aliases:** Entire subnets.
- **Port Aliases:** Groups of TCP/UDP ports.
- **Interface Groups:** Apply rules to multiple physical/logical interfaces simultaneously.

## Advanced Filtering Features

### Outbound NAT (Masquerading)
ZeroWall automatically configures Outbound NAT for internal networks to access the internet via the WAN interface's public IP.

### Anti-Lockout Rule
By default, ZeroWall ensures that an "Anti-Lockout" rule exists on the LAN interface to prevent administrators from accidentally cutting off their own access to the web dashboard.

### Schedule-Based Rules
Rules can be enabled or disabled based on time-of-day or day-of-week schedules.

---

*Previous: [Features](02-features.md) | Next: [System Architecture](04-system-architecture.md)*
