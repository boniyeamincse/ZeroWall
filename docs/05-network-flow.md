cat > /home/claude/docs/05-network-flow.md << 'ENDOFFILE'
# ZeroWall — Network Flow

This document explains how network traffic flows through a ZeroWall deployment, covering physical and logical topologies, VLAN segmentation, routing decisions, and common deployment scenarios.

---

## Table of Contents

1. [Basic Network Topology](#basic-network-topology)
2. [Interface Types](#interface-types)
3. [Packet Flow — LAN to WAN](#packet-flow--lan-to-wan)
4. [Packet Flow — WAN to LAN (Port Forward)](#packet-flow--wan-to-lan-port-forward)
5. [VLAN Segmentation](#vlan-segmentation)
6. [Multi-WAN & Failover](#multi-wan--failover)
7. [DMZ Architecture](#dmz-architecture)
8. [VPN Traffic Flow](#vpn-traffic-flow)
9. [IPS Inline Traffic Path](#ips-inline-traffic-path)
10. [High Availability Traffic Flow](#high-availability-traffic-flow)
11. [Common Deployment Scenarios](#common-deployment-scenarios)

---

## Basic Network Topology

A minimal ZeroWall deployment has two interfaces:

```
          Internet
             │
       ┌─────┴─────┐
       │  ISP/Modem │
       └─────┬─────┘
             │  WAN (em0) — public IP
     ┌───────┴────────┐
     │   Z E R O W A L L   │
     └───────┬────────┘
             │  LAN (em1) — 192.168.1.1/24
       ┌─────┴──────┐
       │   Switch   │
       └─────┬──────┘
     ┌───────┼──────────┐
     │       │          │
   PC-1    PC-2       PC-3
```

ZeroWall sits between all trust zones and enforces policy at each boundary.

---

## Interface Types

| Type | Description | Example |
|------|-------------|---------|
| WAN | Upstream internet connection | DHCP, PPPoE, Static |
| LAN | Internal trusted network | 192.168.1.0/24 |
| OPT | Additional interfaces (DMZ, IoT, Guest) | Configurable |
| VLAN | 802.1Q tagged sub-interfaces | Multiple per physical NIC |
| Bridge | Layer-2 bridge between interfaces | Transparent firewall mode |
| LAGG | Link aggregation (LACP/failover) | Bonded NICs |
| Loopback | Internal routing | lo0, lo1 |
| VPN (WG/tun) | VPN tunnel endpoints | wg0, tun0 |

---

## Packet Flow — LAN to WAN

```
[PC: 192.168.1.50] → sends TCP SYN to 8.8.8.8:443

Step 1: Packet arrives on LAN interface (em1)
        src=192.168.1.50:54321  dst=8.8.8.8:443

Step 2: pf scrub — normalize packet

Step 3: State table lookup — no existing state (new connection)

Step 4: Rule evaluation (LAN interface, inbound direction)
        → Rule "Allow LAN to WAN HTTPS" matches → pass, keep state

Step 5: Routing — kernel routes to WAN default gateway (via em0)

Step 6: Outbound NAT on em0
        src=192.168.1.50:54321 → translated to → src=203.0.113.1:62145
        State entry created mapping the translation

Step 7: Packet leaves em0 toward internet

Return packet: 8.8.8.8:443 → 203.0.113.1:62145
Step 8: Arrives on em0
Step 9: State table lookup — HIT (maps to 192.168.1.50:54321)
Step 10: Reverse NAT applied, forwarded to em1 → PC
```

---

## Packet Flow — WAN to LAN (Port Forward)

```
[External: 198.51.100.42] → TCP SYN to 203.0.113.1:443

Step 1: Packet arrives on WAN interface (em0)

Step 2: pf scrub

Step 3: State table lookup — no existing state

Step 4: NAT (inbound) — port forward rule matches
        dst=203.0.113.1:443 → rewritten to → dst=192.168.1.100:443

Step 5: Rule evaluation — "Allow port forward 443 to web_server" → pass

Step 6: Routing — forward to LAN

Step 7: Packet delivered to 192.168.1.100

Return: 192.168.1.100:443 → 198.51.100.42:src_port
Step 8: Outbound on em0 — reverse NAT restores source to 203.0.113.1:443
```

---

## VLAN Segmentation

ZeroWall supports 802.1Q VLAN tagging for logical network segmentation on a single physical interface:

```
Physical NIC: em1 (trunk port to managed switch)
│
├── em1.10  — VLAN 10: Corporate LAN (192.168.10.0/24)
├── em1.20  — VLAN 20: Guest WiFi    (192.168.20.0/24)
├── em1.30  — VLAN 30: IoT Devices   (192.168.30.0/24)
└── em1.40  — VLAN 40: Server DMZ    (192.168.40.0/24)
```

Each VLAN interface gets its own firewall rule set. Default inter-VLAN policy is **block all** — administrators explicitly permit required traffic.

### Example Inter-VLAN Rules

```
# IoT VLAN cannot reach Corporate LAN
block in on em1.30 from 192.168.30.0/24 to 192.168.10.0/24

# Corporate LAN can reach Server DMZ on specific ports
pass in on em1.10 proto tcp from 192.168.10.0/24 to 192.168.40.0/24 port { 80 443 3306 }

# Guest WiFi: internet only
pass in on em1.20 from 192.168.20.0/24 to !192.168.0.0/16
block in on em1.20 from 192.168.20.0/24 to 192.168.0.0/16
```

---

## Multi-WAN & Failover

ZeroWall supports multiple WAN connections with automatic failover and load balancing:

```
                    ┌──────────────┐
                    │   ZeroWall   │
                    └──┬───────┬───┘
                       │       │
              ISP-A (em0)   ISP-B (em2)
              203.0.113.1   198.51.100.1
              (Primary)     (Failover)
```

### Gateway Groups

A **gateway group** defines failover/load-balance behavior:

| Mode | Behavior |
|------|---------|
| Failover | All traffic uses primary; switch to secondary only on failure |
| Load Balance | Traffic distributed across gateways (per-flow hashing) |
| Tier-based | Priority tiers; all same-tier gateways share load |

### Gateway Monitoring

ZeroWall uses ICMP probes (configurable target, interval, threshold) to determine gateway health. On failure:
1. Gateway marked down in routing table
2. States using failed gateway are cleared
3. Traffic rerouted through next available gateway
4. Alert sent to administrator

---

## DMZ Architecture

A DMZ (demilitarized zone) isolates public-facing servers from the internal network:

```
Internet
    │
    │ WAN (em0)  203.0.113.0/30
    │
┌───┴────────────────────────────┐
│           ZeroWall              │
│                                │
│  Rules:                        │
│  WAN→DMZ: allow :80,:443,:25   │
│  DMZ→LAN: block all            │
│  LAN→DMZ: allow all            │
│  DMZ→WAN: allow :80,:443       │
└───┬────────────┬───────────────┘
    │            │
LAN (em1)      DMZ (em2)
192.168.1.0/24  10.0.0.0/24
    │            │
  Clients    Web/Mail/DNS Servers
```

### DMZ Security Policy

- Servers in DMZ have no route to LAN hosts
- Servers in DMZ may reach internet for updates (filtered)
- If a DMZ server is compromised, LAN remains protected
- Management access to DMZ servers from LAN only, on specific ports

---

## VPN Traffic Flow

### WireGuard Site-to-Site

```
Site A (ZeroWall-A)              Site B (ZeroWall-B)
192.168.1.0/24                   192.168.2.0/24
    │                                │
    │    wg0 tunnel (UDP 51820)      │
    ├────────────────────────────────┤
    │    Encrypted WireGuard packets │
    └────────────────────────────────┘
         over internet (NAT OK)
```

Traffic flow for PC-A (192.168.1.50) → PC-B (192.168.2.50):
1. PC-A sends to 192.168.2.50 — routes to ZeroWall-A
2. ZeroWall-A: routing table has `192.168.2.0/24 via wg0`
3. Packet encrypted, encapsulated in WireGuard UDP
4. Sent to ZeroWall-B's WAN IP on port 51820
5. ZeroWall-B decrypts, decapsulates
6. Routes decapsulated packet to 192.168.2.50

### Road Warrior VPN

```
Mobile User                     ZeroWall HQ
[WireGuard client]              wg0 server (10.10.0.1/24)
Gets peer IP: 10.10.0.2         LAN: 192.168.1.0/24
    │
    │  All traffic through tunnel (split or full tunnel)
    │  WireGuard UDP port 51820
    └──────────────────────────────────────────────────►
                                Routes traffic:
                                10.10.0.2 → internet via NAT
                                10.10.0.2 → 192.168.1.0/24 direct
```

---

## IPS Inline Traffic Path

In IPS mode, Suricata inspects all packets via a pf **divert socket**:

```
Incoming packet on em0
        │
        ▼
pf rule: "divert-to 127.0.0.1 port 8000"
        │
        ▼
Suricata (listening on divert socket :8000)
        │
    ┌───┴───┐
    │ Inspect│
    └───┬───┘
   Pass │  Drop (rule match)
        │         │
        ▼         ▼
  Re-injected   Dropped + logged
  into pf        to EVE JSON
        │
        ▼
  Normal pf rule processing continues
```

This "bump-in-the-wire" approach provides true **inline prevention** without requiring a separate network tap.

---

## High Availability Traffic Flow

### CARP Virtual IPs

```
          Internet
               │
               │  203.0.113.1 (CARP VIP — floats to primary)
               │
    ┌──────────┴──────────┐
    │                     │
┌───┴──────┐         ┌────┴─────┐
│ ZeroWall │ ◄──SYNC──► ZeroWall│
│  PRIMARY │ pfsync   │SECONDARY│
│ (MASTER) │ (state)  │(BACKUP) │
└───┬──────┘         └────┬─────┘
    │                     │
    └──────────┬───────────┘
               │  192.168.1.1 (CARP VIP)
               │
             LAN clients
```

- **Primary** holds all CARP VIPs (lower advskew value wins)
- **Secondary** monitors CARP advertisements; takes over if primary goes silent
- **pfsync** keeps state tables in sync — TCP connections survive failover
- Failover time: typically 1-3 seconds

---

## Common Deployment Scenarios

### Scenario 1: Home / Small Office
```
Modem/ONT → ZeroWall (1 WAN, 1 LAN) → Switch → Devices
```
Single WAN, single LAN. DHCP from ZeroWall. Basic firewall rules. Optional VPN for remote access.

### Scenario 2: Small Business with DMZ
```
Modem → ZeroWall (1 WAN, 1 LAN, 1 DMZ) → { Internal Switch, DMZ Switch }
```
Web/mail servers in DMZ. Firewall separates internal users from servers.

### Scenario 3: Multi-Site Enterprise
```
HQ ZeroWall ←── WireGuard ──→ Branch ZeroWall
     │                               │
   HQ LAN                       Branch LAN
   BGP via FRR ←──────────────→ BGP via FRR
```
Site-to-site VPN with dynamic routing. OSPF or BGP for automatic route propagation.

### Scenario 4: Cloud Gateway
```
ZeroWall VM on cloud (AWS/GCP/Azure/Hetzner)
  │  WAN: cloud public IP
  │  LAN: VPC/private network
  │
  VPN termination for remote workers
  Firewall for cloud workloads
```

### Scenario 5: High Availability Pair
```
     Internet
         │
    ┌────┴─────┐
    │  Layer 2  │  (switch or upstream router)
    └────┬─────┘
         │
    CARP VIP (floats)
    ┌────┴───────────┐
 ZW-Primary       ZW-Secondary
    │   pfsync/CARP  │
    └────────────────┘
         │
    Internal CARP VIP
         │
       LAN Switch
```

---

*Previous: [System Architecture](04-system-architecture.md) | Next: [Modules](06-modules.md)*
ENDOFFILE
Output