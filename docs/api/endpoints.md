# ZeroWall — API Endpoint Reference

This document provides a detailed reference for all ZeroWall REST API (v1) endpoints.

## 1. Authentication
Endpoints for user session management.

### Login
`POST /api/v1/login`
- **Request Body:**
  ```json
  {
    "username": "admin",
    "password": "your_password"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "token": "ey...",
    "expires_at": "2026-03-12T16:00:00Z"
  }
  ```

---

## 2. System Management
Monitor and configure core system settings.

### Get System Status
`GET /api/v1/status`
- **Response (200 OK):**
  ```json
  {
    "status": "online",
    "version": "1.0.0-Beta",
    "uptime": "2d 4h 12m",
    "cpu_usage": 12,
    "memory_usage": 45
  }
  ```

### Health Check
`GET /api/v1/health`
- **Response (200 OK):**
  ```json
  {
    "healthy": true,
    "timestamp": 1773302878
  }
  ```

---

## 3. Firewall
Configure rules, NAT, and traffic shaping.

### List Firewall Rules
`GET /api/v1/firewall/rules`
- **Query Parameters:**
  - `interface`: (Optional) Filter by interface (e.g., `wan`, `lan`).
- **Response (200 OK):**
  ```json
  {
    "rules": [
      {
        "id": "1",
        "action": "pass",
        "interface": "wan",
        "protocol": "tcp",
        "source": "any",
        "destination": "10.0.0.10",
        "port": 443,
        "description": "Allow HTTPS to Web Server"
      }
    ]
  }
  ```

### Create/Update Rule
`POST /api/v1/firewall/rules`
- **Request Body:**
  ```json
  {
    "action": "block",
    "interface": "lan",
    "protocol": "any",
    "source": "192.168.1.100",
    "destination": "any",
    "description": "Block infected host"
  }
  ```

---

## 4. VPN
Interface with WireGuard and OpenVPN services.

### Get VPN Status
`GET /api/v1/vpn/status`
- **Response (200 OK):**
  ```json
  {
    "wireguard": { "active_peers": 5, "total_transfer": "1.2GB" },
    "openvpn": { "active_clients": 2, "status": "running" }
  }
  ```

---

## 5. IDS/IPS (Security)
Monitor and manage the Suricata engine.

### Get IDS Alerts
`GET /api/v1/security/alerts`
- **Query Parameters:**
  - `limit`: Number of alerts to retrieve.
  - `severity`: Filter by severity (e.g., `high`, `medium`).

---

## 6. Real-Time Telemetry (WebSockets)
Endpoints for live streaming data.

- `wss://<IP>/ws/stats`: Live CPU/Network throughput updates.
- `wss://<IP>/ws/logs`: Live firewall and system log stream.
