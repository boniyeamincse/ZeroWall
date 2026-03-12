
# ZeroWall — System Architecture

This document describes the overall system architecture of ZeroWall, including its layered design, component interactions, process topology, and data flows.

---

## Table of Contents

1. [Architectural Overview](#architectural-overview)
2. [Layered Architecture](#layered-architecture)
3. [Process & Service Map](#process--service-map)
4. [Component Interaction Diagram](#component-interaction-diagram)
5. [Configuration Store](#configuration-store)
6. [Backend API Service](#backend-api-service)
7. [Frontend Web Application](#frontend-web-application)
8. [Kernel Subsystems](#kernel-subsystems)
9. [Filesystem Layout](#filesystem-layout)
10. [Boot Sequence](#boot-sequence)
11. [Security Boundaries](#security-boundaries)

---

## Architectural Overview

ZeroWall follows a **layered, service-oriented architecture** where a privileged backend daemon mediates all interactions between the web UI / API and kernel subsystems. No web-facing component runs as root or has direct kernel access.

```
┌──────────────────────────────────────────────────────────────────┐
│                        MANAGEMENT PLANE                          │
│                                                                  │
│   ┌─────────────┐   ┌──────────────┐   ┌───────────────────┐   │
│   │  Browser /  │   │  REST Client │   │  CLI (zwcli)      │   │
│   │  Web UI     │   │  (API users) │   │  Admin Console    │   │
│   └──────┬──────┘   └──────┬───────┘   └────────┬──────────┘   │
│          │                 │                     │               │
│   ┌──────▼─────────────────▼─────────────────────▼──────────┐   │
│   │              zwapi  (Backend API — Go)                   │   │
│   │         Runs as zwapi user, HTTPS :443 / :8080           │   │
│   └──────────────────────────┬───────────────────────────────┘   │
│                              │                                   │
│   ┌──────────────────────────▼───────────────────────────────┐   │
│   │           zwconfigd  (Configuration Daemon — Go)          │   │
│   │       Reads/writes /etc/zerowall/config.xml               │   │
│   │       Generates pf.conf, unbound.conf, etc.               │   │
│   └──────────────────────────┬───────────────────────────────┘   │
│                              │  (Unix socket, restricted)        │
└──────────────────────────────┼───────────────────────────────────┘
                               │
┌──────────────────────────────┼───────────────────────────────────┐
│                         DATA PLANE                                │
│                              │                                   │
│   ┌──────────────────────────▼───────────────────────────────┐   │
│   │              zwsupervisor  (privileged daemon)            │   │
│   │  Runs as root. Executes: pfctl, ifconfig, route, etc.     │   │
│   └──────┬──────┬──────┬──────┬──────┬──────┬────────────────┘   │
│          │      │      │      │      │      │                    │
│        pf    Suricata  WG  OpenVPN Unbound ISC-DHCP              │
│     (kernel) (IPS)   (VPN) (VPN)  (DNS)  (DHCP)                 │
│                                                                   │
│   ┌───────────────────────────────────────────────────────────┐   │
│   │            FreeBSD 14.x Hardened Kernel                   │   │
│   └───────────────────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────┘
```

---

## Layered Architecture

### Layer 1 — Hardware Abstraction
FreeBSD drivers and NIC support. ZeroWall supports Intel (em, igb, ix), Broadcom (bge, bce), Realtek (re, rge), VMware (vmx), VirtIO, and others.

### Layer 2 — FreeBSD Base System
Minimal, hardened FreeBSD installation. Unnecessary packages, daemons, and setuid binaries are removed. Base system is **read-only mounted** from a dedicated partition.

### Layer 3 — ZeroWall Kernel Extensions
- `/boot/loader.conf` tuning for network performance
- `if_wg` (WireGuard) kernel module
- `pfsync` and `carp` modules for HA
- `divert` socket support for IPS inline mode

### Layer 4 — ZeroWall Core Services
| Service | Binary | Role |
|---------|--------|------|
| zwsupervisor | `/usr/local/sbin/zwsupervisor` | Privileged operations executor |
| zwconfigd | `/usr/local/sbin/zwconfigd` | Configuration manager |
| zwapi | `/usr/local/sbin/zwapi` | REST API + web UI server |
| zwlogd | `/usr/local/sbin/zwlogd` | Log aggregator and exporter |

### Layer 5 — Third-Party Services
Suricata, Unbound, ISC-DHCP, OpenVPN, FRRouting — all managed by ZeroWall's rc.d scripts, started/stopped via `zwsupervisor`.

### Layer 6 — Management Interface
Web dashboard and REST API. All interactions flow through `zwapi` which validates, authorizes, and forwards to `zwconfigd`.

---

## Process & Service Map

```
PID   USER        COMMAND
1     root        init
├── zwsupervisor  root        privileged command executor
├── zwconfigd     zwcfg       configuration daemon
├── zwapi         zwapi       API + web server (HTTPS)
├── zwlogd        zwlog       log aggregator
├── suricata      suricata    IDS/IPS engine
├── unbound       unbound     DNS resolver
├── dhcpd         dhcpd       DHCP server
├── openvpn       nobody      OpenVPN daemon(s)
├── frr           frr         Dynamic routing (OSPF/BGP)
└── pflogd        root        pf log reader
```

---

## Component Interaction Diagram

```
Admin Browser
    │  HTTPS/WSS :443
    ▼
zwapi (Go HTTP server)
    │  Validates JWT, RBAC check
    │  Parses request body
    ▼
zwconfigd (Unix socket)
    │  Validates config change
    │  Writes to config.xml (atomic)
    │  Generates subsystem config files
    │
    ├──► /etc/pf.conf  ──► pfctl -f  (via zwsupervisor)
    ├──► /etc/unbound/unbound.conf  ──► unbound-control reload
    ├──► /etc/suricata/suricata.yaml  ──► suricatasc reload-rules
    ├──► /etc/openvpn/*.conf  ──► service openvpn restart
    └──► /etc/wg/*.conf  ──► wg syncconf wg0 /etc/wg/wg0.conf
```

---

## Configuration Store

ZeroWall stores its entire configuration in a single XML file: `/etc/zerowall/config.xml`.

### Design Rationale
- Single-file config enables atomic backup, restore, and replication to HA secondary
- Human-readable for audit and version control
- Git-committable — administrators can track changes with `git log`

### Config Schema (excerpt)
```xml
<?xml version="1.0" encoding="UTF-8"?>
<zerowall version="24.11">
  <system>
    <hostname>fw01</hostname>
    <domain>corp.example.com</domain>
    <timezone>America/New_York</timezone>
  </system>
  <interfaces>
    <wan>
      <if>em0</if>
      <type>dhcp</type>
      <descr>WAN</descr>
    </wan>
    <lan>
      <if>em1</if>
      <type>static</type>
      <ipaddr>192.168.1.1</ipaddr>
      <subnet>24</subnet>
    </lan>
  </interfaces>
  <filter>
    <rule>
      <id>1001</id>
      <action>pass</action>
      <interface>lan</interface>
      <direction>in</direction>
      <proto>tcp</proto>
      <dst_port>443</dst_port>
      <description>Allow HTTPS outbound from LAN</description>
    </rule>
  </filter>
</zerowall>
```

---

## Backend API Service

`zwapi` is written in **Go** and serves:
- The React-based web dashboard (static assets)
- The REST API (`/api/v1/*`)
- WebSocket endpoints for real-time telemetry (`/ws/stats`, `/ws/logs`)

### Key Design Decisions
- Stateless: JWT tokens carry all session info; no server-side sessions
- All mutations validated against JSON Schema before forwarding to `zwconfigd`
- Database is the config.xml file — no external database required for core functionality
- Optional PostgreSQL for extended logging/reporting (see [Database Design](11-database-design.md))

---

## Frontend Web Application

Built with **React 18** and **TypeScript**. Served as a compiled SPA from `/usr/local/www/zerowall/`.

- State management: Zustand
- UI components: Custom component library + Tailwind CSS
- Charts: Recharts + D3.js
- Real-time data: Native WebSocket API

---

## Kernel Subsystems

### Network Stack Tuning
ZeroWall applies the following `sysctl` values on boot:

```sh
# Increase network buffer sizes
kern.ipc.maxsockbuf=16777216
net.inet.tcp.recvbuf_max=16777216
net.inet.tcp.sendbuf_max=16777216

# Enable TCP SYN cookies (SYN flood mitigation)
net.inet.tcp.syncookies=1

# Disable ICMP redirects on WAN
net.inet.ip.redirect=0
net.inet6.ip6.redirect=0

# Enable IP forwarding
net.inet.ip.forwarding=1
net.inet6.ip6.forwarding=1

# Harden TCP
net.inet.tcp.blackhole=2
net.inet.udp.blackhole=1
```

---

## Filesystem Layout

```
/
├── boot/                    FreeBSD bootloader (read-only)
├── etc/
│   ├── zerowall/
│   │   ├── config.xml       Master configuration file
│   │   ├── config.xml.bak   Automatic backup
│   │   └── certs/           Certificates and keys
│   ├── pf.conf              Generated by zwconfigd
│   ├── unbound/             Generated Unbound config
│   └── dhcpd.conf           Generated DHCP config
├── usr/
│   ├── local/
│   │   ├── sbin/            ZeroWall daemons
│   │   ├── www/zerowall/    Web dashboard assets
│   │   └── etc/rc.d/        ZeroWall rc scripts
│   └── share/zerowall/      Templates, schemas
├── var/
│   ├── log/zerowall/        Application logs
│   ├── db/zerowall/         Runtime databases (DHCP leases, etc.)
│   └── run/zerowall/        PID files, Unix sockets
└── zw/                      ZeroWall package overlay (read-only)
```

---

## Boot Sequence

```
1. BIOS/UEFI POST
2. FreeBSD bootloader (loader.conf applied)
3. Kernel loads + NIC drivers attach
4. rc.d sequence:
   a. zwpre        — mount filesystems, apply sysctl
   b. zwnetwork    — configure interfaces (from config.xml)
   c. zwfirewall   — generate pf.conf, load rules (default-deny)
   d. zwconfigd    — start configuration daemon
   e. zwapi        — start API/web server
   f. suricata     — start IDS/IPS (if enabled)
   g. unbound      — start DNS resolver
   h. dhcpd        — start DHCP server
   i. openvpn/wg   — start VPN daemons (if configured)
   j. zwlogd       — start log aggregator
5. System ready — web UI available on LAN
```

---

## Security Boundaries

ZeroWall enforces strict privilege separation:

| Component | User | Capabilities |
|-----------|------|-------------|
| zwapi | zwapi | bind :443, read /etc/zerowall/certs/ |
| zwconfigd | zwcfg | read/write /etc/zerowall/config.xml |
| zwsupervisor | root | execute pfctl, ifconfig, route, wg |
| suricata | suricata | AF_PACKET socket, write /var/log/suricata/ |
| unbound | unbound | bind :53, read zone files |
| dhcpd | dhcpd | bind :67, write lease file |

Communication between components uses **Unix domain sockets** with filesystem-level permissions — no network-exposed inter-process communication for privileged operations.

---

*Previous: [Firewall Engine](03-firewall-engine.md) | Next: [Network Flow](05-network-flow.md)*
ENDOFFILE
Output

