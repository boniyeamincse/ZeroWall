
# ZeroWall — Web Dashboard

This document covers the ZeroWall web management interface: its architecture, navigation, key pages, role-based access control, and customization options.

---

## Table of Contents

1. [Dashboard Overview](#dashboard-overview)
2. [Technology Stack](#technology-stack)
3. [Navigation Structure](#navigation-structure)
4. [Dashboard Home](#dashboard-home)
5. [Firewall Rules Manager](#firewall-rules-manager)
6. [Real-Time Monitoring](#real-time-monitoring)
7. [VPN Management](#vpn-management)
8. [IDS/IPS Interface](#idsips-interface)
9. [System Configuration](#system-configuration)
10. [Role-Based Access Control](#role-based-access-control)
11. [Theming & Customization](#theming--customization)
12. [Keyboard Shortcuts](#keyboard-shortcuts)

---

## Dashboard Overview

The ZeroWall web dashboard provides a **unified management interface** for all firewall, VPN, routing, and monitoring functions. It is designed for both daily operations and initial deployment, offering:

- Responsive layout for desktop, tablet, and mobile
- Real-time data via WebSocket connections
- Wizard-driven setup for complex configurations
- Context-sensitive help integrated throughout

**Default access:** `https://192.168.1.1` (LAN interface, HTTPS)  
**Default credentials:** admin / zerowall (must be changed on first login)

---

## Technology Stack

| Component | Technology |
|-----------|-----------|
| Frontend framework | React 18 + TypeScript |
| UI components | Custom ZeroWall Design System |
| Styling | Tailwind CSS 3.x |
| State management | Zustand |
| Charts & graphs | Recharts + D3.js |
| Real-time data | WebSocket API |
| API communication | Axios + React Query |
| Build tool | Vite |

---

## Navigation Structure

```
ZeroWall Dashboard
│
├── Dashboard (Home)
│   ├── Overview
│   ├── Widgets (customizable)
│   └── Quick Actions
│
├── Firewall
│   ├── Rules
│   │   ├── WAN Rules
│   │   ├── LAN Rules
│   │   ├── Floating Rules
│   │   └── Rule Scheduler
│   ├── Aliases
│   ├── NAT
│   │   ├── Port Forwards
│   │   ├── Outbound NAT
│   │   └── 1:1 NAT
│   ├── Traffic Shaper
│   └── Virtual IPs
│
├── Interfaces
│   ├── Interface Assignments
│   ├── VLANs
│   ├── Bridges
│   ├── LAGGs
│   └── Wireless
│
├── Services
│   ├── DNS Resolver (Unbound)
│   ├── DHCP Server
│   ├── Dynamic DNS
│   ├── NTP
│   └── SNMP
│
├── VPN
│   ├── WireGuard
│   │   ├── Servers
│   │   └── Peers
│   ├── OpenVPN
│   │   ├── Servers
│   │   ├── Clients
│   │   └── Client Export
│   └── IPsec
│
├── IDS/IPS
│   ├── Overview
│   ├── Rules & Rulesets
│   ├── Suppression List
│   └── Alert Log
│
├── Routing
│   ├── Static Routes
│   ├── Gateways
│   ├── Gateway Groups
│   └── Dynamic Routing (FRR)
│
├── Monitoring
│   ├── Traffic Graphs
│   ├── System Health
│   ├── Firewall States
│   ├── DHCP Leases
│   └── Logs
│       ├── Firewall Log
│       ├── System Log
│       ├── VPN Log
│       └── IDS/IPS Log
│
├── System
│   ├── General Setup
│   ├── High Availability
│   ├── Certificates
│   ├── Backup & Restore
│   ├── Firmware Update
│   ├── Package Manager
│   └── Advanced
│       ├── Admin Access
│       ├── Firewall (pf options)
│       ├── Networking
│       └── System Tunables
│
└── Diagnostics
    ├── Ping / Traceroute
    ├── DNS Lookup
    ├── Packet Capture
    ├── pfInfo
    ├── Reboot / Shutdown
    └── Command Prompt (admin only)
```

---

## Dashboard Home

The dashboard home page displays configurable widgets:

### Default Widget Layout

```
┌──────────────────────────────────────────────────────────┐
│  ZeroWall  fw01.corp.example.com          [admin] [logout]│
├──────────────────────────────────────────────────────────┤
│ Dashboard │ Firewall │ Interfaces │ VPN │ IDS │ System   │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ┌─────────────────┐  ┌─────────────────┐               │
│  │  System Info    │  │  Interface Stats │               │
│  │  CPU:  12%      │  │  WAN:  45Mbps ↓ │               │
│  │  RAM:  1.2GB    │  │       8Mbps ↑   │               │
│  │  Disk: 4.1GB    │  │  LAN:  100Mbps  │               │
│  │  Uptime: 14d    │  │  wg0:  12Mbps   │               │
│  └─────────────────┘  └─────────────────┘               │
│                                                          │
│  ┌──────────────────────────────────────────┐           │
│  │  WAN Traffic (24h)                        │           │
│  │  ╭────────────────────────────────────╮  │           │
│  │  │   ▓▓        ▓▓▓▓    ▓▓       ░░   │  │           │
│  │  │ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░  │  │           │
│  │  ╰────────────────────────────────────╯  │           │
│  └──────────────────────────────────────────┘           │
│                                                          │
│  ┌─────────────────┐  ┌─────────────────┐               │
│  │  Firewall       │  │  IDS Alerts     │               │
│  │  States: 12,441 │  │  Last hour: 3   │               │
│  │  Rules: 47      │  │  Last 24h: 28   │               │
│  │  Blocked/hr: 82 │  │  [View Alerts]  │               │
│  └─────────────────┘  └─────────────────┘               │
│                                                          │
│  ┌──────────────────────────────────────────┐           │
│  │  Active VPN Connections                   │           │
│  │  alice-laptop  10.200.0.2  wg0  Online   │           │
│  │  bob-phone     10.200.0.3  wg0  Online   │           │
│  │  branch-office 10.200.0.10 wg0  Online   │           │
│  └──────────────────────────────────────────┘           │
└──────────────────────────────────────────────────────────┘
```

### Customizable Widgets

Available widgets (drag-and-drop layout):
- System Information
- Interface Traffic (per interface, selectable)
- Traffic Graph (line chart, configurable period)
- Firewall State Table summary
- Top Blocked IPs
- IDS/IPS Alert Summary
- Active VPN Sessions
- DHCP Lease Map
- Gateway Status
- CPU/Memory/Disk gauges
- Recent Log Entries

---

## Firewall Rules Manager

The rules manager provides full CRUD for firewall rules with drag-to-reorder functionality.

### Rules Table View

```
Interface: LAN    [+ Add Rule]  [+ Add Separator]
──────────────────────────────────────────────────────────
# │ En │ Action │ Proto │ Source      │ Dest       │ Port  │ Desc
──┼────┼────────┼───────┼─────────────┼────────────┼───────┼──────
⠿ │ ✓  │  PASS  │  any  │ LAN net     │ LAN addr   │  any  │ Anti-lockout
⠿ │ ✓  │  PASS  │  *    │ RFC1918     │ !RFC1918   │  any  │ LAN to WAN
⠿ │ ✓  │  PASS  │ TCP   │ LAN net     │ web_servers│ 443   │ HTTPS to DMZ
⠿ │ ✗  │  BLOCK │  *    │ IoT_VLAN    │ LAN net    │  any  │ IoT isolation
──────────────────────────────────────────────────────────
```

### Rule Edit Form

The rule editor uses a structured form with inline validation:

```
Edit Firewall Rule
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Action:       [● Pass] [○ Block] [○ Reject]
Disabled:     [ ] Disable this rule
Quick:        [ ] Stop processing on match
Interface:    [LAN ▼]
Direction:    [● In] [○ Out]
TCP/IP Ver:   [● IPv4] [○ IPv6] [○ Both]
Protocol:     [TCP ▼]

Source
  Type:       [● Network]
  Address:    [192.168.1.0/24___________]

Destination
  Type:       [● Single host or alias]
  Address:    [web_servers_____________]
  Port:       [443_____________________]

Options
  Log:        [✓] Log matching packets
  State:      [● Keep State]
  Description:[HTTPS to DMZ servers____]

[Save]  [Cancel]
```

---

## Real-Time Monitoring

### Traffic Graphs

Live bandwidth graphs update every second via WebSocket:

- **Interface view**: Per-interface In/Out rates (bps, pps)
- **Protocol breakdown**: TCP, UDP, ICMP percentages
- **Top Talkers**: Top 10 IPs by bandwidth (last 60 seconds)
- **Historical**: 1 hour, 24 hour, 1 week, 1 month views (stored in RRD or InfluxDB)

### Firewall State Browser

```
Filter: [src IP ________] [dst IP ________] [proto ▼] [Search]

Proto │ Source              │ Destination         │ State       │ Age  │ Bytes
──────┼─────────────────────┼─────────────────────┼─────────────┼──────┼───────
TCP   │ 192.168.1.50:51240  │ 93.184.216.34:443   │ ESTABLISHED │  4m  │ 142KB
TCP   │ 192.168.1.51:49871  │ 8.8.8.8:443         │ ESTABLISHED │  1m  │ 3.2KB
UDP   │ 192.168.1.1:33412   │ 1.1.1.1:53          │ MULTIPLE    │  0s  │ 120B
```

States are refreshed in real-time. Click a state to view full details or kill the connection.

### Packet Capture (Diagnostics)

```
Interface:    [WAN ▼]
Filter expr:  [host 8.8.8.8 and port 53____________]
Max packets:  [100_]

[Start Capture]  → Downloads .pcap file openable in Wireshark
```

---

## VPN Management

### WireGuard Peer Status

```
Server: wg0 (10.200.0.1/24) — Port 51820 — Active

Peer              Tunnel IP    Last Handshake  Transfer       Status
alice-laptop      10.200.0.2   2 min ago       ↓42MB ↑8MB    Online
bob-phone         10.200.0.3   14 min ago      ↓5MB  ↑2MB    Online
branch-office-gw  10.200.0.10  1 min ago       ↓1.2GB↑890MB  Online
charlie-laptop    10.200.0.4   4 days ago      ↓0    ↑0       Offline
```

### Client Config Export

For each WireGuard peer or OpenVPN user, ZeroWall generates:
- **QR Code** (displayed on screen, scan with WireGuard mobile app)
- **.conf file** (download for desktop WireGuard client)
- **.ovpn bundle** (OpenVPN — certificate + config in one file)
- **Installation guide** (per-platform instructions popup)

---

## IDS/IPS Interface

### Alert Stream

```
Live Alerts (updating)
────────────────────────────────────────────────────────────────
Time     │ Action  │ Signature                         │ Src IP
─────────┼─────────┼───────────────────────────────────┼────────────────
14:23:01 │ BLOCKED │ ET EXPLOIT EternalBlue SMB Probe  │ 198.51.100.42
14:22:58 │ ALERT   │ ET SCAN Nmap TCP SYN Scan          │ 10.0.0.99
14:22:45 │ BLOCKED │ ET MALWARE CobaltStrike Beacon     │ 192.168.30.15
14:21:10 │ ALERT   │ ET DNS Query for Known C2 Domain   │ 192.168.1.88
```

Clicking an alert shows:
- Full packet details
- Suricata rule text
- Options: Suppress this SID, block source IP, create firewall rule

### Ruleset Management

```
Ruleset              Version     Enabled    Last Updated    Alerts/24h
Emerging Threats     2024-11-15  Yes        15 Nov 14:00    28
ET Pro (requires key) —          No         —               —
Snort VRT           2024-11-10   Yes        10 Nov 08:00    4
Custom Rules        —            Yes        —               1

[Update All Rulesets]  [Add Custom Rule]  [Manage Suppression]
```

---

## System Configuration

### Backup & Restore

```
Configuration Backup
━━━━━━━━━━━━━━━━━━━
[Download Backup]   → Downloads config.xml (optionally encrypted)

Automatic Backups
  Schedule:         [Daily ▼]
  Retention:        [30 days]
  Remote Storage:   [SFTP ▼]
  Host:             [backup.example.com]
  Path:             [/backups/zerowall/]

[Save]
```

### Firmware Update

```
Current Version:  ZeroWall 24.11.1
Latest Version:   ZeroWall 24.11.2  (update available)
Release Notes:    [View]

Update Method:
  ● Download and Install automatically
  ○ Download only (install manually)
  ○ Upload image file

[Check for Updates]  [Install Update]
```

---

## Role-Based Access Control

### Built-In Roles

| Role | Access Level |
|------|-------------|
| Administrator | Full access to all configuration and operations |
| Operator | Can modify rules, VPN, monitoring; cannot change system/users |
| Read-Only | View-only access to all pages; cannot make changes |
| VPN Manager | Manage VPN servers and peers only |
| Helpdesk | View logs, monitoring, DHCP leases; no config changes |
| Custom | Granular per-page permissions defined by admin |

### Custom Role Configuration

```
Custom Role: "NOC Operator"
━━━━━━━━━━━━━━━━━━━━━━━━━━
Pages:
  [✓] Dashboard > Overview
  [✓] Monitoring > Traffic Graphs
  [✓] Monitoring > Logs (view only)
  [✓] IDS/IPS > Overview
  [ ] Firewall > Rules (no access)
  [ ] System > anything (no access)

[Save Role]
```

### User Management

Users can authenticate via:
1. Local user database
2. LDAP/Active Directory (group-to-role mapping)
3. RADIUS (role attribute in response)

---

## Theming & Customization

### Themes

ZeroWall ships with:
- **ZeroWall Dark** (default) — deep navy and cyan accent
- **ZeroWall Light** — clean white with blue accents
- **High Contrast** — accessibility-optimized

Custom CSS injection supported via System > Advanced > Custom CSS.

### Login Page Branding

Administrators can customize:
- Logo image
- Product name displayed
- Login page background
- Welcome message

Used by managed security service providers (MSSPs) to white-label the interface.

---

## Keyboard Shortcuts

| Shortcut | Action |
|---------|--------|
| `?` | Show keyboard shortcuts help |
| `g d` | Go to Dashboard |
| `g f` | Go to Firewall > Rules |
| `g v` | Go to VPN |
| `g m` | Go to Monitoring |
| `g s` | Go to System |
| `/` | Focus search |
| `Ctrl+S` | Save current form |
| `Ctrl+Z` | Revert unsaved changes |
| `Esc` | Close modal / cancel |

---

*Previous: [VPN System](08-vpn-system.md) | Next: [API Design](10-api-design.md)*
