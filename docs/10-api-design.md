# ZeroWall — API Design

The ZeroWall REST API (v1) allows for complete programmatic management and monitoring of the firewall system. This document outlines the design principles, authentication mechanisms, and core endpoints of the API.

---

## 1. Design Principles
- **RESTful Architecture:** Uses standard HTTP methods (GET, POST, PUT, DELETE, PATCH) to represent operations on resources.
- **Statelessness:** No session state is stored on the server. Every request must contain all information necessary for authentication and processing.
- **JSON Format:** All request bodies and responses are formatted as JSON.
- **Predictable URIs:** Resources are organized logically (e.g., `/api/v1/firewall/rules`, `/api/v1/vpn/peers`).
- **Idempotency:** Repeated GET, PUT, and DELETE requests to the same resource yield the same result (where applicable).

## 2. Authentication
ZeroWall API uses **JSON Web Tokens (JWT)** for secure authentication.

### Obtaining a Token
`POST /api/v1/auth/login`
- **Body:** `{ "username": "admin", "password": "..." }`
- **Response:** `{ "token": "ey...", "expires_at": "..." }`

### Using the Token
All subsequent requests must include the token in the `Authorization` header:
`Authorization: Bearer <TOKEN>`

## 3. Core Resource Endpoints

| Category | Endpoint Base | Description |
|----------|---------------|-------------|
| **System** | `/api/v1/system` | CPU, Memory, Uptime, Hostname, Updates |
| **Interfaces** | `/api/v1/interfaces` | List interfaces, status, IP configuration |
| **Firewall** | `/api/v1/firewall` | Rules management, NAT, Aliases, States |
| **VPN** | `/api/v1/vpn` | WireGuard/OpenVPN status and configuration |
| **IDS/IPS** | `/api/v1/ids` | Alerts, rulesets, and engine control |
| **Services** | `/api/v1/services` | DHCP leases, DNS resolver, NTP |
| **Logs** | `/api/v1/logs` | Real-time and historical log retrieval |

## 4. Example Request: Fetching Firewall Rules
`GET /api/v1/firewall/rules?interface=lan`

**Response:**
```json
{
  "total": 12,
  "rules": [
    {
      "id": "1001",
      "action": "pass",
      "interface": "lan",
      "source": "192.168.1.0/24",
      "destination": "any",
      "protocol": "any",
      "description": "Default allow LAN to any"
    }
  ]
}
```

## 5. WebSockets for Real-Time Data
For high-frequency telemetry and log streaming, ZeroWall provides WebSocket endpoints:
- `wss://<IP>/ws/v1/stats`: Live CPU/Network throughput.
- `wss://<IP>/ws/v1/logs`: Real-time firewall and system logs.
- `wss://<IP>/ws/v1/alerts`: Real-time IDS/IPS alert stream.

## 6. Error Handling
ZeroWall uses standard HTTP status codes for error reporting:
- **400 Bad Request:** Malformed JSON or invalid parameter.
- **401 Unauthorized:** Invalid or expired JWT.
- **403 Forbidden:** Authenticated user lacks sufficient permissions for the resource.
- **404 Not Found:** Resource does not exist.
- **500 Internal Server Error:** Unexpected server-side failure.

---

## 7. Detailed Documentation
For detailed endpoint specifications and design analysis, see:
- [API Endpoint Reference](api/endpoints.md)
- [API Design Analysis](api/design-analysis.md)

---

*Previous: [Web Dashboard](09-web-dashboard.md) | Next: [Database Design](11-database-design.md)*
