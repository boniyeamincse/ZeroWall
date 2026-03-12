# ZeroWall — Security Model

Security is at the foundation of ZeroWall. This document describes the security architecture, privilege separation, hardening measures, and best practices implemented to ensure the integrity and safety of the firewall and the networks it protects.

---

## 1. Privilege Separation
ZeroWall employs a strict tiered privilege model to minimize the impact of a potential vulnerability in any single component.

- **Non-Privileged: `zwapi` (Web UI/API)**
  - Runs as a dedicated `zwapi` user.
  - Has NO direct access to kernel interfaces or privileged system commands.
  - Communicates with the configuration logic via restricted Unix sockets.
- **Intermediate: `zwconfigd` (Configuration Daemon)**
  - Runs with limited privileges required to read and write ZeroWall's `config.xml`.
  - Responsible for generating configuration files for system services (dhcpd, unbound, pf.conf).
- **Privileged: `zwsupervisor`**
  - Runs as `root`.
  - Only executes a pre-defined set of safe commands (pftcl, ifconfig, service restarts).
  - Validates all requests before execution.

## 2. System Hardening
The underlying FreeBSD system is hardened following industry best practices:

- **Minimal Attack Surface:** Only essential services are installed. Unused binaries and setuid files are removed.
- **Kernel Hardening:** Sysctl tuning to prevent IP spoofing, mitigate SYN floods, and disable ICMP redirects.
- **Filesystem Security:** The base system is mounted read-only during normal operation. Only specific directories (`/var`, `/etc/zerowall`) are writable.
- **Secure Defaults:** The firewall defaults to a "deny all" policy for all incoming traffic.

## 3. Administrative Security
- **HTTPS Enforcement:** The web dashboard is served over TLS by default.
- **Role-Based Access Control (RBAC):** Administrators can define granular permissions for users, limiting access to sensitive configuration areas.
- **Anti-Lockout:** A protected rule ensures that local network administrative access to the dashboard is never accidentally blocked.
- **JWT Authentication:** The API uses signed JSON Web Tokens with short expiration times for secure, stateless session management.

## 4. Encryption
- **VPN:** WireGuard (ChaCha20-Poly1305) and OpenVPN (AES-256-GCM) provide state-of-the-art encryption for remote access and site-to-site links.
- **Passwords:** All passwords in the configuration file are hashed using Argon2 or BCrypt.
- **Secrets Management:** Sensitive keys and certificates are stored with restrictive filesystem permissions (0600) and owned by the relevant service user.

## 5. IDS/IPS Integration
Inline Intrusion Prevention (Suricata) monitors all traffic passing through the firewall, identifying and blocking known threats, zero-day exploits, and malicious command-and-control (C2) patterns based on real-time threat intelligence.

---

*Previous: [Modules](06-modules.md) | Next: [VPN System](08-vpn-system.md)*
